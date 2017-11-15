package daemon

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                     Copyright (c) 2009-2017 ESSENTIAL KAOS                         //
//        Essential Kaos Open Source License <https://essentialkaos.com/ekol>         //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"os"
	"runtime"
	"strings"

	"pkg.re/essentialkaos/ek.v9/fmtc"
	"pkg.re/essentialkaos/ek.v9/fsutil"
	"pkg.re/essentialkaos/ek.v9/jsonutil"
	"pkg.re/essentialkaos/ek.v9/knf"
	"pkg.re/essentialkaos/ek.v9/log"
	"pkg.re/essentialkaos/ek.v9/options"
	"pkg.re/essentialkaos/ek.v9/pid"
	"pkg.re/essentialkaos/ek.v9/signal"
	"pkg.re/essentialkaos/ek.v9/usage"

	"github.com/essentialkaos/sonar/slack"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Basic info
const (
	APP  = "Sonar"
	VER  = "1.0.0"
	DESC = "Utility for showing user Slack status in Jira"
)

// Options
const (
	OPT_CONFIG   = "c:config"
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VERSION  = "v:version"
)

// Configuration file props
const (
	MAIN_ENABLED  = "main:enabled"
	MAIN_MAPPINGS = "main:mappings"
	SLACK_TOKEN   = "slack:token"
	HTTP_IP       = "http:ip"
	HTTP_PORT     = "http:port"
	LOG_DIR       = "log:dir"
	LOG_FILE      = "log:file"
	LOG_PERMS     = "log:perms"
	LOG_LEVEL     = "log:level"
)

// Pid info
const (
	PID_DIR  = "/var/run/sonar"
	PID_FILE = "sonar.pid"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Options map
var optMap = options.Map{
	OPT_CONFIG:   {Value: "/etc/sonar.knf"},
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL, Alias: "u:usage"},
	OPT_VERSION:  {Type: options.BOOL, Alias: "ver"},
}

var enabled bool
var mappings map[string]string

// ////////////////////////////////////////////////////////////////////////////////// //

func Init() {
	runtime.GOMAXPROCS(8)

	_, errs := options.Parse(optMap)

	if len(errs) != 0 {
		for _, err := range errs {
			printError(err.Error())
		}

		os.Exit(1)
	}

	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}

	if options.GetB(OPT_VERSION) {
		showAbout()
		return
	}

	if options.GetB(OPT_HELP) {
		showUsage()
		return
	}

	loadConfig()
	validateConfig()
	registerSignalHandlers()
	setupLogger()
	createPidFile()

	log.Aux(strings.Repeat("-", 88))
	log.Aux("%s %s starting...", APP, VER)

	enabled = knf.GetB(MAIN_ENABLED, true)

	start()
}

// loadConfig read and parse configuration file
func loadConfig() {
	err := knf.Global(options.GetS(OPT_CONFIG))

	if err != nil {
		printErrorAndExit(err.Error())
	}
}

// validateConfig validate configuration file values
func validateConfig() {
	var permsChecker = func(config *knf.Config, prop string, value interface{}) error {
		if !fsutil.CheckPerms(value.(string), config.GetS(prop)) {
			switch value.(string) {
			case "DW":
				return fmtc.Errorf("Property %s must be path to writable directory", prop)
			case "DX":
				return fmtc.Errorf("Property %s must be path to executable directory", prop)
			}
		}

		return nil
	}

	errs := knf.Validate([]*knf.Validator{
		{SLACK_TOKEN, knf.Empty, nil},
		{LOG_DIR, knf.Empty, nil},
		{LOG_FILE, knf.Empty, nil},
		{HTTP_PORT, knf.Empty, nil},

		{HTTP_PORT, knf.Less, 1024},
		{HTTP_PORT, knf.Greater, 65535},

		{SLACK_TOKEN, knf.NotLen, 42},
		{SLACK_TOKEN, knf.NotPrefix, "xoxb-"},

		{LOG_DIR, permsChecker, "DW"},
		{LOG_DIR, permsChecker, "DX"},
		{LOG_LEVEL, knf.NotContains, []string{"debug", "info", "warn", "error", "crit"}},
	})

	if len(errs) != 0 {
		printError("Error while configuration file validation:")

		for _, err := range errs {
			printError("  %v", err)
		}

		os.Exit(1)
	}
}

// registerSignalHandlers register signal handlers
func registerSignalHandlers() {
	signal.Handlers{
		signal.TERM: termSignalHandler,
		signal.INT:  intSignalHandler,
		signal.HUP:  hupSignalHandler,
	}.TrackAsync()
}

// setupLogger setup logger
func setupLogger() {
	err := log.Set(knf.GetS(LOG_FILE), knf.GetM(LOG_PERMS, 644))

	if err != nil {
		printErrorAndExit(err.Error())
	}

	err = log.MinLevel(knf.GetS(LOG_LEVEL))

	if err != nil {
		printErrorAndExit(err.Error())
	}
}

// loadMappings load mappings data
func loadMappings() {
	if knf.GetS(MAIN_MAPPINGS) == "" {
		return
	}

	mappings = make(map[string]string)

	err := jsonutil.DecodeFile(knf.GetS(MAIN_MAPPINGS), &mappings)

	if err != nil {
		log.Error(err.Error())
	}
}

// createPidFile create PID file
func createPidFile() {
	pid.Dir = PID_DIR

	err := pid.Create(PID_FILE)

	if err != nil {
		printErrorAndExit(err.Error())
	}
}

// start start service
func start() {
	loadMappings()

	err := slack.StartObserver(knf.GetS(SLACK_TOKEN), mappings)

	if err != nil {
		log.Crit(err.Error())
		shutdown(1)
	}

	err = startHTTPServer(
		knf.GetS(HTTP_IP),
		knf.GetS(HTTP_PORT),
	)

	if err != nil {
		log.Crit(err.Error())
		shutdown(1)
	}

	shutdown(0)
}

// INT signal handler
func intSignalHandler() {
	log.Aux("Received INT signal, shutdown...")
	shutdown(0)
}

// TERM signal handler
func termSignalHandler() {
	log.Aux("Received TERM signal, shutdown...")
	shutdown(0)
}

// HUP signal handler
func hupSignalHandler() {
	log.Info("Received HUP signal, log will be reopened...")
	log.Reopen()
	log.Info("Log reopened by HUP signal")
}

// printError prints error message to console
func printError(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{r}"+f+"{!}\n", a...)
}

// printError prints warning message to console
func printWarn(f string, a ...interface{}) {
	fmtc.Fprintf(os.Stderr, "{y}"+f+"{!}\n", a...)
}

// printErrorAndExit print error mesage and exit with exit code 1
func printErrorAndExit(f string, a ...interface{}) {
	printError(f, a...)
	os.Exit(1)
}

// shutdown stop deamon
func shutdown(code int) {
	pid.Remove(PID_FILE)
	os.Exit(code)
}

// ////////////////////////////////////////////////////////////////////////////////// //

func showUsage() {
	info := usage.NewInfo()

	info.AddOption(OPT_CONFIG, "Path to configuraion file", "file")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VERSION, "Show version")

	info.Render()
}

func showAbout() {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2006,
		Owner:   "ESSENTIAL KAOS",
		License: "Essential Kaos Open Source License <https://essentialkaos.com/ekol>",
	}

	about.Render()
}
