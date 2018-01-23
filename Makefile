# Copyright © 2018 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

HEADER_EXTRA_FILES = Makefile

PRE_COMMIT = headers.check-staged
COMMIT_MSG = git.commit-msg-log git.commit-msg-length git.commit-msg-empty git.commit-msg-prefix git.commit-msg-phrase git.commit-msg-casing

include .make/log.make
include .make/general.make
include .make/git.make
include .make/versions.make
include .make/headers.make
include .make/go/main.make
include .make/protos/main.make

is.database-drop:
	@$(log) "Dropping Cockroach database `echo $(NAME)`"
	cockroach sql --insecure --execute="DROP DATABASE $(NAME) CASCADE;"

ci.encrypt-variables:
	keybase encrypt -b -i ci/variables.yml -o ci/variables.yml.encrypted johanstokking htdvisser romeovs ericgo

ci.decrypt-variables:
	keybase decrypt -i ci/variables.yml.encrypted -o ci/variables.yml

messages:
	@$(GO) run ./pkg/errors/generate_messages.go --filename config/messages.json

dev-deps: go.dev-deps

dev-cert:
	go run $$(go env GOROOT)/src/crypto/tls/generate_cert.go -ca -host localhost

deps: go.deps

test: go.test

quality: go.quality

# Cache build artifacts (speeds up dev builds)
cache: go.install

# example binary
ttn-example: MAIN=./cmd/ttn-example/main.go
ttn-example: $(RELEASE_DIR)/ttn-example-$(GOOS)-$(GOARCH)

# stack binary
ttn-stack: MAIN=./cmd/ttn-stack/main.go
ttn-stack: $(RELEASE_DIR)/ttn-stack-$(GOOS)-$(GOARCH)

# All binaries
build-all: GO_FLAGS=-i -installsuffix ttn_prod
build-all: go.clean-build ttn-stack

# All supported platforms
build-all-platforms:
	GOOS=linux GOARCH=amd64 make build-all
	GOOS=linux GOARCH=386 make build-all
	GOOS=linux GOARCH=arm make build-all
	GOOS=linux GOARCH=arm64 make build-all
	GOOS=darwin GOARCH=amd64 make build-all
	GOOS=windows GOARCH=amd64 make build-all
	GOOS=windows GOARCH=386 make build-all
