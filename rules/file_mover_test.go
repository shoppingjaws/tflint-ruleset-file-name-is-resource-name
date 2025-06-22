package rules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformVariableFileNameRule_WithFix(t *testing.T) {
	tests := []struct {
		name          string
		files         map[string]string
		expectedIssues helper.Issues
		expectedFixes map[string]string
	}{
		{
			name: "variable in main.tf should be moved to variables.tf with fix",
			files: map[string]string{
				"main.tf": `
variable "example" {
  description = "An example variable"
  type        = string
  default     = "hello"
}

resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = var.example
}`,
			},
			expectedIssues: helper.Issues{
				{
					Rule:    NewTerraformVariableFileNameRule(),
					Message: "Variable block should be declared in variables.tf, not in main.tf",
				},
			},
			expectedFixes: map[string]string{
				"variables.tf": `
variable "example" {
  description = "An example variable"
  type        = string
  default     = "hello"
}
`,
				"main.tf": `

resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = var.example
}`,
			},
		},
	}

	rule := NewTerraformVariableFileNameRule()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tempDir := t.TempDir()
			
			for filename, content := range test.files {
				fullPath := filepath.Join(tempDir, filename)
				if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to write test file %s: %s", filename, err)
				}
			}

			runner := helper.TestRunner(t, test.files)

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			if len(runner.Issues) != len(test.expectedIssues) {
				t.Fatalf("Expected %d issues, got %d", len(test.expectedIssues), len(runner.Issues))
			}

			for i, issue := range runner.Issues {
				expectedIssue := test.expectedIssues[i]
				if issue.Message != expectedIssue.Message {
					t.Errorf("Expected message '%s', got '%s'", expectedIssue.Message, issue.Message)
				}
			}
		})
	}
}