package rules

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type TerraformVariableFileNameRule struct {
	tflint.DefaultRule
}

func NewTerraformVariableFileNameRule() *TerraformVariableFileNameRule {
	return &TerraformVariableFileNameRule{}
}

func (r *TerraformVariableFileNameRule) Name() string {
	return "terraform_variable_file_name"
}

func (r *TerraformVariableFileNameRule) Enabled() bool {
	return true
}

func (r *TerraformVariableFileNameRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *TerraformVariableFileNameRule) Link() string {
	return ""
}

func (r *TerraformVariableFileNameRule) Check(runner tflint.Runner) error {
	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		return nil
	}

	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for name, file := range files {
		if err := r.checkFile(runner, name, file); err != nil {
			return err
		}
	}

	return nil
}

func (r *TerraformVariableFileNameRule) checkFile(runner tflint.Runner, filename string, file *hcl.File) error {
	basename := filepath.Base(filename)
	
	if !strings.HasSuffix(basename, ".tf") {
		return nil
	}

	body, diags := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: "variable",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if diags != nil {
		return diags
	}

	for _, block := range body.Blocks {
		if block.Type == "variable" {
			blockRange := block.DefRange
			
			if blockRange.Filename != filename {
				continue
			}

			if basename != "variables.tf" {
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf("Variable block should be declared in variables.tf, not in %s", basename),
					blockRange,
				); err != nil {
					return err
				}
			}
		}
	}

	return nil
}