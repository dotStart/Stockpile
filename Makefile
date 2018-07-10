APPLICATION_BRAND := vanilla
APPLICATION_VERSION := 2.0-alpha.1
APPLICATION_COMMIT_HASH := `git log -1 --pretty=format:"%H"`
APPLICATION_TIMESTAMP := `date --utc "+%s"`

GIT := $(shell command -v git 2> /dev/null)
DEP := $(shell command -v dep 2> /dev/null)
GO := $(shell command -v go 2> /dev/null)
NODE := $(shell command -v node 2> /dev/null)
NPM := $(shell command -v npm 2> /dev/null)
export

PLUGINS := $(wildcard plugins/*/.)

all: check-env print-config install-dependencies core core-plugins package

check-env:
	@echo "==> Checking prerequisites"
	@echo -n "Checking for git ... "
ifndef GIT
	@echo "Not found"
	$(error "git is unavailable")
endif
	@echo $(GIT)
	@echo -n "Checking for go ... "
ifndef GO
	@echo "Not Found"
	$(error "go is unavailable")
endif
	@echo $(GO)
	@echo -n "Checking for dep ... "
ifndef DEP
	@echo "Not Found"
	$(error "dep is unavailable")
endif
	@echo $(DEP)
	@echo -n "Checking for node ... "
ifndef NODE
	@echo "Not Found"
	$(error "node is unavailable")
endif
	@echo $(NODE)
	@echo -n "Checking for npm ... "
ifndef NPM
	@echo "Not Found"
	$(error "npm is unavailable")
endif
	@echo $(NPM)
	@echo ""

print-config:
	@echo "==> Build Configuration"
	@echo ""
	@echo "       Brand: ${APPLICATION_BRAND}"
	@echo "     Version: ${APPLICATION_VERSION}"
	@echo "  Commit SHA: ${APPLICATION_COMMIT_HASH}"
	@echo "   Timestamp: ${APPLICATION_TIMESTAMP}"
	@echo ""

clean:
	@echo "==> Clearing previous build data"
	@rm -rf build/ || true
	@$(GO) clean -cache

install-dependencies:
	@echo "==> Installing dependencies"
	@$(GO) get -u github.com/elazarl/go-bindata-assetfs/go-bindata-assetfs
	@$(GO) get -u github.com/jteeuwen/go-bindata/go-bindata
	@$(DEP) ensure -v
	@echo ""

core:
	@echo "==> Building stockpile"
	$(MAKE) -C stockpile/

core-plugins:
	@echo "==> Building core plugins"
	@for dir in $(PLUGINS); do \
        "$(MAKE)" -C $$dir; \
    done

package:
	@echo "==> Creating distribution packages"
	@for dir in build/*; do if [ -d "$$dir" ]; then tar -czvf "$(basename "$$dir").tar.gz" --xform="s,$$dir/,," "$$dir"; fi; done
	@echo ""

.PHONY: all
