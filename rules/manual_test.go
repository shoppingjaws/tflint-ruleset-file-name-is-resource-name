package rules

import (
	"os"
	"path/filepath"
	"testing"
)

// Test_Manual_CreateTestFiles creates test files in testdata for manual verification
// This test is designed to help developers verify the --fix behavior manually
// Run: go test -run Test_Manual_CreateTestFiles
func Test_Manual_CreateTestFiles(t *testing.T) {
	baseDir := "testdata/manual_verification"
	
	scenarios := map[string]map[string]string{
		"scenario1_variable_in_main": {
			"main.tf": `variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t2.micro"
}

variable "region" {
  description = "AWS region"
  type        = string
  default     = "us-west-2"
}

resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = var.instance_type
}`,
		},
		"scenario2_mixed_blocks_in_variables": {
			"variables.tf": `variable "valid_var" {
  type = string
}

resource "aws_instance" "should_not_be_here" {
  ami = "ami-12345678"
}

output "should_not_be_here" {
  value = "test"
}`,
		},
		"scenario3_existing_variables_tf": {
			"variables.tf": `variable "existing" {
  type = string
}`,
			"main.tf": `variable "new_var" {
  type = number
}

resource "aws_instance" "example" {
  ami = "ami-12345678"
}`,
		},
	}

	for scenarioName, files := range scenarios {
		scenarioDir := filepath.Join(baseDir, scenarioName)
		if err := os.MkdirAll(scenarioDir, 0755); err != nil {
			t.Fatalf("Failed to create scenario directory %s: %s", scenarioDir, err)
		}

		for filename, content := range files {
			fullPath := filepath.Join(scenarioDir, filename)
			if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
				t.Fatalf("Failed to write file %s: %s", fullPath, err)
			}
		}

		t.Logf("Created scenario: %s in %s", scenarioName, scenarioDir)
	}

	t.Log("Manual test files created in rules/testdata/manual_verification/")
	t.Log("To test manually:")
	t.Log("1. cd rules/testdata/manual_verification/scenario1_variable_in_main")
	t.Log("2. tflint --init")
	t.Log("3. tflint")
	t.Log("4. tflint --fix")
}

// Test_Manual_FileOperationsDemo demonstrates the file mover functionality
func Test_Manual_FileOperationsDemo(t *testing.T) {
	demoDir := "testdata/demo"
	if err := os.MkdirAll(demoDir, 0755); err != nil {
		t.Fatalf("Failed to create demo directory: %s", err)
	}

	// Demo 1: Create variables.tf from scratch
	t.Run("demo_create_variables_tf", func(t *testing.T) {
		targetFile := filepath.Join(demoDir, "variables.tf")
		
		// Remove if exists
		os.Remove(targetFile)
		
		content := `variable "demo_var" {
  description = "A demo variable"
  type        = string
  default     = "demo"
}`

		blockMover := &BlockMover{}
		if err := blockMover.appendToTargetFile(targetFile, content); err != nil {
			t.Fatalf("Failed to create variables.tf: %s", err)
		}

		t.Logf("Created %s", targetFile)
	})

	// Demo 2: Append to existing variables.tf
	t.Run("demo_append_to_variables_tf", func(t *testing.T) {
		targetFile := filepath.Join(demoDir, "variables.tf")
		
		additionalContent := `
variable "another_var" {
  type = number
  default = 42
}`

		blockMover := &BlockMover{}
		if err := blockMover.appendToTargetFile(targetFile, additionalContent); err != nil {
			t.Fatalf("Failed to append to variables.tf: %s", err)
		}

		// Read and display final content
		finalContent, err := os.ReadFile(targetFile)
		if err != nil {
			t.Fatalf("Failed to read final content: %s", err)
		}

		t.Logf("Final content of %s:\n%s", targetFile, string(finalContent))
	})

	t.Logf("Demo files created in %s", demoDir)
}