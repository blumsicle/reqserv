MODULE_PATH := github.com/blumsicle/reqserv
APP_NAME    := $(shell basename $(MODULE_PATH))

DOCKER_USER ?= blumsicle8
PLATFORM ?= $(shell uname -m)

BRANCH  := $(shell git rev-parse --abbrev-ref HEAD)
VERSION := $(shell basename $(BRANCH))
COMMIT  := $(shell git rev-parse --short HEAD)

LDFLAGS ?= '-X $(MODULE_PATH)/cmd.Name=$(APP_NAME)  \
		   -X $(MODULE_PATH)/cmd.Version=$(VERSION) \
		   -X $(MODULE_PATH)/cmd.Commit=$(COMMIT)'

SRC_PATH  := .
DEST_PATH := ./bin/$(APP_NAME)-$(VERSION)

install: generate
	CGO_ENABLED=0 go install -ldflags $(LDFLAGS) $(SRC_PATH)

generate:
	go generate ./...

deps:
	go mod download

build: generate
	CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -o $(DEST_PATH) $(SRC_PATH)

docker-build:
	docker build . \
		--load \
		--platform linux/$(PLATFORM) \
		--tag $(DOCKER_USER)/$(APP_NAME):$(VERSION) \
		--tag $(DOCKER_USER)/$(APP_NAME):latest

docker-push: docker-build
	docker push --all-tags $(DOCKER_USER)/$(APP_NAME)

docker-push-multi:
	docker build . \
		--push \
		--platform linux/arm64,linux/amd64 \
		--tag $(DOCKER_USER)/$(APP_NAME):$(VERSION) \
		--tag $(DOCKER_USER)/$(APP_NAME):latest

.PHONY: install generate deps build docker-build docker-push
