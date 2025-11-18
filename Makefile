.PHONY: default
default: run

.ONESHELL:

# Configurable variables (can be overridden on command line)
WORKERS ?= 4
FILE ?=
DEBUG ?= false
INJECT ?= false

define CHECK_INJECT
	if [ "$(INJECT)" = "true" ]; then \
		if [ -z "$(FILE)" ]; then \
			echo "Error: FILE must be set when INJECT=true" >&2; exit 1; \
		fi; \
	fi
endef


GO_ARGS = \
	--workers=$(WORKERS) \
	$(if $(filter true,$(DEBUG)),--debug,) \
	$(if $(filter true,$(INJECT)),--inject,) \
	$(if $(FILE),$(FILE),)

.PHONY: run
run:
	@$(CHECK_INJECT)
	@go run . $(GO_ARGS)

.PHONY: scan
scan:
	@$(CHECK_INJECT)
	@go run . --mode=scanning $(GO_ARGS)

.PHONY: parse
parse:
	@$(CHECK_INJECT)
	@go run . --mode=parsing $(GO_ARGS)

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

.PHONY: inject
inject:
	$(MAKE) run INJECT=true

.PHONY: help
help:
	@echo "Parameters:"
	@echo "  WORKERS=<n>   Number of workers for parallel scanning (default: 4)"
	@echo "  FILE=<path>   Input file to process"
	@echo "  INJECT=true   Enable inject mode (requires FILE)"
	@echo "  DEBUG=true    Enable debug output"
	@echo ""
	@echo "Commands:"
	@echo "  make run      - Execute in normal mode or process FILE"
	@echo "  make scan     - Run lexical analysis on FILE"
	@echo "  make parse    - Run parsing on FILE"
	@echo "  make inject   - Inject FILE and continue in REPL"
	@echo "  make build    - Build the forky binary"
	@echo "  make test     - Run Go tests"
	@echo "  make clean    - Remove build artifacts"
	@echo "  make help     - Show this help"
