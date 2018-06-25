APPLICATION_BRAND := vanilla
APPLICATION_VERSION := 2.0.0
APPLICATION_COMMIT_HASH := `git log -1 --pretty=format:"%H"`
APPLICATION_TIMESTAMP := `date --utc "+%s"`
export

PLUGINS := $(wildcard plugins/*/.)

all: print-config install-dependencies core core-plugins package

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
	@go clean -cache

install-dependencies:
	@echo "==> Installing dependencies"
	@dep ensure -v
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
