package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_FileNameIsResourceName(t *testing.T) {
	tests := []struct {
		FileName string
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
	}{
		{
			FileName: "variables.tf",
			Name:     "accept variable if file name is variables.tf with config",
			Content: `
variable "variable_name" {}
`,
			Config: `
rule "file_name_is_resource_name" {
  enabled = true
  variable_file_name_pattern = "^variables.tf$"
}`,
			Expected: helper.Issues{},
		},
		{
			FileName: "variable.tf",
			Name:     "accept variable if file name is variable.tf with config",
			Content: `
variable "variable_name" {}
`,
			Expected: helper.Issues{},
		},
		{
			FileName: "locals.tf",
			Name:     "accept locals block if file name is locals.tf without config",
			Content: `
locals {}
`,
			Expected: helper.Issues{},
		},
		{
			FileName: "aws_instance.tf",
			Name:     "accept data if file name is data_$type",
			Content: `
resource "aws_instance" "web" {
		instance_type = "t2.micro"
}`,
			Expected: helper.Issues{},
		},
		{
			FileName: "aws_instance.tf",
			Name:     "decline resource if file name is $invalid_type",
			Content: `
resource "aws_security_group" "web" {
}`,
			Expected: helper.Issues{
				{
					Rule:    NewFileNameIsResourceNameRule(),
					Message: "File name should be the same as the resource type aws_security_group.tf",
					Range: hcl.Range{
						Filename: "aws_instance.tf", Start: hcl.Pos{Line: 2, Column: 1}, End: hcl.Pos{Line: 2, Column: 36}},
				},
			},
		},
		{
			FileName: "data_aws_instance.tf",
			Name:     "decline data if file name is data_$invalid_type",
			Content: `
data "aws_security_group" "web" {
}`,
			Expected: helper.Issues{
				{
					Rule:    NewFileNameIsResourceNameRule(),
					Message: "File name should be the same as the data type aws_security_group.tf",
					Range: hcl.Range{
						Filename: "data_aws_instance.tf", Start: hcl.Pos{Line: 2, Column: 1}, End: hcl.Pos{Line: 2, Column: 32}},
				},
			},
		},
		{
			FileName: "data_aws_instance.tf",
			Name:     "accept data if file name is data_$type",
			Content: `
data "aws_instance" "web" {
}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewFileNameIsResourceNameRule()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{test.FileName: test.Content, ".tflint.hcl": test.Config})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, test.Expected, runner.Issues)
		})
	}
}
