MODULE_PATH := github.com/blumsicle/reqserv
APP_NAME := $(shell basename $(MODULE_PATH))

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
VERSION := $(shell basename $(BRANCH))
COMMIT := $(shell git rev-parse --short HEAD)

LDFLAGS ?= '-X $(MODULE_PATH)/cmd.Name=$(APPNAME) -X $(MODULE_PATH)/cmd.Version=$(VERSION) -X $(MODULE_PATH)/cmd.Commit=$(COMMIT)'

SRC_PATH := .
DEST_PATH := ./bin/$(APP_NAME)-$(VERSION)

install: generate
	CGO_ENABLED=0 go install -ldflags $(LDFLAGS) $(SRC_PATH)

generate:
	go generate ./...


