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
		Expected helper.Issues
	}{
		{
			FileName: "resource.tf",
			Name:     "Fail if file name is not resource type",
			Content: `
resource "aws_instance" "web" {
		instance_type = "t2.micro"
}`,
			Expected: helper.Issues{},
		},
		{
			FileName: "variable.tf",
			Name:     "decline the declaration of non variable block with variable.tf",
			Content: `
resource "aws_instance" "web" {
    instance_type = "t2.micro"
}
resource "aws_instance" "db" {
    instance_type = "t2.micro"
}
`,
			Expected: helper.Issues{{
				Rule:    NewFileNameIsResourceNameRule(),
				Message: "Do not declare anything other than variable block in ^variable.tf$",
				Range:   hcl.Range{Filename: "variable.tf", Start: hcl.Pos{Line: 2, Column: 1}, End: hcl.Pos{Line: 2, Column: 30}},
			}},
		},
	}

	rule := NewFileNameIsResourceNameRule()

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{test.FileName: test.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, test.Expected, runner.Issues)
		})
	}
}
