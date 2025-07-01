# Default target
default: build

# =============================================================================
# Development Commands
# =============================================================================

.PHONY: build
build:
	@echo "Building TFLint ruleset..."
	go build -o tflint-ruleset-file-name-is-resource-name

.PHONY: test
test:
	@echo "Running tests..."
	go test ./...

.PHONY: test-verbose
test-verbose:
	@echo "Running tests with verbose output..."
	go test -v ./...

.PHONY: install
install: build
	@echo "Installing plugin to ~/.tflint.d/plugins/"
	mkdir -p ~/.tflint.d/plugins
	cp ./tflint-ruleset-file-name-is-resource-name ~/.tflint.d/plugins/
	@echo "Plugin installed successfully"

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -f ./tflint-ruleset-file-name-is-resource-name
	@echo "Build artifacts cleaned"

# =============================================================================
# Testing Infrastructure
# =============================================================================

.PHONY: clean-testdata
clean-testdata:
	@echo "Cleaning test data..."
	rm -rf rules/testdata/working
	@echo "Test data cleaned"

.PHONY: reset-scenarios
reset-scenarios:
	@echo "Resetting test scenarios..."
	./scripts/reset-scenarios.sh

# Dynamic scenario testing - usage: make test-scenario SCENARIO=scenario_name
.PHONY: test-scenario
test-scenario:
	@if [ -z "$(SCENARIO)" ]; then \
		echo "Usage: make test-scenario SCENARIO=<scenario_name>"; \
		echo ""; \
		echo "Available scenarios:"; \
		if [ -d "rules/testdata/templates" ]; then \
			for scenario in rules/testdata/templates/*; do \
				if [ -d "$$scenario" ]; then \
					echo "  - $$(basename "$$scenario")"; \
				fi; \
			done; \
		fi; \
		exit 1; \
	fi
	./scripts/test-scenario.sh $(SCENARIO)

.PHONY: test-all-scenarios
test-all-scenarios: reset-scenarios
	@echo "Testing all scenarios..."
	@for scenario in rules/testdata/templates/*; do \
		if [ -d "$$scenario" ]; then \
			scenario_name=$$(basename "$$scenario"); \
			echo ""; \
			./scripts/test-scenario.sh "$$scenario_name"; \
		fi; \
	done
	@echo ""
	@echo "All scenarios tested!"

.PHONY: test-scenarios
test-scenarios: reset-scenarios
	@echo "Running scenario tests and comparing with expected results..."
	@for scenario in rules/testdata/templates/*; do \
		if [ -d "$$scenario" ]; then \
			scenario_name=$$(basename "$$scenario"); \
			echo "Testing scenario: $$scenario_name"; \
			SKIP_RESET=1 ./scripts/test-scenario.sh "$$scenario_name"; \
			if [ -d "rules/testdata/answer/$$scenario_name" ]; then \
				echo "Comparing results for $$scenario_name..."; \
				if diff -r "rules/testdata/working/$$scenario_name" "rules/testdata/answer/$$scenario_name" > /dev/null 2>&1; then \
					echo "✓ $$scenario_name passed"; \
				else \
					echo "✗ $$scenario_name failed - differences found:"; \
					diff -r "rules/testdata/working/$$scenario_name" "rules/testdata/answer/$$scenario_name"; \
					exit 1; \
				fi; \
			else \
				echo "⚠ No answer directory found for $$scenario_name"; \
			fi; \
		fi; \
	done
	@echo ""
	@echo "All scenario tests passed!"

.PHONY: validate-test-system
validate-test-system:
	@echo "Validating that the test system correctly detects failures..."
	./scripts/validate-test-system.sh

# =============================================================================
# Quick Testing Shortcuts
# =============================================================================

.PHONY: test-scenario1
test-scenario1:
	make test-scenario SCENARIO=scenario1_blocks_in_wrong_files

.PHONY: test-scenario2
test-scenario2:
	make test-scenario SCENARIO=scenario2_custom_file_names


# =============================================================================
# Development Utilities
# =============================================================================

.PHONY: lint
lint:
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found, running basic go vet..."; \
		go vet ./...; \
	fi

.PHONY: fmt
fmt:
	@echo "Formatting Go code..."
	go fmt ./...

.PHONY: mod-tidy
mod-tidy:
	@echo "Tidying Go modules..."
	go mod tidy

# =============================================================================
# Help
# =============================================================================

.PHONY: help
help:
	@echo "TFLint Ruleset File Name Development"
	@echo ""
	@echo "Development Commands:"
	@echo "  build              Build the plugin binary"
	@echo "  test               Run all tests"
	@echo "  test-verbose       Run tests with verbose output"
	@echo "  install            Build and install plugin to ~/.tflint.d/plugins/"
	@echo "  clean              Clean build artifacts"
	@echo "  lint               Run linters (golangci-lint or go vet)"
	@echo "  fmt                Format Go code"
	@echo "  mod-tidy           Tidy Go modules"
	@echo ""
	@echo "Testing Infrastructure:"
	@echo "  reset-scenarios    Reset test scenarios to initial state"
	@echo "  test-scenario      Test specific scenario (requires SCENARIO=name)"
	@echo "  test-all-scenarios Test all available scenarios"
	@echo "  test-scenarios     Test all scenarios and compare with expected results"
	@echo "  validate-test-system Validate that test system correctly detects failures"
	@echo "  clean-testdata     Clean test data directories"
	@echo ""
	@echo "Quick Testing:"
	@echo "  test-scenario1     Test scenario1_blocks_in_wrong_files"
	@echo "  test-scenario2     Test scenario2_custom_file_names"
	@echo ""
	@echo "Available scenarios:"
	@if [ -d "rules/testdata/templates" ]; then \
		for scenario in rules/testdata/templates/*; do \
			if [ -d "$$scenario" ]; then \
				echo "  - $$(basename "$$scenario")"; \
			fi; \
		done; \
	fi