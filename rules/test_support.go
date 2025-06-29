package rules

import (
	"os"
	"path/filepath"
	"testing"
)

// TestScenario represents a test scenario for file operations
type TestScenario struct {
	Name          string
	Files         map[string]string
	ExpectedFiles map[string]string
}

// CreateTestFiles creates test files in a specified directory
func CreateTestFiles(t *testing.T, baseDir string, scenarios map[string]map[string]string) {
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
}


// GetTestScenarios returns predefined test scenarios
func GetTestScenarios() map[string]map[string]string {
	return map[string]map[string]string{
		"scenario1_blocks_in_wrong_files": {
			"main.tf": `variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t2.micro"
}

output "instance_id" {
  description = "The ID of the EC2 instance"
  value       = aws_instance.example.id
}

locals {
  common_tags = {
    Environment = "dev"
    Project     = "example"
  }
}

resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = var.instance_type
  tags          = local.common_tags
}`,
		},
		"scenario2_mixed_blocks_in_dedicated_files": {
			"variables.tf": `variable "valid_var" {
  type = string
}

resource "aws_instance" "should_not_be_here" {
  ami = "ami-12345678"
}

output "should_not_be_here" {
  value = "test"
}`,
			"outputs.tf": `output "valid_output" {
  value = "test"
}

variable "should_not_be_here" {
  type = string
}

locals {
  should_not_be_here = "test"
}`,
			"locals.tf": `locals {
  valid_local = "test"
}

output "should_not_be_here" {
  value = "test"
}

variable "should_not_be_here" {
  type = string
}`,
		},
		"scenario3_existing_files_with_new_blocks": {
			"variables.tf": `variable "existing" {
  type = string
}`,
			"outputs.tf": `output "existing_output" {
  value = "existing"
}`,
			"locals.tf": `locals {
  existing_local = "existing"
}`,
			"main.tf": `variable "new_var" {
  type = number
}

output "new_output" {
  value = "new"
}

locals {
  new_local = "new"
}

resource "aws_instance" "example" {
  ami = "ami-12345678"
}`,
		},
		"scenario4_resource_specific_files": {
			"aws_instance.tf": `resource "aws_instance" "correct" {
  ami = "ami-12345678"
}

resource "aws_s3_bucket" "wrong_type" {
  bucket = "wrong-resource-type"
}

variable "should_not_be_here" {
  type = string
}`,
			"aws_s3_bucket.tf": `resource "aws_s3_bucket" "data" {
  bucket = "my-data-bucket"
}`,
			"main.tf": `resource "aws_iam_role" "lambda_role" {
  name = "lambda-execution-role"
}

resource "aws_s3_bucket" "wrong_place" {
  bucket = "should-be-in-s3-file"
}`,
		},
	}
}
