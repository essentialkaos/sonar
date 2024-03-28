package daemon

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2024 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"os"

	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/jsonutil"
	"github.com/essentialkaos/ek/v12/knf"
	"github.com/essentialkaos/ek/v12/log"
	"github.com/essentialkaos/ek/v12/options"
	"github.com/essentialkaos/ek/v12/pid"
	"github.com/essentialkaos/ek/v12/signal"
	"github.com/essentialkaos/ek/v12/support"
	"github.com/essentialkaos/ek/v12/support/deps"
	"github.com/essentialkaos/ek/v12/terminal/tty"
	"github.com/essentialkaos/ek/v12/usage"

	knfv "github.com/essentialkaos/ek/v12/knf/validators"
	knff "github.com/essentialkaos/ek/v12/knf/validators/fs"

	"github.com/essentialkaos/sonar/slack"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Basic info
const (
	APP  = "Sonar"
	VER  = "1.8.2"
	DESC = "Utility for showing user Slack status in JIRA"
)

// Options
const (
	OPT_CONFIG   = "c:config"
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"

	OPT_VERB_VER = "vv:verbose-version"
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
	OPT_HELP:     {Type: options.BOOL},
	OPT_VER:      {Type: options.MIXED},

	OPT_VERB_VER: {Type: options.BOOL},
}

var (
	enabled  bool
	mappings map[string]string
	bots     map[string]bool
	token    []byte
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Run is main daemon function
func Run(gitRev string, gomod []byte) {
	preConfigureUI()

	_, errs := options.Parse(optMap)

	if len(errs) != 0 {
		printError(errs[0].Error())
		os.Exit(1)
	}

	configureUI()

	switch {
	case options.GetB(OPT_VER):
		genAbout(gitRev).Print(options.GetS(OPT_VER))
		os.Exit(0)
	case options.GetB(OPT_HELP):
		genUsage().Print()
		os.Exit(0)
	case options.GetB(OPT_VERB_VER):
		support.Collect(APP, VER).
			WithRevision(gitRev).
			WithDeps(deps.Extract(gomod)).
			Print()
		os.Exit(0)
	}

	loadConfig()
	validateConfig()
	registerSignalHandlers()
	setupLogger()
	createPidFile()

	log.Divider()
	log.Aux("%s %s starting…", APP, VER)

	enabled = knf.GetB(MAIN_ENABLED, true)

	start()
}

// preConfigureUI preconfigures UI based on information about user terminal
func preConfigureUI() {
	if !tty.IsTTY() {
		fmtc.DisableColors = true
	}
}

// configureUI configures user interface
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}
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

// printErrorAndExit print error message and exit with exit code 1
func printErrorAndExit(f string, a ...interface{}) {
	printError(f, a...)
	os.Exit(1)
}

// shutdown stops daemon
func shutdown(code int) {
	pid.Remove(PID_FILE)
	os.Exit(code)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo()

	info.AddOption(OPT_CONFIG, "Path to configuration file", "file")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	return info
}

// genAbout generates info about version
func genAbout(gitRev string) *usage.About {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2006,
		Owner:   "ESSENTIAL KAOS",
		License: "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
	}

	if gitRev != "" {
		about.Build = "git:" + gitRev
	}

	return about
}
