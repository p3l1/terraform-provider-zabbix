# ABOUTME: Build targets for the Zabbix Terraform provider.
# ABOUTME: Provides build, test, lint, and documentation generation commands.

BINARY_NAME=terraform-provider-zabbix

default: build

build:
	go build -v -o $(BINARY_NAME) .

install: build
	cp $(BINARY_NAME) $(shell go env GOPATH)/bin/

lint:
	golangci-lint run

generate:
	go generate ./...

fmt:
	gofmt -s -w .

test:
	go test -v -cover -timeout=120s -parallel=4 ./...

testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

.PHONY: default build install lint generate fmt test testacc
