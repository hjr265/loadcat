# Copyright 2021 Furqan Software Ltd. All rights reserved.

.PHONY: build test clean loadcatd

build: loadcatd

test:
	go test -mod=vendor -v ./...

clean:
	go clean -i ./...

loadcatd:
	GOOS=linux GOARCH=amd64 go build -mod=vendor -v -o loadcatd ./cmd/loadcatd
