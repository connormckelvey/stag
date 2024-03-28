SRC=$(shell find . -name "*.go")

.PHONY: fmt lint test examples deps

default: all

all: fmt lint test examples

fmt:
	$(info * [checking formatting] **************************************)
	@test -z $(shell gofmt -l $(SRC)) || (gofmt -d $(SRC); exit 1)

lint:
	$(info * [running lint tools] ***************************************)
	golangci-lint run -v

test: deps
	$(info * [running tests] ********************************************)
	go test -v $(shell go list ./... | grep -v /examples$)

examples: deps
	$(info * [running tests] ********************************************)
	go test -v ./examples/...

deps:
	$(info * {downloading dependencies} *********************************)
	go get -v ./...

