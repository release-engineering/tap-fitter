SHELL = /bin/bash
GO := go
OBJ := tap-fitter
SPECIFIC_UNIT_TEST := $(if $(TEST),-run $(TEST),)
extra_env := $(GOENV)

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: clean
clean: ## Remove binaries and test artifacts
	@rm -rf $(OBJ) coverage.out

.PHONY: fmt
fmt: ## Run go fmt against code.
	$(GO) fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	$(GO) vet -tags '$(GO_BUILD_TAGS)' ./...

.PHONY: test
test-unit:
	$(GO) test -coverprofile=coverage.out $(SPECIFIC_UNIT_TEST) -count=1 ./...

.PHONY: tidy
tidy: ## Update dependencies
	$(GO) mod tidy

.PHONY: lint
lint:
	find . -name '*.go' | xargs goimports -w

.PHONY: verify
verify: tidy fmt vet lint ## Verify the current code and lint
	git diff --exit-code

##@ Build

$(OBJ):
	$(extra_env) $(GO) build $(extra_flags) -o $@ ./cmd/tap-fitter

.PHONY: build
build: clean $(OBJ)
