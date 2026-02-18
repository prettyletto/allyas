
BINARY_NAME = allyas
SRC_DIR = cmd/allyas

.DEFAULT_GOAL := all

.PHONY: all build clean install run help

all: build

build:
	@echo "Building the $(BINARY_NAME) binary..."
	go build -o $(BINARY_NAME) ./$(SRC_DIR)

clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)

install: build
	@if [ -f $(BINARY_NAME) ]; then \
		echo "Installing $(BINARY_NAME) globally..."; \
		sudo mv $(BINARY_NAME) /usr/local/bin/; \
	else \
		echo "$(BINARY_NAME) not found! Build it first."; \
		exit 1; \
	fi

run:
	@if [ -f $(BINARY_NAME) ]; then \
		echo "Running $(BINARY_NAME)..."; \
		./$(BINARY_NAME); \
	else \
		echo "$(BINARY_NAME) not found! Build it first."; \
		exit 1; \
	fi

help:
	@echo "Makefile commands:"
	@echo "  all        - Build the binary (default target)"
	@echo "  build      - Build the binary from the source"
	@echo "  clean      - Clean build artifacts"
	@echo "  install    - Install the binary globally"
	@echo "  run        - Run the application"
	@echo "  help       - Show this help message"

