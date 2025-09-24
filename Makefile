.PHONY: default
default: repl

.PHONY: repl
repl:
	@echo "Starting Forky REPL..."
	go run .

.PHONY: build
build:
	@echo "Building forky binary..."
	go build -o forky .
	@echo "Binary created: ./forky"

.PHONY: test
test:
	@echo "Running tests..."
	go test ./...

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -f forky
	go clean

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make          - Start REPL (default)"
	@echo "  make repl     - Start REPL"
	@echo "  make build    - Build forky binary"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make help     - Show this help"