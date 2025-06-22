default: build

.PHONY: test
test:
	go test ./...
	@# Clean up test artifacts
	@rm -f rules/variables.tf rules/output.tf rules/module.tf rules/provider.tf

.PHONY: test-keep-files
test-keep-files:
	go test ./...
	@echo "Test files preserved in rules/testdata/ for inspection"

.PHONY: build
build:
	go build

.PHONY: install
install: build
	mkdir -p ~/.tflint.d/plugins
	mv ./tflint-ruleset-file-name-is-resource-name ~/.tflint.d/plugins

.PHONY: clean
clean:
	rm -f ./tflint-ruleset-file-name-is-resource-name
	rm -f rules/variables.tf rules/output.tf rules/module.tf rules/provider.tf  
	rm -rf rules/testdata

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

.PHONY: test-scenario1
test-scenario1:
	./scripts/test-scenario.sh scenario1_variable_in_main

.PHONY: test-scenario2
test-scenario2:
	./scripts/test-scenario.sh scenario2_mixed_blocks_in_variables

.PHONY: test-scenario3
test-scenario3:
	./scripts/test-scenario.sh scenario3_existing_variables_tf

.PHONY: test-all-scenarios
test-all-scenarios: reset-scenarios
	@echo "Testing all scenarios..."
	@./scripts/test-scenario.sh scenario1_variable_in_main
	@echo ""
	@./scripts/test-scenario.sh scenario2_mixed_blocks_in_variables  
	@echo ""
	@./scripts/test-scenario.sh scenario3_existing_variables_tf
	@echo ""
	@echo "All scenarios tested!"
