.PHONY: build clean install

VERSION ?= 1.5.4
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_TIME ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

LDFLAGS := -X main.version=$(VERSION) \
           -X main.commit=$(COMMIT) \
           -X main.buildTime=$(BUILD_TIME)

ifeq ($(OS),Windows_NT)
  BIN_EXT := .exe
  MKDIR_CMD := if not exist dist mkdir dist
  RM_FILE := del /F /Q
  RM_DIR := rmdir /S /Q
  DIST_PRIMARY := dist\indus$(BIN_EXT)
  RSRC_FILE := cmd\indus-terminal\rsrc.syso
  NULLDEV := nul
else
  BIN_EXT :=
  MKDIR_CMD := mkdir -p dist
  RM_FILE := rm -f
  RM_DIR := rm -rf
  DIST_PRIMARY := dist/indus$(BIN_EXT)
  RSRC_FILE := cmd/indus-terminal/rsrc.syso
  NULLDEV := /dev/null
endif

build:
	$(MKDIR_CMD)
	-rsrc -ico build/icon.ico -o cmd/indus-terminal/rsrc.syso
	go build -ldflags "$(LDFLAGS) -H windowsgui" -o dist/indus$(BIN_EXT) ./cmd/indus-terminal

install:
	@echo Use install.bat on Windows or copy dist/indus$(BIN_EXT) into your PATH on Unix-like systems.

clean:
	-$(RM_FILE) $(DIST_PRIMARY) 2>$(NULLDEV)
	-$(RM_FILE) $(RSRC_FILE) 2>$(NULLDEV)
	-$(RM_DIR) dist 2>$(NULLDEV)
