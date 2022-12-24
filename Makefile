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
DEBUG_DIR := ${MKFILE_DIR}dbg
DOCKER_TAG := golayout

# Version
RELEASE?=0.0.1
ifndef GIT_COMMIT
  GIT_COMMIT := git-$(shell git rev-parse --short HEAD)
endif
GIT_REPO_INFO=$(shell git config --get remote.origin.url)

# Build Flags
GO_REL_LD_FLAGS= "-s -w -X ${PROJECT_NAME}/pkg/version.RELEASE=${RELEASE} -X ${PROJECT_NAME}/pkg/version.COMMIT=${GIT_COMMIT} -X ${PROJECT_NAME}/pkg/version.REPO=${GIT_REPO_INFO} -X ${PROJECT_NAME}/pkg/version.BUILDTIME=${DATETIME} -X ${PROJECT_NAME}/pkg/version.SERVICENAME=$@"
GO_DBG_LD_FLAGS= "-X ${PROJECT_NAME}/pkg/version.RELEASE=${RELEASE} -X ${PROJECT_NAME}/pkg/version.COMMIT=${GIT_COMMIT} -X ${PROJECT_NAME}/pkg/version.REPO=${GIT_REPO_INFO} -X ${PROJECT_NAME}/pkg/version.BUILDTIME=${DATETIME} -X ${PROJECT_NAME}/pkg/version.SERVICENAME=$@"
CGO_SWITCH := 0

# Binaries to build
BINS = demo

# Packages to build
PKGS = library

RELEASE = $(BINS:%=release-%)
DEBUG = $(BINS:%=debug-%)
RUN = $(BINS:%=run-%)
TEST = $(PKGS:%=test-%)

release: $(RELEASE)

$(RELEASE):
	cd ${MKFILE_DIR} && \
	CGO_ENABLED=${CGO_SWITCH} go build -v -trimpath -ldflags ${GO_REL_LD_FLAGS} \
	-o ${RELEASE_DIR}/$(@:release-%=%) ${MKFILE_DIR}cmd/$(@:release-%=%)/

run: $(RELEASE) $(RUN)

$(RUN):
	cd ${MKFILE_DIR} && \
	${RELEASE_DIR}/$(@:run-%=%) ${MKFILE_DIR}data/books.json

debug: $(DEBUG)

$(DEBUG):
	cd ${MKFILE_DIR} && \
	CGO_ENABLED=${CGO_SWITCH} go build -v -trimpath -gcflags "-N" \
	-o ${DEBUG_DIR}/$(@:debug-%=%) ${MKFILE_DIR}cmd/$(@:debug-%=%)/ && \
	gdb ${DEBUG_DIR}/$(@:debug-%=%)

test: $(TEST)

$(TEST):
	cd ${MKFILE_DIR} && \
	go test ${PROJECT_NAME}/pkg/$(@:test-%=%)

clean:
	@rm -f ${RELEASE_DIR}/*
	@rm -f ${DEBUG_DIR}/*

.PHONY: release
