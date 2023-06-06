help:
	@echo "Usage: make [target]"
	@echo "Targets:"
	@echo "  help                                      Display this help message"
	@echo
	@echo "  install-assignments                       Install the assignments executable to /usr/bin/"
	@echo "  install-assignments-local                 Install the assignments executable to ~/.local/bin/"
	@echo "  install-assignments-gopath                Install the assignments executable to GOPATH/bin/"
	@echo "  build-assignments                         Build the assignments executable"
	@echo "  run-assignments                           Run the assignments executable"

install-assignments:
	go build -v -o assignments ./assignment/cmd/main.go && sudo mv assignments /usr/bin/

install-assignments-local:
	go build -v -o assignments ./assignment/cmd/main.go && mv assignments ~/.local/bin/

install-assignments-gopath:
	go build -v -o assignments ./assignment/cmd/main.go && mv assignments $(GOPATH)/bin/

build-assignments:
	go build -v -o assignments ./assignment/cmd/main.go

run-assignments:
	go run ./assignment/cmd/main.go
