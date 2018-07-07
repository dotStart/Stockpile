PLATFORMS := darwin/386 darwin/amd64 linux/386 linux/amd64 linux/arm windows/386/.exe windows/amd64/.exe
LDFLAGS :=-ldflags "-X github.com/dotStart/Stockpile/stockpile/metadata.brand=${APPLICATION_BRAND} -X github.com/dotStart/Stockpile/stockpile/metadata.version=${APPLICATION_VERSION} -X github.com/dotStart/Stockpile/stockpile/metadata.commitHash=${APPLICATION_COMMIT_HASH} -X \"github.com/dotStart/Stockpile/stockpile/metadata.timestampRaw=${APPLICATION_TIMESTAMP}\""

DF = $(shell command -v df 2> /dev/null)

# magical formula:
temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))
ext = $(word 3, $(temp))

build: build-ui $(PLATFORMS)

.ONESHELL:
build-ui:
	@cd ui
	@NPM_FLAGS=
ifdef DF
	@if $(DF) -t vboxsf . > /dev/null 2> /dev/null; then \
		echo "VBox filesystem detected - Enabling compatibility flags"; \
		NPM_FLAGS=--no-bin-links; \
	fi
endif
	@"$(NPM)" install $(NPM_FLAGS)
	@"$(NPM)" run build

$(PLATFORMS):
	@export GOOS=$(os); export GOARCH=$(arch); $(GO) generate -v github.com/dotStart/Stockpile/stockpile/server/rpc github.com/dotStart/Stockpile/stockpile/server/ui && $(GO) build -v ${LDFLAGS} -o ../build/$(os)-$(arch)/stockpile$(ext)

.PHONY: build
