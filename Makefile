VERSION=0.1.0
GITCOMMIT?=$(shell git describe --dirty --always)
LDFLAGS=-ldflags "-w -s -X main.version=${VERSION} -X main.commit=${GITCOMMIT}"

all: check-dns-multi

.PHONY: check-dns-multi

check-dns-multi: cmd/check-dns-multi/*.go
	go build $(LDFLAGS) -o check-dns-multi ./cmd/check-dns-multi/

linux: cmd/check-dns-multi/*.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o check-dns-multi ./cmd/check-dns-multi/

check:
	go test -v ./...
	go test -race ./...


