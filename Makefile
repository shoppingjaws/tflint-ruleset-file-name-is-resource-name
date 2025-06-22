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
	@echo "Test scenarios created in rules/testdata/manual_verification/"
	@echo "Use these directories to manually test tflint --fix functionality"

.PHONY: demo-file-ops
demo-file-ops:
	go test ./rules -run Test_Manual_FileOperationsDemo -v
	@echo "File operation demo completed. Check rules/testdata/demo/"
