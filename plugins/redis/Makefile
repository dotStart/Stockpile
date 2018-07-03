PLUGIN_VERSION := "1.0"
LDFLAGS := -ldflags "-X github.com/dotStart/Stockpile/plugins/redis.version=${PLUGIN_VERSION}"

UNAME := $(shell uname)
ifeq ($(UNAME), Linux)
PLATFORMS := linux/amd64/so
else
PLATFORMS :=
endif

# magical formula:
temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))
ext = $(word 3, $(temp))

build: $(PLATFORMS)

$(PLATFORMS):
	@export GOOS=$(os); export GOARCH=$(arch); export CGO_ENABLED=1; go build -v -buildmode=plugin ${LDFLAGS} -o ../../build/$(os)-$(arch)/plugins/redis.$(ext)

.PHONY: build
