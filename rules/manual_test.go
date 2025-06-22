package rules

import (
	"testing"
)

// Test_Manual_CreateTestFiles creates test templates for repeatable testing
// This test is designed to help developers verify the --fix behavior manually
// Run: go test -run Test_Manual_CreateTestFiles
func Test_Manual_CreateTestFiles(t *testing.T) {
	baseDir := "testdata/templates"
	scenarios := GetTestScenarios()

	CreateTestFiles(t, baseDir, scenarios)

	t.Log("Test templates created in rules/testdata/templates/")
	t.Log("To test with automatic reset:")
	t.Log("1. make reset-scenarios")
	t.Log("2. make test-scenario1  # or test-scenario2, test-scenario3")
	t.Log("3. make test-all-scenarios  # to test all at once")
	t.Log("")
	t.Log("For manual testing:")
	t.Log("1. make reset-scenarios")
	t.Log("2. cd rules/testdata/working/scenario1_variable_in_main")
	t.Log("3. tflint --fix")
	t.Log("4. make reset-scenarios  # to reset for next test")
}

// Test_Manual_FileOperationsDemo demonstrates the block manager functionality
func Test_Manual_FileOperationsDemo(t *testing.T) {
	demoDir := "testdata/demo"
	CreateDemoFiles(t, demoDir)
}