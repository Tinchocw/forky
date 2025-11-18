.PHONY: default
default: run

.ONESHELL:

# Configurable variables (can be overridden on the make command line)
WORKERS ?= 4
MODE ?= repl
FILE ?=
DEBUG ?= false

.PHONY: repl
repl:
	@echo "Starting Forky REPL..."
	go run .

.PHONY: run
run:
	@echo "Running with MODE=$(MODE) WORKERS=$(WORKERS) FILE=$(FILE) DEBUG=$(DEBUG)"
	# If FILE is set, pass it as a positional argument; otherwise run REPL
	if [ -z "$(FILE)" ]; then \
		echo "Starting REPL (workers=$(WORKERS))"; \
		go run . --workers=$(WORKERS) $(if $(filter true,$(DEBUG)),--debug,); \
	else \
		# Include MODE if provided, pass FILE as positional arg
		go run . $(if $(MODE),--mode=$(MODE),) --workers=$(WORKERS) $(if $(filter true,$(DEBUG)),--debug,) $(FILE); \
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
	go run . --scanning --workers=$(WORKERS) $(if $(FILE),$(FILE),)

.PHONY: parse
parse:
	@echo "Running parser mode (workers=$(WORKERS))"
	# Pass FILE as positional arg when set
	go run . --parsing --workers=$(WORKERS) $(if $(FILE),$(FILE),)


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
	@echo "  MODE=<mode>   Mode to run (repl|scanning|parsing). Default: repl"
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