# Copyright © 2017 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

YEAR = 2017
HEADER = $(COMMENT) Copyright © 2017 The Things Network Foundation, distributed under the MIT license (see LICENSE file)
HEADER_EXTRA_FILES = Makefile

PRE_COMMIT = headers.check-staged
COMMIT_MSG = git.commit-msg-log git.commit-msg-length git.commit-msg-empty git.commit-msg-prefix git.commit-msg-phrase git.commit-msg-casing

include .make/log.make
include .make/general.make
include .make/git.make
include .make/headers.make
include .make/go/main.make

ci.encrypt-variables:
	keybase encrypt -b -i ci/variables.yml -o ci/variables.yml.encrypted johanstokking htdvisser romeovs ericgo

ci.decrypt-variables:
	keybase decrypt -i ci/variables.yml.encrypted -o ci/variables.yml

messages:
	@$(GO) run ./pkg/errors/generate_messages.go --filename config/messages.json

dev-deps: go.dev-deps

deps: go.deps

test: go.test

quality: go.quality

# Cache build artifacts (speeds up dev builds)
cache: go.install

# Example builds
ttn-example: MAIN=./cmd/ttn-example/main.go
ttn-example: $(RELEASE_DIR)/ttn-example-$(GOOS)-$(GOARCH)
ttn-example: go.link

ttn-example-dev: MAIN=./cmd/ttn-example/main.go
ttn-example-dev: $(RELEASE_DIR)/ttn-example-$(GOOS)-$(GOARCH)-dev
ttn-example-dev: go.link

ttn-example-docker: MAIN=./cmd/ttn-example/main.go
ttn-example-docker: $(RELEASE_DIR)/ttn-example-linux-amd64

build-all: ttn-example
