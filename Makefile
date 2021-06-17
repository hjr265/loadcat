# Copyright 2021 Furqan Software Ltd. All rights reserved.

.PHONY: build test clean loadcatd

GO_BUILD := docker run --rm -ti -e GOOS -e GOARCH -v `pwd`:/loadcat -v /tmp/loadcat-go-build:/root/.cache/go-build -w /loadcat cardboard/golang:1.16 go build
ifeq ($(SKIP_DOCKER),true)
	GO_BUILD := go build
endif

build: loadcatd

test:
	go test -mod=vendor -v ./...

clean:
	go clean -i ./...

loadcatd:
	GOOS=linux GOARCH=amd64 $(GO_BUILD) -mod=vendor -v -o loadcatd ./cmd/loadcatd
