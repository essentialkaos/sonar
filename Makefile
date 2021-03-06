################################################################################

# This Makefile generated by GoMakeGen 1.3.0 using next command:
# gomakegen --dep .
#
# More info: https://kaos.sh/gomakegen

################################################################################

.DEFAULT_GOAL := help
.PHONY = fmt vet all clean git-config deps dep-init dep-update help

################################################################################

all: sonar ## Build all binaries

sonar: ## Build sonar binary
	go build sonar.go

install: ## Install all binaries
	cp sonar /usr/bin/sonar

uninstall: ## Uninstall all binaries
	rm -f /usr/bin/sonar

git-config: ## Configure git redirects for stable import path services
	git config --global http.https://pkg.re.followRedirects true

deps: git-config dep-update ## Download dependencies

dep-init: ## Initialize dep workspace
	which dep &>/dev/null || go get -u -v github.com/golang/dep/cmd/dep
	dep init

dep-update: ## Update packages and dependencies through dep
	which dep &>/dev/null || go get -u -v github.com/golang/dep/cmd/dep
	test -s Gopkg.toml || dep init
	test -s Gopkg.lock && dep ensure -update || dep ensure

fmt: ## Format source code with gofmt
	find . -name "*.go" -exec gofmt -s -w {} \;

vet: ## Runs go vet over sources
	go vet -composites=false -printfuncs=LPrintf,TLPrintf,TPrintf,log.Debug,log.Info,log.Warn,log.Error,log.Critical,log.Print ./...

clean: ## Remove generated files
	rm -f sonar

help: ## Show this info
	@echo -e '\n\033[1mSupported targets:\033[0m\n'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[33m%-12s\033[0m %s\n", $$1, $$2}'
	@echo -e ''
	@echo -e '\033[90mGenerated by GoMakeGen 1.3.0\033[0m\n'

################################################################################
