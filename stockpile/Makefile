PLATFORMS := darwin/386 darwin/amd64 linux/386 linux/amd64 linux/arm windows/386/.exe windows/amd64/.exe
LDFLAGS :=-ldflags "-X github.com/dotStart/Stockpile/stockpile/metadata.brand=${APPLICATION_BRAND} -X github.com/dotStart/Stockpile/stockpile/metadata.version=${APPLICATION_VERSION} -X github.com/dotStart/Stockpile/stockpile/metadata.commitHash=${APPLICATION_COMMIT_HASH} -X \"github.com/dotStart/Stockpile/stockpile/metadata.timestampRaw=${APPLICATION_TIMESTAMP}\""

# magical formula:
temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))
ext = $(word 3, $(temp))

build: $(PLATFORMS)

$(PLATFORMS):
	@export GOOS=$(os); export GOARCH=$(arch); go build -v ${LDFLAGS} -o ../build/$(os)-$(arch)/stockpile$(ext)

.PHONY: build
