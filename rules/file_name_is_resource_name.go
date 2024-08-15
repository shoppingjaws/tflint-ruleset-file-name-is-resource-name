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
	ProviderFileNamePattern string `hclext:"provider_file_name_pattern,optional"`
	OutputFileNamePattern   string `hclext:"output_file_name_pattern,optional"`
	ModuleFileNamePattern   string `hclext:"module_file_name_pattern,optional"`
	DataFileNamePattern     string `hclext:"data_file_name_pattern,optional"`
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
	if config.ProviderFileNamePattern == "" {
		config.ProviderFileNamePattern = `^provider.tf$`
		logger.Debug("Config: provider =  %s", config.ProviderFileNamePattern)
	}
	if config.OutputFileNamePattern == "" {
		config.OutputFileNamePattern = `^output.tf$`
		logger.Debug("Config: output =  %s", config.OutputFileNamePattern)
	}
	if config.ModuleFileNamePattern == "" {
		config.ModuleFileNamePattern = `^module.tf$`
		logger.Debug("Config: module =  %s", config.ModuleFileNamePattern)
	}
	if config.DataFileNamePattern == "" {
		config.ModuleFileNamePattern = `^data_.*.tf$`
		logger.Debug("Config: data =  %s", config.DataFileNamePattern)
	}
	if err := runner.DecodeRuleConfig(r.Name(), config); err != nil {
		logger.Error("Error decoding rule config: %s", err)
		return err
	}
	variableRe := regexp.MustCompile(config.VariableFileNamePattern)
	localsRe := regexp.MustCompile(config.LocalsFileNamePattern)
	providerRe := regexp.MustCompile(config.ProviderFileNamePattern)
	outputRe := regexp.MustCompile(config.OutputFileNamePattern)
	moduleRe := regexp.MustCompile(config.ModuleFileNamePattern)
	dataRe := regexp.MustCompile(`^data.tf$`)
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
		} else if providerRe.MatchString(filename) { //provider
			logger.Debug("provider found %s", filename)
			if len(*blocks) != len(blocks.Filter(Provider)) {
				return runner.EmitIssue(r, `Do not declare anything other than Provider block in `+filename, blocks.Exclude(Provider)[0].Range)
			}
		} else if outputRe.MatchString(filename) { //output
			logger.Debug("output found %s", filename)
			if len(*blocks) != len(blocks.Filter(Output)) {
				return runner.EmitIssue(r, `Do not declare anything other than output block in `+filename, blocks.Exclude(Output)[0].Range)
			}
		} else if moduleRe.MatchString(filename) { //module
			logger.Debug("module found %s", filename)
			if len(*blocks) != len(blocks.Filter(Module)) {
				return runner.EmitIssue(r, `Do not declare anything other than Module block in `+filename, blocks.Exclude(Module)[0].Range)
			}
		} else if dataRe.MatchString(filename) { //data
			logger.Debug("data found %s", filename)
			if len(*blocks) != len(blocks.Filter(Data)) {
				return runner.EmitIssue(r, `Do not declare anything other than Data block of `+toBlockName(filename)+` in `+filename, blocks.Exclude(Data)[0].Range)
			}
			for _, data := range blocks.Filter(Data) {
				if *data.Type+".tf" != filename {
					return runner.EmitIssue(r, `File name should be the same as the data type `+*data.Type+".tf", data.Range)
				}
			}
		} else { // reosurce
			logger.Debug("resource found %s", filename)
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
