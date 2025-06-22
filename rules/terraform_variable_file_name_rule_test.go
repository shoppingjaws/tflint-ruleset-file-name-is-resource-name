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
					Message: "Output block should be declared in outputs.tf, not in variables.tf",
					Range: hcl.Range{
						Filename: "variables.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 21},
					},
				},
			},
		},
		{
			name: "output block in main.tf - invalid",
			files: map[string]string{
				"main.tf": `
output "instance_id" {
  description = "The ID of the EC2 instance"
  value       = aws_instance.example.id
}`,
			},
			expected: helper.Issues{
				{
					Rule:    NewTerraformVariableFileNameRule(),
					Message: "Output block should be declared in outputs.tf, not in main.tf",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 21},
					},
				},
			},
		},
		{
			name: "output block in outputs.tf - valid",
			files: map[string]string{
				"outputs.tf": `
output "instance_id" {
  description = "The ID of the EC2 instance"
  value       = aws_instance.example.id
}`,
			},
			expected: helper.Issues{},
		},
		{
			name: "locals block in main.tf - invalid",
			files: map[string]string{
				"main.tf": `
locals {
  common_tags = {
    Environment = "dev"
    Project     = "example"
  }
}`,
			},
			expected: helper.Issues{
				{
					Rule:    NewTerraformVariableFileNameRule(),
					Message: "Locals block should be declared in locals.tf, not in main.tf",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 7},
					},
				},
			},
		},
		{
			name: "locals block in locals.tf - valid",
			files: map[string]string{
				"locals.tf": `
locals {
  common_tags = {
    Environment = "dev"
    Project     = "example"
  }
}`,
			},
			expected: helper.Issues{},
		},
		{
			name: "variable block in outputs.tf - invalid",
			files: map[string]string{
				"outputs.tf": `
variable "should_not_be_here" {
  type = string
}`,
			},
			expected: helper.Issues{
				{
					Rule:    NewTerraformVariableFileNameRule(),
					Message: "Variable block should be declared in variables.tf, not in outputs.tf",
					Range: hcl.Range{
						Filename: "outputs.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 30},
					},
				},
			},
		},
		{
			name: "resource block in outputs.tf - invalid",
			files: map[string]string{
				"outputs.tf": `
resource "aws_instance" "should_not_be_here" {
  ami = "ami-12345678"
}`,
			},
			expected: helper.Issues{
				{
					Rule:    NewTerraformVariableFileNameRule(),
					Message: "Only output blocks should be declared in outputs.tf, found resource block",
					Range: hcl.Range{
						Filename: "outputs.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 45},
					},
				},
			},
		},
		{
			name: "output block in locals.tf - invalid",
			files: map[string]string{
				"locals.tf": `
output "should_not_be_here" {
  value = "test"
}`,
			},
			expected: helper.Issues{
				{
					Rule:    NewTerraformVariableFileNameRule(),
					Message: "Output block should be declared in outputs.tf, not in locals.tf",
					Range: hcl.Range{
						Filename: "locals.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 28},
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
