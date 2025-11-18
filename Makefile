.PHONY: default
default: run

.ONESHELL:

# Configurable variables (can be overridden on the make command line)
WORKERS ?= 4
FILE ?=
DEBUG ?= false

.PHONY: run
run:
	@echo "Running with WORKERS=$(WORKERS) FILE=$(FILE) DEBUG=$(DEBUG)"
	if [ -z "$(FILE)" ]; then \
		go run . --workers=$(WORKERS) $(if $(filter true,$(DEBUG)),--debug,); \
	else \
		go run . --workers=$(WORKERS) $(if $(filter true,$(DEBUG)),--debug,) $(FILE); \
	fi

.PHONY: build
build:
	@echo "Building forky binary..."
	go build -o forky .
	@echo "Binary created: ./forky"

.PHONY: test
test:
	@echo "Running tests..."
	go test ./...

.PHONY: scan
scan:
	@echo "Running scanner mode (workers=$(WORKERS))"
	# Pass FILE as positional arg when set
	go run . --mode=scanning --workers=$(WORKERS) $(if $(FILE),$(FILE),)

.PHONY: parse
parse:
	@echo "Running parser mode (workers=$(WORKERS))"
	# Pass FILE as positional arg when set
	go run . --mode=parsing --workers=$(WORKERS) $(if $(FILE),$(FILE),)


.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -f forky
	go clean

.PHONY: help
help:
	@echo "Parameters (can be set on the make command line):"
	@echo "  WORKERS=<n>   Number of workers (default: 4)"
	@echo "  FILE=<path>   File to process (passed as positional argument)"
	@echo "  DEBUG=true    Enable debug output"
	@echo ""
	@echo "Usage:"
	@echo "  make          - Same as 'make run'"
	@echo "  make run      - Start REPL or run FILE if FILE is set"
	@echo "  make scan     - Run in scanning mode"
	@echo "  make parse    - Run in parsing mode"
	@echo ""
	@echo "Utilities:"
	@echo "  make build    - Build forky binary"
	@echo "  make test     - Run tests (go test ./...)"
	@echo "  make clean    - Remove build artifacts and run go clean"
	@echo "  make help     - Show this help"