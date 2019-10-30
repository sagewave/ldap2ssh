NAME=ldap2ssh
VERSION=0.3

BIN_DIR := $(CURDIR)/bin
SOURCE_FILES?=$$(go list ./... | grep -v /vendor/)

travis-ci: prepare build dist 

install:
	go install .

build: mod
	@rm -rf build/
	@$(BIN_DIR)/gox -ldflags "-X main.Version=$(VERSION)" \
		-osarch="darwin/amd64" \
		-osarch="linux/amd64" \
		-osarch="windows/amd64" \
		-output "build/{{.Dir}}_$(VERSION)_{{.OS}}_{{.Arch}}/$NAME" \
		${SOURCE_FILES}

dist:
	$(eval FILES := $(shell ls build))
	@rm -rf dist && mkdir dist
	@for f in $(FILES); do \
		(cd $(shell pwd)/build/$$f && tar cvzf ../../dist/$$f.tar.gz *); \
		echo $$f; \
	done

clean:
	@rm -rf build/
	@rm -rf dist/
	@rm -rf bin/

mod:
	GO111MODULE=on go mod download
	GO111MODULE=on go mod tidy

prepare:
	GOBIN=$(BIN_DIR) GO111MODULE=on go get github.com/buildkite/github-release
	GOBIN=$(BIN_DIR) GO111MODULE=on go get github.com/mitchellh/gox
	GOBIN=$(BIN_DIR) GO111MODULE=on go get github.com/axw/gocov/gocov
	GOBIN=$(BIN_DIR) GO111MODULE=on go get golang.org/x/tools/cmd/cover

.PHONY: default prepare mod build dist clean 