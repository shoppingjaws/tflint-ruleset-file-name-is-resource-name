package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformVariableFileNameRule(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected helper.Issues
	}{
		{
			name: "variable in variables.tf - valid",
			files: map[string]string{
				"variables.tf": `
variable "example" {
  description = "An example variable"
  type        = string
  default     = "hello"
}`,
			},
			expected: helper.Issues{},
		},
		{
			name: "variable in main.tf - invalid",
			files: map[string]string{
				"main.tf": `
variable "example" {
  description = "An example variable"
  type        = string
  default     = "hello"
}`,
			},
			expected: helper.Issues{
				{
					Rule:    NewTerraformVariableFileNameRule(),
					Message: "Variable block should be declared in variables.tf, not in main.tf",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 19},
					},
				},
			},
		},
		{
			name: "multiple variables in wrong file - invalid",
			files: map[string]string{
				"config.tf": `
variable "first" {
  type = string
}

variable "second" {
  type = number
}`,
			},
			expected: helper.Issues{
				{
					Rule:    NewTerraformVariableFileNameRule(),
					Message: "Variable block should be declared in variables.tf, not in config.tf",
					Range: hcl.Range{
						Filename: "config.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 17},
					},
				},
				{
					Rule:    NewTerraformVariableFileNameRule(),
					Message: "Variable block should be declared in variables.tf, not in config.tf",
					Range: hcl.Range{
						Filename: "config.tf",
						Start:    hcl.Pos{Line: 6, Column: 1},
						End:      hcl.Pos{Line: 6, Column: 18},
					},
				},
			},
		},
		{
			name: "no variables - valid",
			files: map[string]string{
				"main.tf": `
resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = "t2.micro"
}`,
			},
			expected: helper.Issues{},
		},
		{
			name: "mixed content with variables in correct file - valid",
			files: map[string]string{
				"variables.tf": `
variable "instance_type" {
  type    = string
  default = "t2.micro"
}`,
				"main.tf": `
resource "aws_instance" "example" {
  ami           = "ami-12345678"
  instance_type = var.instance_type
}`,
			},
			expected: helper.Issues{},
		},
		{
			name: "resource block in variables.tf - invalid",
			files: map[string]string{
				"variables.tf": `
variable "instance_type" {
  type = string
}

resource "aws_instance" "example" {
  ami = "ami-12345678"
}`,
			},
			expected: helper.Issues{
				{
					Rule:    NewTerraformVariableFileNameRule(),
					Message: "Only variable blocks should be declared in variables.tf, found resource block",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 6, Column: 1},
						End:      hcl.Pos{Line: 6, Column: 34},
					},
				},
			},
		},
		{
			name: "output block in variables.tf - invalid",
			files: map[string]string{
				"variables.tf": `
output "instance_id" {
  value = "test"
}`,
			},
			expected: helper.Issues{
				{
					Rule:    NewTerraformVariableFileNameRule(),
					Message: "Only variable blocks should be declared in variables.tf, found output block",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 21},
					},
				},
			},
		},
	}

	rule := NewTerraformVariableFileNameRule()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			runner := helper.TestRunner(t, test.files)

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, test.expected, runner.Issues)
		})
	}
}
