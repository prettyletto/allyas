BINARY_NAME := allyas
SRC_DIR := ./cmd/allyas
PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin
VERSION ?= 0.1.0-beta.1

.DEFAULT_GOAL := all

.PHONY: all build clean install run test fmt vet check

all: build

build:
	go build -trimpath -ldflags "-X main.version=$(VERSION)" -o $(BINARY_NAME) $(SRC_DIR)

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

check: fmt vet test

clean:
	rm -f $(BINARY_NAME)

install: build
	install -Dm755 $(BINARY_NAME) "$(DESTDIR)$(BINDIR)/$(BINARY_NAME)"

run: build
	./$(BINARY_NAME)
