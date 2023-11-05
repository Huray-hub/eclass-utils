.DEFAULT-GOAL := help

# Define variables
EXECUTABLE_NAME := assignments
BUILD_DIR := .

# Help target
.PHONY: help
help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  help                                   Display this help message"
	@echo
	@echo "  install-assignments [PREFIX=/usr/bin]  Install the assignments executable (default: /usr/bin)"
	@echo "  install-assignments-local              Install the assignments executable to ~/.local/bin/"
	@echo "  install-assignments-gopath             Install the assignments executable to GOPATH/bin/"
	@echo "  build-assignments                      Build the assignments executable"
	@echo "  run-assignments                        Run the assignments executable"

# Build target
.PHONY: build-assignments
build-assignments:
	go build -v -o $(BUILD_DIR)/$(EXECUTABLE_NAME) ./assignment/cmd/assignments/main.go

# Run target
.PHONY: run-assignments
run-assignments:
	go run ./assignment/cmd/assignments/main.go

# Install target
.PHONY: install-assignments
install-assignments: build-assignments
	install -D $(BUILD_DIR)/$(EXECUTABLE_NAME) $(PREFIX)/$(EXECUTABLE_NAME)

# Install local target
.PHONY: install-assignments-local
install-assignments-local: build-assignments
	install -D $(BUILD_DIR)/$(EXECUTABLE_NAME) ~/.local/bin/$(EXECUTABLE_NAME)

# Install GOPATH target
.PHONY: install-assignments-gopath
install-assignments-gopath: build-assignments
	install -D $(BUILD_DIR)/$(EXECUTABLE_NAME) $(GOPATH)/bin/$(EXECUTABLE_NAME)
