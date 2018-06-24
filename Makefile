APPLICATION_BRAND := vanilla
APPLICATION_VERSION := 2.0.0
APPLICATION_COMMIT_HASH := `git log -1 --pretty=format:"%H"`
APPLICATION_TIMESTAMP := `date --utc "+%s"`

LDFLAGS :=-ldflags "-X github.com/dotStart/Stockpile/metadata.brand=${APPLICATION_BRAND} -X github.com/dotStart/Stockpile/metadata.version=${APPLICATION_VERSION} -X github.com/dotStart/Stockpile/metadata.commitHash=${APPLICATION_COMMIT_HASH} -X \"github.com/dotStart/Stockpile/metadata.timestampRaw=${APPLICATION_TIMESTAMP}\""

.PHONY: print-config install-dependencies generate-sources build package

clean:
	@echo "==> Clearing previous build data"
	@rm -rf build/ || true
	@go clean -cache

print-config:
	@echo "==> Build Configuration"
	@echo ""
	@echo "       Brand: ${APPLICATION_BRAND}"
	@echo "     Version: ${APPLICATION_VERSION}"
	@echo "  Commit SHA: ${APPLICATION_COMMIT_HASH}"
	@echo "   Timestamp: ${APPLICATION_TIMESTAMP}"
	@echo ""
	@echo "  Linker Flags: ${LDFLAGS}"
	@echo ""

install-dependencies:
	@echo "==> Installing dependencies"
	@dep ensure -v
	@echo ""

generate-sources:
	@echo "==> Generating protobuf sources"
	# todo
	@echo ""

build: build/mac32/stockpile build/mac64/stockpile \
	   build/linux32/stockpile build/linux64/stockpile build/linuxarm/stockpile \
	   build/win32/stockpile.exe build/win64/stockpile.exe

build/mac32/stockpile:
	@echo "==> Compiling Stockpile for Mac OS (32-Bit)"
	@export GOOS="darwin"; export GOARCH="386"; go build -v ${LDFLAGS} -o build/mac32/stockpile
	@echo ""

build/mac64/stockpile:
	@echo "==> Compiling Stockpile for Mac OS (64-Bit)"
	@export GOOS="darwin"; export GOARCH="amd64"; go build -v ${LDFLAGS} -o build/mac64/stockpile
	@echo ""

build/linux32/stockpile:
	@echo "==> Compiling Stockpile for Linux (32-Bit)"
	@export GOOS="linux"; export GOARCH="386"; go build -v ${LDFLAGS} -o build/linux32/stockpile
	@echo ""

build/linux64/stockpile:
	@echo "==> Compiling Stockpile for Linux (64-Bit)"
	@export GOOS="linux"; export GOARCH="amd64"; go build -v ${LDFLAGS} -o build/linux64/stockpile
	@echo ""

build/linuxarm/stockpile:
	@echo "==> Compiling Stockpile for Linux (ARM)"
	@export GOOS="linux"; export GOARCH="arm"; go build -v ${LDFLAGS} -o build/linuxarm/stockpile
	@echo ""

build/win32/stockpile.exe:
	@echo "==> Compiling Stockpile for Windows (32-Bit)"
	@export GOOS="windows"; export GOARCH="386"; go build -v ${LDFLAGS} -o build/win32/stockpile.exe
	@echo ""

build/win64/stockpile.exe:
	@echo "==> Compiling Stockpile for Windows (64-Bit)"
	@export GOOS="windows"; export GOARCH="amd64"; go build -v ${LDFLAGS} -o build/win64/stockpile.exe
	@echo ""

package:
	@echo "==> Creating distribution packages"
	@for dir in build/*; do if [ -d "$$dir" ]; then tar -czvf "$(basename "$$dir").tar.gz" --xform="s,$$dir/,," "$$dir"; fi; done
	@echo ""
