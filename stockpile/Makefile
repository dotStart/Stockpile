LDFLAGS :=-ldflags "-X github.com/dotStart/Stockpile/metadata.brand=${APPLICATION_BRAND} -X github.com/dotStart/Stockpile/metadata.version=${APPLICATION_VERSION} -X github.com/dotStart/Stockpile/metadata.commitHash=${APPLICATION_COMMIT_HASH} -X \"github.com/dotStart/Stockpile/metadata.timestampRaw=${APPLICATION_TIMESTAMP}\""

build: ../build/mac32/stockpile ../build/mac64/stockpile \
	   ../build/linux32/stockpile ../build/linux64/stockpile ../build/linuxarm/stockpile \
	   ../build/win32/stockpile.exe ../build/win64/stockpile.exe

../build/mac32/stockpile:
	@export GOOS="darwin"; export GOARCH="386"; go build -v ${LDFLAGS} -o ../build/mac32/stockpile
	@echo ""

../build/mac64/stockpile:
	@export GOOS="darwin"; export GOARCH="amd64"; go build -v ${LDFLAGS} -o ../build/mac64/stockpile
	@echo ""

../build/linux32/stockpile:
	@export GOOS="linux"; export GOARCH="386"; go build -v ${LDFLAGS} -o ../build/linux32/stockpile
	@echo ""

../build/linux64/stockpile:
	@export GOOS="linux"; export GOARCH="amd64"; go build -v ${LDFLAGS} -o ../build/linux64/stockpile
	@echo ""

../build/linuxarm/stockpile:
	@export GOOS="linux"; export GOARCH="arm"; go build -v ${LDFLAGS} -o ../build/linuxarm/stockpile
	@echo ""

../build/win32/stockpile.exe:
	@export GOOS="windows"; export GOARCH="386"; go build -v ${LDFLAGS} -o ../build/win32/stockpile.exe
	@echo ""

../build/win64/stockpile.exe:
	@export GOOS="windows"; export GOARCH="amd64"; go build -v ${LDFLAGS} -o ../build/win64/stockpile.exe
	@echo ""

.PHONY: build
