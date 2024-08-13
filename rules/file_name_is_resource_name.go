package rules

import (
	"regexp"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// FileNameIsResourceNameRule checks whether ...
type FileNameIsResourceNameRule struct {
	tflint.DefaultRule
}

// NewFileNameIsResourceNameRule returns a new rule
func NewFileNameIsResourceNameRule() *FileNameIsResourceNameRule {
	return &FileNameIsResourceNameRule{}
}

// Name returns the rule name
func (r *FileNameIsResourceNameRule) Name() string {
	return "file_name_as_resource_name"
}

// Enabled returns whether the rule is enabled by default
func (r *FileNameIsResourceNameRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *FileNameIsResourceNameRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *FileNameIsResourceNameRule) Link() string {
	return ""
}

// Check checks whether ...
func (r *FileNameIsResourceNameRule) Check(runner tflint.Runner) error {
	// This rule is an example to get a top-level resource attribute.
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}
	variablePattern := `^variable.tf$`
	variableRe, err := regexp.Compile(variablePattern)
	if err != nil {
		return err
	}
	for filename, file := range files {
		blocks, err := GetBlocksFromBody(file.Body)
		if err != nil {
			return err
		}
		// variable.tf
		if variableRe.MatchString(filename) {
			if len(*blocks) != len(blocks.Filter(Variable)) {
				return runner.EmitIssue(r, `Do not declare anything other than variable block in `+variablePattern, blocks.Exclude(Variable)[0].Range)
			}
		}
	}
	return nil
}
