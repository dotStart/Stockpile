PLATFORMS := darwin/386 darwin/amd64 linux/386 linux/amd64 linux/arm windows/386 windows/amd64
LDFLAGS :=-ldflags "-X github.com/dotStart/Stockpile/metadata.brand=${APPLICATION_BRAND} -X github.com/dotStart/Stockpile/stockpile/metadata.version=${APPLICATION_VERSION} -X github.com/dotStart/Stockpile/stockpile/metadata.commitHash=${APPLICATION_COMMIT_HASH} -X \"github.com/dotStart/Stockpile/stockpile/metadata.timestampRaw=${APPLICATION_TIMESTAMP}\""

# magical formula:
temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

build: $(PLATFORMS)

$(PLATFORMS):
	@export GOOS=$(os); export GOARCH=$(arch); go build -v ${LDFLAGS} -o ../build/$(os)-$(arch)/stockpile

.PHONY: build
