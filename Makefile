################################################################################

# This Makefile generated by GoMakeGen 1.0.0 using next command:
# gomakegen --dep .
#
# More info: https://kaos.sh/gomakegen

################################################################################

.DEFAULT_GOAL := help
.PHONY = fmt all clean git-config deps dep-init dep-update help

################################################################################

all: sonar ## Build all binaries

sonar: ## Build sonar binary
	go build sonar.go

install: ## Install binaries
	cp sonar /usr/bin/sonar

uninstall: ## Uninstall binaries
	rm -f /usr/bin/sonar

git-config: ## Configure git redirects for stable import path services
	git config --global http.https://pkg.re.followRedirects true

deps: git-config dep-update ## Download dependencies

dep-init: ## Initialize dep workspace
	which dep &>/dev/null || (echo -e '\e[31mDep is not installed\e[0m' ; exit 1)
	dep init

dep-update: ## Update packages and dependencies through dep
	which dep &>/dev/null || (echo -e '\e[31mDep is not installed\e[0m' ; exit 1)
	test -s Gopkg.toml || dep init
	dep ensure -update

fmt: ## Format source code with gofmt
	find . -name "*.go" -exec gofmt -s -w {} \;

clean: ## Remove generated files
	rm -f sonar

help: ## Show this info
	@echo -e '\nSupported targets:\n'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[33m%-12s\033[0m %s\n", $$1, $$2}'
	@echo -e ''

################################################################################
