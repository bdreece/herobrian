# herobrian v1.2.0
# Copyright (C) 2024 Brian Reece

NPROCS = $(shell grep -c 'processor' /proc/cpuinfo)
MAKEFLAGS += -j$(NPROCS)

all: build

# ================ #
# SCRIPTS
# ================ #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage: make [SCRIPT]'
	@echo 'Scripts:'
	@sed -n 's/^## //p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/  /'

## restore: restore dependencies
.PHONY: restore restore/go restore/npm
restore: restore/go restore/npm

restore/go: generate
	go mod download
	go mod verify

restore/npm:
	npm ci

## build: build the application
.PHONY: build build/go build/npm
build: build/go build/npm

build/go: restore/go
	@mkdir -p bin
	go build -v -o $(abspath ./bin) ./...
	
build/npm: restore/npm generate
	npm run -ws --if-present build

## generate: run go source generator
.PHONY: generate
generate:
	go generate -v ./...

## run: run the application
.PHONY: run
run:
	go run -v $(abspath ./cmd/herobrian)

## test: run all tests
.PHONY: test
test:
	go test -v ./...
	npm run -ws --if-present test

## watch: run the application with live-reload
.PHONY: watch
watch:
	air -c configs/air.toml

## install: install the application
.PHONY: install
install:
	install -m 0644 -Dt /usr/share/herobrian/app web/app/dist 
	install -m 0644 -Dt /usr/share/herobrian/static web/static
	install -m 0644 -Dt /usr/share/herobrian/templates web/templates
	install -m 0644 -Dt /etc/herobrian configs/schema.sql
	install -m 0644 -Dt /etc/herobrian configs/settings.production.yml
	install -m 0755 -t /usr/bin bin/herobrian
