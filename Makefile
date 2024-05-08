GO = golang:1.21 go
LNXOS = alpine:3.7
BINARY = ./bin/service
MAIN = ./main.go
DB_NET = db-net
DOCKER = docker run --rm -v $(shell pwd):/svc -w /svc
OS ?=
PRIVATE_KEY ?=

## display help message
.PHONY: help
help:
	@echo ''
	@echo 'Management Commands for Go Reference Service:'
	@echo
	@echo 'Usage:'
	@echo '  ## Develop / Test Commands ##'
	@echo '  clean              Run clean up.'
	@echo '  build              Build the Service Binary'
	@echo '  fmt                Run code formatter.'
	@echo '  check              Run static code analysis (lint).'
	@echo '  generate-api       Generating API layer Boilerplate Code'
	@echo ''

## Clean up
.PHONY: clean
clean:
	@echo '==> Cleaning...'
	${DOCKER} ${LNXOS} rm -f coverage.out report.json function.zip
	${DOCKER} ${LNXOS} rm -Rf vendor
	${DOCKER} ${LNXOS} rm -f ${BINARY}

## Build the Service
.PHONY: build
build:
ifeq ($(strip $(OS)),osx)
	@echo '==> Building Service Binary for MacOS'
	GOOS=darwin GOARCH=amd64 go build -o ${BINARY} ${MAIN}
endif
ifeq ($(strip $(OS)),linux)
	@echo '==> Building Service Binary for linux'
	GOOS=linux GOARCH=amd64 go build -o ${BINARY} ${MAIN}
endif
ifeq ($(strip $(OS)),windows)
	@echo '==> Building Service Binary for Windows'
	GOOS=windows GOARCH=amd64 go build -o ${BINARY} ${MAIN}
endif
ifeq ($(strip $(OS)),)
	@echo ''
	@echo 'Please provide a value for os:'
	@echo ''
	@echo 'OS=osx       For Mac Operating System'
	@echo 'OS=linux     For Linux Operating System'
	@echo 'OS=windows   For Windows Operating System'
	@echo ''
endif

## Run code formatter
.PHONY: fmt
fmt:
	@echo '==> Code formatting...'
	${DOCKER} cytopia/goimports:latest-0.3 -l -w .

.PHONY: check
check: fmt install-linter
	@echo '==> Code check...'
	bin/golangci-lint run -v -c .golangci.yml
	@$(MAKE) validate-api

## validate api
.PHONY: validate-api
validate-api:
	@echo "==> Validate OpenAPI Spec"
	${DOCKER} openapitools/openapi-generator-cli:v7.2.0 validate -i api.yaml

## generate api
.PHONY: generate-api
generate-api:
	@echo "==> Generate API(s)"
	@if [ -d api/handlers ]; then \
  		echo "==> handlers directory exist" ; \
	else \
	  	echo "==> creating handlers directory" ; \
  		${DOCKER} ${LNXOS} mkdir api/handlers ; \
	fi
	${DOCKER} openapitools/openapi-generator-cli:v7.2.0 generate \
		-i api.yaml \
		-g go-server \
		-t templates \
		--additional-properties=featureCORS=true \
		-o api/gen
	${DOCKER} ${LNXOS} rm -rf api/gen/api api/gen/.openapi-generator
	-${DOCKER} ${LNXOS} mv api/gen/go/*_service.go api/handlers
	${DOCKER} ${LNXOS} mv api/gen/go/* api/gen
	${DOCKER} ${LNXOS} rm -rf api/gen/go
	@${MAKE} import-package-update
	@${MAKE} fmt

## update imports and package names
.PHONY: import-package-update
import-package-update:
	@echo "==> updating imports and package names"
	@for f in $(shell ls api/handlers); do \
  		if ! grep -q "$${f}" api/gen/.openapi-generator-ignore; then \
			${DOCKER} ${LNXOS} sed -i 's/package openapi/package handlers/' api/handlers/$${f} ; \
			${DOCKER} ${LNXOS} echo 'go/'$${f} >> api/gen/.openapi-generator-ignore ; \
			${DOCKER} ${LNXOS} sed -i 's/package openapi/package handlers/' api/handlers/$${f} ; \
			${DOCKER} ${LNXOS} sed -i 's~import (~import (. "github.com/alighm/sample-service/api/gen"~' api/handlers/$${f} ; \
		fi \
	done

.PHONY: install-linter
install-linter:
	@echo '==> Install Linter...'
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.55.2
