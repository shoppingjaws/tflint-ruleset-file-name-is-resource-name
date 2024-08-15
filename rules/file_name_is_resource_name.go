package rules

import (
	"regexp"

	"github.com/terraform-linters/tflint-plugin-sdk/logger"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// FileNameIsResourceNameRule checks whether ...
type FileNameIsResourceNameRule struct {
	tflint.DefaultRule
}

type FileNameIsResourceNameRuleConfig struct {
	VariableFileNamePattern string `hclext:"variable_file_name_pattern,optional"`
	LocalsFileNamePattern   string `hclext:"locals_file_name_pattern,optional"`
	// provider_file_name_pattern string `hclext:"provider_file_name_pattern,optional"`
	// output_file_name_pattern   string `hclext:"output_file_name_pattern,optional"`
	// module_file_name_pattern   string `hclext:"module_file_name_pattern,optional"`
	// data_file_name_pattern     string `hclext:"data_file_name_pattern,optional"`
}

// NewFileNameIsResourceNameRule returns a new rule
func NewFileNameIsResourceNameRule() *FileNameIsResourceNameRule {
	return &FileNameIsResourceNameRule{}
}

// Name returns the rule name
func (r *FileNameIsResourceNameRule) Name() string {
	return "file_name_is_resource_name"
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
	logger.Debug("Check init")
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}
	config := &FileNameIsResourceNameRuleConfig{}
	if config.VariableFileNamePattern == "" {
		config.VariableFileNamePattern = `^variable.tf$`
		logger.Debug("Config: variable = %s", config.VariableFileNamePattern)
	}
	if config.LocalsFileNamePattern == "" {
		config.LocalsFileNamePattern = `^locals.tf$`
		logger.Debug("Config: local =  %s", config.LocalsFileNamePattern)
	}
	if err := runner.DecodeRuleConfig(r.Name(), config); err != nil {
		logger.Error("Error decoding rule config: %s", err)
		return err
	}
	variableRe := regexp.MustCompile(config.VariableFileNamePattern)
	localsRe := regexp.MustCompile(config.LocalsFileNamePattern)
	for filename, file := range files {
		logger.Debug("File: %s", filename)
		blocks, err := GetBlocksFromBody(file.Body)
		if err != nil {
			return err
		}
		if variableRe.MatchString(filename) { // variable
			logger.Debug("variable found %s", filename)
			if len(*blocks) != len(blocks.Filter(Variable)) {
				return runner.EmitIssue(r, `Do not declare anything other than Variable block in `+filename, blocks.Exclude(Variable)[0].Range)
			}
		} else if localsRe.MatchString(filename) { //locals
			logger.Debug("locals found %s", filename)
			if len(*blocks) != len(blocks.Filter(Locals)) {
				return runner.EmitIssue(r, `Do not declare anything other than Locals block in `+filename, blocks.Exclude(Locals)[0].Range)
			}
		} else { // reosurce
			logger.Debug("resource")
			// check if there is any block other than resource
			if len(*blocks) != len(blocks.Filter(Resource)) {
				return runner.EmitIssue(r, `Do not declare anything other than Resource block of `+toBlockName(filename)+` in `+filename, blocks.Exclude(Resource)[0].Range)
			}
			for _, resource := range blocks.Filter(Resource) {
				if *resource.Type+".tf" != filename {
					return runner.EmitIssue(r, `File name should be the same as the resource type `+*resource.Type+".tf", resource.Range)
				}
			}
		}
	}
	return nil
}
