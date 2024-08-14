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
	variable_file_name_pattern string `hclext:"variable_file_name_pattern,optional"`
	locals_file_name_pattern   string `hclext:"locals_file_name_pattern,optional"`
	provider_file_name_pattern string `hclext:"provider_file_name_pattern,optional"`
	output_file_name_pattern   string `hclext:"output_file_name_pattern,optional"`
	module_file_name_pattern   string `hclext:"module_file_name_pattern,optional"`
	data_file_name_pattern     string `hclext:"data_file_name_pattern,optional"`
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
	logger.Debug("Config %s", config.data_file_name_pattern)
	if err := runner.DecodeRuleConfig(r.Name(), &config); err != nil {
		logger.Error("Error decoding rule config: %s", err)
		return err
	}
	var variablePattern string
	if config.variable_file_name_pattern == "" {
		variablePattern = `^variable.tf$`
	} else {
		variablePattern = config.variable_file_name_pattern
	}
	variableRe, err := regexp.Compile(variablePattern)
	if err != nil {
		return err
	}

	var localsPattern string
	if config.locals_file_name_pattern == "" {
		localsPattern = `^locals.tf$`
	} else {
		localsPattern = config.locals_file_name_pattern
	}
	localsRe, err := regexp.Compile(localsPattern)
	if err != nil {
		return err
	}

	var providerPattern string
	if config.provider_file_name_pattern == "" {
		providerPattern = `^provider.tf$`
	} else {
		providerPattern = config.provider_file_name_pattern
	}
	providerRe, err := regexp.Compile(providerPattern)
	if err != nil {
		return err
	}

	var outputPattern string
	if config.output_file_name_pattern == "" {
		outputPattern = `^output.tf$`
	} else {
		outputPattern = config.output_file_name_pattern
	}
	outputRe, err := regexp.Compile(outputPattern)
	if err != nil {
		return err
	}

	var modulePattern string
	if config.module_file_name_pattern == "" {
		modulePattern = `^module.tf$`
	} else {
		modulePattern = config.module_file_name_pattern
	}
	moduleRe, err := regexp.Compile(modulePattern)
	if err != nil {
		return err
	}

	var dataPattern string
	if config.data_file_name_pattern == "" {
		dataPattern = `^data_.*.tf$`
	} else {
		dataPattern = config.data_file_name_pattern
	}
	dataRe, err := regexp.Compile(dataPattern)
	if err != nil {
		return err
	}
	for filename, file := range files {
		logger.Debug("File: %s", filename)
		blocks, err := GetBlocksFromBody(file.Body)
		if err != nil {
			return err
		}
		// variable.tf
		if variableRe.MatchString(filename) {
			logger.Debug("variable")
			// check if there is any block other than variable
			if len(*blocks) != len(blocks.Filter(Variable)) {
				return runner.EmitIssue(r, `Do not declare anything other than Variable block in `+filename, blocks.Exclude(Variable)[0].Range)
			}
		} else
		// locals.tf
		if localsRe.MatchString(filename) {
			logger.Debug("locals")
			if len(*blocks) != len(blocks.Filter(Locals)) {
				return runner.EmitIssue(r, `Do not declare anything other than Locals block in `+filename, blocks.Exclude(Locals)[0].Range)
			}
		} else
		// provider.tf
		if providerRe.MatchString(filename) {
			logger.Debug("provider")
			if len(*blocks) != len(blocks.Filter(Provider)) {
				return runner.EmitIssue(r, `Do not declare anything other than Provider block in `+filename, blocks.Exclude(Provider)[0].Range)
			}
		} else
		// output.tf
		if outputRe.MatchString(filename) {
			logger.Debug("output")
			if len(*blocks) != len(blocks.Filter(Output)) {
				return runner.EmitIssue(r, `Do not declare anything other than Output block in `+filename, blocks.Exclude(Output)[0].Range)
			}
		} else
		// module.tf
		if moduleRe.MatchString(filename) {
			logger.Debug("module")
			if len(*blocks) != len(blocks.Filter(Module)) {
				return runner.EmitIssue(r, `Do not declare anything other than Module block in `+filename, blocks.Exclude(Module)[0].Range)
			}
		} else
		// data
		if dataRe.MatchString(filename) {
			logger.Debug("data")
			if len(*blocks) != len(blocks.Filter(Data)) {
				return runner.EmitIssue(r, `Do not declare anything other than Data block in `+filename, blocks.Exclude(Data)[0].Range)
			}
			for _, data := range blocks.Filter(Data) {
				if "data_"+*data.Type+".tf" != filename {
					return runner.EmitIssue(r, `File name should be the same as the data type `+"data_"+*data.Type+".tf", data.Range)
				}
			}
		} else
		// resource
		{
			logger.Debug("else")
			// check if there is any block other than resource
			if len(*blocks) != len(blocks.Filter(Resource)) {
				return runner.EmitIssue(r, `Do not declare anything other than Resource block in `+filename, blocks.Exclude(Data)[0].Range)
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
