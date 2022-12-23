# Copyright David Huseby. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

SHELL:=/bin/sh
PROJECT_NAME := library
DATETIME = $(shell date '+%Y%m%d_%H%M%S')
PROTOS = `ls proto`

# Path Related
MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR := $(dir $(MKFILE_PATH))
RELEASE_DIR := ${MKFILE_DIR}bin
DOCKER_TAG := golayout

# Version
RELEASE?=0.0.1
ifndef GIT_COMMIT
  GIT_COMMIT := git-$(shell git rev-parse --short HEAD)
endif
GIT_REPO_INFO=$(shell git config --get remote.origin.url)

# Build Flags
GO_LD_FLAGS= "-s -w -X ${PROJECT_NAME}/pkg/version.RELEASE=${RELEASE} -X ${PROJECT_NAME}/pkg/version.COMMIT=${GIT_COMMIT} -X ${PROJECT_NAME}/pkg/version.REPO=${GIT_REPO_INFO} -X ${PROJECT_NAME}/pkg/version.BUILDTIME=${DATETIME} -X ${PROJECT_NAME}/pkg/version.SERVICENAME=$@"
CGO_SWITCH := 0

# Binaries to build
BINS = demo

# Packages to build
PKGS = library

BUILDBINS = $(BINS:%=build-%)
RUNBINS = $(BINS:%=run-%)
TESTPKGS = $(PKGS:%=test-%)

build: $(BUILDBINS)

$(BUILDBINS):
	cd ${MKFILE_DIR} && \
	CGO_ENABLED=${CGO_SWITCH} go build -v -trimpath -ldflags ${GO_LD_FLAGS} \
	-o ${RELEASE_DIR}/$(@:build-%=%) ${MKFILE_DIR}cmd/$(@:build-%=%)/

run: $(BUILDBINS) $(RUNBINS)

$(RUNBINS):
	cd ${MKFILE_DIR} && \
	${RELEASE_DIR}/$(@:run-%=%)

test: $(TESTPKGS)

$(TESTPKGS):
	cd ${MKFILE_DIR} && \
	go test ${PROJECT_NAME}/pkg/$(@:test-%=%)

clean:
	@rm -f ${RELEASE_DIR}/*

.PHONY: build
