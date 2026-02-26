.PHONY: build clean install

VERSION ?= 1.0.0
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_TIME ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -X main.version=$(VERSION) \
           -X main.commit=$(COMMIT) \
           -X main.buildTime=$(BUILD_TIME)

build:
	rsrc -ico build/icon.ico -o cmd/indus-terminal/rsrc.syso
	go build -ldflags "$(LDFLAGS)" -o indus.exe ./cmd/indus-terminal

install:
	@echo Run install.bat as administrator to install INDUS Terminal

clean:
	del /F /Q indus.exe 2>nul
	del /F /Q cmd\indus-terminal\rsrc.syso 2>nul
	rmdir /S /Q dist 2>nul
