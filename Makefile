default: build

.PHONY: test
test:
	go test ./...
	@# Clean up test artifacts
	@rm -f rules/variables.tf rules/output.tf rules/module.tf rules/provider.tf

.PHONY: build
build:
	go build

.PHONY: install
install: build
	mkdir -p ~/.tflint.d/plugins
	mv ./tflint-ruleset-file-name-is-resource-name ~/.tflint.d/plugins

.PHONY: clean-testdata
clean-testdata:
	rm -rf rules/testdata
	@echo "Test data directories cleaned"

.PHONY: create-test-scenarios
create-test-scenarios:
	go test ./rules -run Test_Manual_CreateTestFiles
	@echo "Test templates created in rules/testdata/templates/"
	@echo "Use 'make reset-scenarios' and 'make test-scenario1' for repeatable testing"

.PHONY: demo-file-ops
demo-file-ops:
	go test ./rules -run Test_Manual_FileOperationsDemo -v
	@echo "File operation demo completed. Check rules/testdata/demo/"

.PHONY: reset-scenarios
reset-scenarios:
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
