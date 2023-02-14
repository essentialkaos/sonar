package daemon

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2022 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"os"
	"strings"

	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/jsonutil"
	"github.com/essentialkaos/ek/v12/knf"
	"github.com/essentialkaos/ek/v12/log"
	"github.com/essentialkaos/ek/v12/options"
	"github.com/essentialkaos/ek/v12/pid"
	"github.com/essentialkaos/ek/v12/signal"
	"github.com/essentialkaos/ek/v12/usage"

	knfv "github.com/essentialkaos/ek/v12/knf/validators"
	knff "github.com/essentialkaos/ek/v12/knf/validators/fs"

	"github.com/essentialkaos/sonar/slack"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Basic info
const (
	APP  = "Sonar"
	VER  = "1.8.0"
	DESC = "Utility for showing user Slack status in JIRA"
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
	MAIN_BOTS     = "main:bots"
	MAIN_TOKEN    = "main:token"
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

var (
	enabled  bool
	mappings map[string]string
	bots     map[string]bool
	token    []byte
)

// ////////////////////////////////////////////////////////////////////////////////// //

func Init() {
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
	log.Aux("%s %s starting…", APP, VER)

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
	errs := knf.Validate([]*knf.Validator{
		{MAIN_TOKEN, knfv.Empty, nil},
		{SLACK_TOKEN, knfv.Empty, nil},
		{LOG_DIR, knfv.Empty, nil},
		{LOG_FILE, knfv.Empty, nil},
		{HTTP_PORT, knfv.Empty, nil},

		{HTTP_PORT, knfv.Less, 1024},
		{HTTP_PORT, knfv.Greater, 65535},

		{SLACK_TOKEN, knfv.NotPrefix, "xoxb-"},

		{LOG_DIR, knff.Perms, "DW"},
		{LOG_DIR, knff.Perms, "DX"},

		{LOG_LEVEL, knfv.NotContains, []string{"debug", "info", "warn", "error", "crit"}},
	})

	if len(errs) != 0 {
		printError("Error while configuration file validation:")

		for _, err := range errs {
			printError("  %v", err)
		}

		os.Exit(1)
	}
}

// registerSignalHandlers registers signal handlers
func registerSignalHandlers() {
	signal.Handlers{
		signal.TERM: termSignalHandler,
		signal.INT:  intSignalHandler,
		signal.HUP:  hupSignalHandler,
	}.TrackAsync()
}

// setupLogger setups logger
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

// loadMappings loads mappings data
func loadMappings() {
	if knf.GetS(MAIN_MAPPINGS) == "" {
		return
	}

	mappings = make(map[string]string)

	err := jsonutil.Read(knf.GetS(MAIN_MAPPINGS), &mappings)

	if err != nil {
		log.Error(err.Error())
	}

	if len(mappings) != 0 {
		for orig, alias := range mappings {
			log.Info("Added mapping %s → %s", orig, alias)
		}
	}
}

// loadBots loads bots emails
func loadBots() {
	if knf.GetS(MAIN_BOTS) == "" {
		return
	}

	bots = make(map[string]bool)

	err := jsonutil.Read(knf.GetS(MAIN_BOTS), &bots)

	if err != nil {
		log.Error(err.Error())
	}
}

// createPidFile creates PID file
func createPidFile() {
	pid.Dir = PID_DIR

	err := pid.Create(PID_FILE)

	if err != nil {
		printErrorAndExit(err.Error())
	}
}

// start starts the service
func start() {
	loadMappings()
	loadBots()

	err := slack.StartObserver(knf.GetS(SLACK_TOKEN), mappings)

	if err != nil {
		log.Crit(err.Error())
		shutdown(1)
	}

	token = []byte(knf.GetS(MAIN_TOKEN))

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

// intSignalHandler is INT signal handler
func intSignalHandler() {
	log.Aux("Received INT signal, shutdown…")
	shutdown(0)
}

// termSignalHandler is TERM signal handler
func termSignalHandler() {
	log.Aux("Received TERM signal, shutdown…")
	shutdown(0)
}

// hupSignalHandler is HUP signal handler
func hupSignalHandler() {
	log.Info("Received HUP signal, log will be reopened…")
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

// shutdown stops deamon
func shutdown(code int) {
	pid.Remove(PID_FILE)
	os.Exit(code)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// showUsage prints usage info
func showUsage() {
	info := usage.NewInfo()

	info.AddOption(OPT_CONFIG, "Path to configuraion file", "file")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VERSION, "Show version")

	info.Render()
}

// showAbout prints information about license and version
func showAbout() {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2006,
		Owner:   "ESSENTIAL KAOS",
		License: "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
	}

	about.Render()
}
