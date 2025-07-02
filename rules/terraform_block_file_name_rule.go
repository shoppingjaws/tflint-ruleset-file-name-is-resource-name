package rules

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type TerraformBlockFileNameRule struct {
	tflint.DefaultRule
}

type TerraformBlockFileNameRuleConfig struct {
	VariableFile    string `hclext:"variable_file,optional"`
	OutputFile      string `hclext:"output_file,optional"`
	LocalsFile      string `hclext:"locals_file,optional"`
	TerraformFile   string `hclext:"terraform_file,optional"`
	ProviderFile    string `hclext:"provider_file,optional"`
	ModuleFile      string `hclext:"module_file,optional"`
	DataPrefix      string `hclext:"data_prefix,optional"`
	EphemeralPrefix string `hclext:"ephemeral_prefix,optional"`
	ImportFile      string `hclext:"import_file,optional"`
	MovedFile       string `hclext:"moved_file,optional"`
	RemovedFile     string `hclext:"removed_file,optional"`
	CheckFile       string `hclext:"check_file,optional"`
}

func NewTerraformBlockFileNameRule() *TerraformBlockFileNameRule {
	return &TerraformBlockFileNameRule{}
}

func (r *TerraformBlockFileNameRule) Name() string {
	return "terraform_block_file_name"
}

func (r *TerraformBlockFileNameRule) Enabled() bool {
	return true
}

func (r *TerraformBlockFileNameRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *TerraformBlockFileNameRule) Link() string {
	return ""
}

func (r *TerraformBlockFileNameRule) Check(runner tflint.Runner) error {
	path, err := runner.GetModulePath()
	if err != nil {
		return err
	}
	if !path.IsRoot() {
		return nil
	}

	// Decode custom configuration
	config := &TerraformBlockFileNameRuleConfig{}
	if err := runner.DecodeRuleConfig(r.Name(), config); err != nil {
		return err
	}

	// Set default values if not configured
	if config.VariableFile == "" {
		config.VariableFile = "variables.tf"
	}
	if config.OutputFile == "" {
		config.OutputFile = "outputs.tf"
	}
	if config.LocalsFile == "" {
		config.LocalsFile = "locals.tf"
	}
	if config.TerraformFile == "" {
		config.TerraformFile = "terraform.tf"
	}
	if config.ProviderFile == "" {
		config.ProviderFile = "provider.tf"
	}
	if config.ModuleFile == "" {
		config.ModuleFile = "module.tf"
	}
	if config.DataPrefix == "" {
		config.DataPrefix = "data_"
	}
	if config.EphemeralPrefix == "" {
		config.EphemeralPrefix = "ephemeral_"
	}
	if config.ImportFile == "" {
		config.ImportFile = "imports.tf"
	}
	if config.MovedFile == "" {
		config.MovedFile = "moved.tf"
	}
	if config.RemovedFile == "" {
		config.RemovedFile = "removed.tf"
	}
	if config.CheckFile == "" {
		config.CheckFile = "checks.tf"
	}

	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for name, file := range files {
		if err := r.checkFile(runner, name, file, config); err != nil {
			return err
		}
	}

	return nil
}

func (r *TerraformBlockFileNameRule) checkFile(runner tflint.Runner, filename string, file *hcl.File, config *TerraformBlockFileNameRuleConfig) error {
	basename := filepath.Base(filename)

	if !strings.HasSuffix(basename, ".tf") {
		return nil
	}

	body, diags := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
				Body:       &hclext.BodySchema{},
			},
			{
				Type:       "resource",
				LabelNames: []string{"type", "name"},
				Body:       &hclext.BodySchema{},
			},
			{
				Type:       "data",
				LabelNames: []string{"type", "name"},
				Body:       &hclext.BodySchema{},
			},
			{
				Type:       "ephemeral",
				LabelNames: []string{"type", "name"},
				Body:       &hclext.BodySchema{},
			},
			{
				Type:       "module",
				LabelNames: []string{"name"},
				Body:       &hclext.BodySchema{},
			},
			{
				Type:       "output",
				LabelNames: []string{"name"},
				Body:       &hclext.BodySchema{},
			},
			{
				Type:       "provider",
				LabelNames: []string{"name"},
				Body:       &hclext.BodySchema{},
			},
			{
				Type: "terraform",
				Body: &hclext.BodySchema{},
			},
			{
				Type: "locals",
				Body: &hclext.BodySchema{},
			},
			{
				Type: "import",
				Body: &hclext.BodySchema{},
			},
			{
				Type: "moved",
				Body: &hclext.BodySchema{},
			},
			{
				Type: "removed",
				Body: &hclext.BodySchema{},
			},
			{
				Type:       "check",
				LabelNames: []string{"name"},
				Body:       &hclext.BodySchema{},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if diags != nil {
		return diags
	}

	for _, block := range body.Blocks {
		blockRange := block.DefRange

		if blockRange.Filename != filename {
			continue
		}

		// First check special files - these take priority
		if basename == config.VariableFile {
			if block.Type != "variable" {
				if block.Type == "output" {
					blockManager := NewBlockManager(runner)
					fixFunc := blockManager.CreateFixFunction(block, config.OutputFile)
					if err := runner.EmitIssueWithFix(
						r,
						fmt.Sprintf("Output block should be declared in %s, not in %s", config.OutputFile, basename),
						blockRange,
						fixFunc,
					); err != nil {
						return err
					}
				} else if block.Type == "locals" {
					blockManager := NewBlockManager(runner)
					fixFunc := blockManager.CreateFixFunction(block, config.LocalsFile)
					if err := runner.EmitIssueWithFix(
						r,
						fmt.Sprintf("Locals block should be declared in %s, not in %s", config.LocalsFile, basename),
						blockRange,
						fixFunc,
					); err != nil {
						return err
					}
				} else if block.Type == "resource" && len(block.Labels) >= 1 {
					resourceType := block.Labels[0]
					expectedFilename := resourceType + ".tf"
					blockManager := NewBlockManager(runner)
					fixFunc := blockManager.CreateFixFunction(block, expectedFilename)
					if err := runner.EmitIssueWithFix(
						r,
						fmt.Sprintf("Resource block '%s' should be declared in %s, not in %s", resourceType, expectedFilename, basename),
						blockRange,
						fixFunc,
					); err != nil {
						return err
					}
				} else if block.Type == "data" && len(block.Labels) >= 1 {
					dataType := block.Labels[0]
					expectedFilename := config.DataPrefix + dataType + ".tf"
					blockManager := NewBlockManager(runner)
					fixFunc := blockManager.CreateFixFunction(block, expectedFilename)
					if err := runner.EmitIssueWithFix(
						r,
						fmt.Sprintf("Data block '%s' should be declared in %s, not in %s", dataType, expectedFilename, basename),
						blockRange,
						fixFunc,
					); err != nil {
						return err
					}
				} else if block.Type == "ephemeral" && len(block.Labels) >= 1 {
					ephemeralType := block.Labels[0]
					expectedFilename := config.EphemeralPrefix + ephemeralType + ".tf"
					blockManager := NewBlockManager(runner)
					fixFunc := blockManager.CreateFixFunction(block, expectedFilename)
					if err := runner.EmitIssueWithFix(
						r,
						fmt.Sprintf("Ephemeral block '%s' should be declared in %s, not in %s", ephemeralType, expectedFilename, basename),
						blockRange,
						fixFunc,
					); err != nil {
						return err
					}
				} else if block.Type == "import" {
					blockManager := NewBlockManager(runner)
					fixFunc := blockManager.CreateFixFunction(block, config.ImportFile)
					if err := runner.EmitIssueWithFix(
						r,
						fmt.Sprintf("Import block should be declared in %s, not in %s", config.ImportFile, basename),
						blockRange,
						fixFunc,
					); err != nil {
						return err
					}
				} else if block.Type == "moved" {
					blockManager := NewBlockManager(runner)
					fixFunc := blockManager.CreateFixFunction(block, config.MovedFile)
					if err := runner.EmitIssueWithFix(
						r,
						fmt.Sprintf("Moved block should be declared in %s, not in %s", config.MovedFile, basename),
						blockRange,
						fixFunc,
					); err != nil {
						return err
					}
				} else if block.Type == "removed" {
					blockManager := NewBlockManager(runner)
					fixFunc := blockManager.CreateFixFunction(block, config.RemovedFile)
					if err := runner.EmitIssueWithFix(
						r,
						fmt.Sprintf("Removed block should be declared in %s, not in %s", config.RemovedFile, basename),
						blockRange,
						fixFunc,
					); err != nil {
						return err
					}
				} else if block.Type == "check" {
					blockManager := NewBlockManager(runner)
					fixFunc := blockManager.CreateFixFunction(block, config.CheckFile)
					if err := runner.EmitIssueWithFix(
						r,
						fmt.Sprintf("Check block should be declared in %s, not in %s", config.CheckFile, basename),
						blockRange,
						fixFunc,
					); err != nil {
						return err
					}
				} else {
					if err := runner.EmitIssue(
						r,
						fmt.Sprintf("Only variable blocks should be declared in %s, found %s block", config.VariableFile, block.Type),
						blockRange,
					); err != nil {
						return err
					}
				}
			}
		} else if basename == config.OutputFile {
			if block.Type != "output" {
				if block.Type == "variable" {
					blockManager := NewBlockManager(runner)
					fixFunc := blockManager.CreateFixFunction(block, config.VariableFile)
					if err := runner.EmitIssueWithFix(
						r,
						fmt.Sprintf("Variable block should be declared in %s, not in %s", config.VariableFile, basename),
						blockRange,
						fixFunc,
					); err != nil {
						return err
					}
				} else if block.Type == "locals" {
					blockManager := NewBlockManager(runner)
					fixFunc := blockManager.CreateFixFunction(block, config.LocalsFile)
					if err := runner.EmitIssueWithFix(
						r,
						fmt.Sprintf("Locals block should be declared in %s, not in %s", config.LocalsFile, basename),
						blockRange,
						fixFunc,
					); err != nil {
						return err
					}
				} else if block.Type == "resource" && len(block.Labels) >= 1 {
					resourceType := block.Labels[0]
					expectedFilename := resourceType + ".tf"
					blockManager := NewBlockManager(runner)
					fixFunc := blockManager.CreateFixFunction(block, expectedFilename)
					if err := runner.EmitIssueWithFix(
						r,
						fmt.Sprintf("Resource block '%s' should be declared in %s, not in %s", resourceType, expectedFilename, basename),
						blockRange,
						fixFunc,
					); err != nil {
						return err
					}
				} else if block.Type == "data" && len(block.Labels) >= 1 {
					dataType := block.Labels[0]
					expectedFilename := config.DataPrefix + dataType + ".tf"
					blockManager := NewBlockManager(runner)
					fixFunc := blockManager.CreateFixFunction(block, expectedFilename)
					if err := runner.EmitIssueWithFix(
						r,
						fmt.Sprintf("Data block '%s' should be declared in %s, not in %s", dataType, expectedFilename, basename),
						blockRange,
						fixFunc,
					); err != nil {
						return err
					}
				} else if block.Type == "ephemeral" && len(block.Labels) >= 1 {
					ephemeralType := block.Labels[0]
					expectedFilename := config.EphemeralPrefix + ephemeralType + ".tf"
					blockManager := NewBlockManager(runner)
					fixFunc := blockManager.CreateFixFunction(block, expectedFilename)
					if err := runner.EmitIssueWithFix(
						r,
						fmt.Sprintf("Ephemeral block '%s' should be declared in %s, not in %s", ephemeralType, expectedFilename, basename),
						blockRange,
						fixFunc,
					); err != nil {
						return err
					}
				} else {
					if err := runner.EmitIssue(
						r,
						fmt.Sprintf("Only output blocks should be declared in %s, found %s block", config.OutputFile, block.Type),
						blockRange,
					); err != nil {
						return err
					}
				}
			}
		} else if basename == config.LocalsFile {
			if block.Type != "locals" {
				if block.Type == "variable" {
					blockManager := NewBlockManager(runner)
					fixFunc := blockManager.CreateFixFunction(block, config.VariableFile)
					if err := runner.EmitIssueWithFix(
						r,
						fmt.Sprintf("Variable block should be declared in %s, not in %s", config.VariableFile, basename),
						blockRange,
						fixFunc,
					); err != nil {
						return err
					}
				} else if block.Type == "output" {
					blockManager := NewBlockManager(runner)
					fixFunc := blockManager.CreateFixFunction(block, config.OutputFile)
					if err := runner.EmitIssueWithFix(
						r,
						fmt.Sprintf("Output block should be declared in %s, not in %s", config.OutputFile, basename),
						blockRange,
						fixFunc,
					); err != nil {
						return err
					}
				} else if block.Type == "resource" && len(block.Labels) >= 1 {
					resourceType := block.Labels[0]
					expectedFilename := resourceType + ".tf"
					blockManager := NewBlockManager(runner)
					fixFunc := blockManager.CreateFixFunction(block, expectedFilename)
					if err := runner.EmitIssueWithFix(
						r,
						fmt.Sprintf("Resource block '%s' should be declared in %s, not in %s", resourceType, expectedFilename, basename),
						blockRange,
						fixFunc,
					); err != nil {
						return err
					}
				} else if block.Type == "data" && len(block.Labels) >= 1 {
					dataType := block.Labels[0]
					expectedFilename := config.DataPrefix + dataType + ".tf"
					blockManager := NewBlockManager(runner)
					fixFunc := blockManager.CreateFixFunction(block, expectedFilename)
					if err := runner.EmitIssueWithFix(
						r,
						fmt.Sprintf("Data block '%s' should be declared in %s, not in %s", dataType, expectedFilename, basename),
						blockRange,
						fixFunc,
					); err != nil {
						return err
					}
				} else if block.Type == "ephemeral" && len(block.Labels) >= 1 {
					ephemeralType := block.Labels[0]
					expectedFilename := config.EphemeralPrefix + ephemeralType + ".tf"
					blockManager := NewBlockManager(runner)
					fixFunc := blockManager.CreateFixFunction(block, expectedFilename)
					if err := runner.EmitIssueWithFix(
						r,
						fmt.Sprintf("Ephemeral block '%s' should be declared in %s, not in %s", ephemeralType, expectedFilename, basename),
						blockRange,
						fixFunc,
					); err != nil {
						return err
					}
				} else {
					if err := runner.EmitIssue(
						r,
						fmt.Sprintf("Only locals blocks should be declared in %s, found %s block", config.LocalsFile, block.Type),
						blockRange,
					); err != nil {
						return err
					}
				}
			}
		} else if basename == config.TerraformFile {
			if block.Type != "terraform" {
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf("Only terraform blocks should be declared in %s, found %s block", config.TerraformFile, block.Type),
					blockRange,
				); err != nil {
					return err
				}
			}
		} else if basename == config.ProviderFile {
			if block.Type != "provider" {
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf("Only provider blocks should be declared in %s, found %s block", config.ProviderFile, block.Type),
					blockRange,
				); err != nil {
					return err
				}
			}
		} else if basename == config.ModuleFile {
			if block.Type != "module" {
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf("Only module blocks should be declared in %s, found %s block", config.ModuleFile, block.Type),
					blockRange,
				); err != nil {
					return err
				}
			}
		} else if basename == config.ImportFile {
			if block.Type != "import" {
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf("Only import blocks should be declared in %s, found %s block", config.ImportFile, block.Type),
					blockRange,
				); err != nil {
					return err
				}
			}
		} else if basename == config.MovedFile {
			if block.Type != "moved" {
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf("Only moved blocks should be declared in %s, found %s block", config.MovedFile, block.Type),
					blockRange,
				); err != nil {
					return err
				}
			}
		} else if basename == config.RemovedFile {
			if block.Type != "removed" {
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf("Only removed blocks should be declared in %s, found %s block", config.RemovedFile, block.Type),
					blockRange,
				); err != nil {
					return err
				}
			}
		} else if basename == config.CheckFile {
			if block.Type != "check" {
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf("Only check blocks should be declared in %s, found %s block", config.CheckFile, block.Type),
					blockRange,
				); err != nil {
					return err
				}
			}
		} else {
			// Handle non-special files
			switch block.Type {
			case "variable":
				blockManager := NewBlockManager(runner)
				fixFunc := blockManager.CreateFixFunction(block, config.VariableFile)

				if err := runner.EmitIssueWithFix(
					r,
					fmt.Sprintf("Variable block should be declared in %s, not in %s", config.VariableFile, basename),
					blockRange,
					fixFunc,
				); err != nil {
					return err
				}
			case "output":
				blockManager := NewBlockManager(runner)
				fixFunc := blockManager.CreateFixFunction(block, config.OutputFile)

				if err := runner.EmitIssueWithFix(
					r,
					fmt.Sprintf("Output block should be declared in %s, not in %s", config.OutputFile, basename),
					blockRange,
					fixFunc,
				); err != nil {
					return err
				}
			case "locals":
				blockManager := NewBlockManager(runner)
				fixFunc := blockManager.CreateFixFunction(block, config.LocalsFile)

				if err := runner.EmitIssueWithFix(
					r,
					fmt.Sprintf("Locals block should be declared in %s, not in %s", config.LocalsFile, basename),
					blockRange,
					fixFunc,
				); err != nil {
					return err
				}
			case "resource":
				if len(block.Labels) >= 1 {
					resourceType := block.Labels[0]
					expectedFilename := resourceType + ".tf"

					if basename != expectedFilename {
						blockManager := NewBlockManager(runner)
						fixFunc := blockManager.CreateFixFunction(block, expectedFilename)

						if err := runner.EmitIssueWithFix(
							r,
							fmt.Sprintf("Resource block '%s' should be declared in %s, not in %s", resourceType, expectedFilename, basename),
							blockRange,
							fixFunc,
						); err != nil {
							return err
						}
					}
				}
			case "data":
				if len(block.Labels) >= 1 {
					dataType := block.Labels[0]
					expectedFilename := config.DataPrefix + dataType + ".tf"

					if basename != expectedFilename {
						blockManager := NewBlockManager(runner)
						fixFunc := blockManager.CreateFixFunction(block, expectedFilename)

						if err := runner.EmitIssueWithFix(
							r,
							fmt.Sprintf("Data block '%s' should be declared in %s, not in %s", dataType, expectedFilename, basename),
							blockRange,
							fixFunc,
						); err != nil {
							return err
						}
					}
				}
			case "ephemeral":
				if len(block.Labels) >= 1 {
					ephemeralType := block.Labels[0]
					expectedFilename := config.EphemeralPrefix + ephemeralType + ".tf"

					if basename != expectedFilename {
						blockManager := NewBlockManager(runner)
						fixFunc := blockManager.CreateFixFunction(block, expectedFilename)

						if err := runner.EmitIssueWithFix(
							r,
							fmt.Sprintf("Ephemeral block '%s' should be declared in %s, not in %s", ephemeralType, expectedFilename, basename),
							blockRange,
							fixFunc,
						); err != nil {
							return err
						}
					}
				}
			case "terraform":
				blockManager := NewBlockManager(runner)
				fixFunc := blockManager.CreateFixFunction(block, config.TerraformFile)

				if err := runner.EmitIssueWithFix(
					r,
					fmt.Sprintf("Terraform block should be declared in %s, not in %s", config.TerraformFile, basename),
					blockRange,
					fixFunc,
				); err != nil {
					return err
				}
			case "provider":
				blockManager := NewBlockManager(runner)
				fixFunc := blockManager.CreateFixFunction(block, config.ProviderFile)

				if err := runner.EmitIssueWithFix(
					r,
					fmt.Sprintf("Provider block should be declared in %s, not in %s", config.ProviderFile, basename),
					blockRange,
					fixFunc,
				); err != nil {
					return err
				}
			case "module":
				blockManager := NewBlockManager(runner)
				fixFunc := blockManager.CreateFixFunction(block, config.ModuleFile)

				if err := runner.EmitIssueWithFix(
					r,
					fmt.Sprintf("Module block should be declared in %s, not in %s", config.ModuleFile, basename),
					blockRange,
					fixFunc,
				); err != nil {
					return err
				}
			case "import":
				blockManager := NewBlockManager(runner)
				fixFunc := blockManager.CreateFixFunction(block, config.ImportFile)

				if err := runner.EmitIssueWithFix(
					r,
					fmt.Sprintf("Import block should be declared in %s, not in %s", config.ImportFile, basename),
					blockRange,
					fixFunc,
				); err != nil {
					return err
				}
			case "moved":
				blockManager := NewBlockManager(runner)
				fixFunc := blockManager.CreateFixFunction(block, config.MovedFile)

				if err := runner.EmitIssueWithFix(
					r,
					fmt.Sprintf("Moved block should be declared in %s, not in %s", config.MovedFile, basename),
					blockRange,
					fixFunc,
				); err != nil {
					return err
				}
			case "removed":
				blockManager := NewBlockManager(runner)
				fixFunc := blockManager.CreateFixFunction(block, config.RemovedFile)

				if err := runner.EmitIssueWithFix(
					r,
					fmt.Sprintf("Removed block should be declared in %s, not in %s", config.RemovedFile, basename),
					blockRange,
					fixFunc,
				); err != nil {
					return err
				}
			case "check":
				blockManager := NewBlockManager(runner)
				fixFunc := blockManager.CreateFixFunction(block, config.CheckFile)

				if err := runner.EmitIssueWithFix(
					r,
					fmt.Sprintf("Check block should be declared in %s, not in %s", config.CheckFile, basename),
					blockRange,
					fixFunc,
				); err != nil {
					return err
				}
			}

			// Check if this is a resource-specific file that contains non-matching blocks
			if r.isResourceFile(basename, config.DataPrefix, config.EphemeralPrefix) {
				expectedResourceType := strings.TrimSuffix(basename, ".tf")
				if block.Type == "resource" {
					if len(block.Labels) >= 1 && block.Labels[0] != expectedResourceType {
						if err := runner.EmitIssue(
							r,
							fmt.Sprintf("Only '%s' resource blocks should be declared in %s, found '%s' resource block", expectedResourceType, basename, block.Labels[0]),
							blockRange,
						); err != nil {
							return err
						}
					}
				} else {
					if err := runner.EmitIssue(
						r,
						fmt.Sprintf("Only '%s' resource blocks should be declared in %s, found %s block", expectedResourceType, basename, block.Type),
						blockRange,
					); err != nil {
						return err
					}
				}
			}

			// Check if this is a data-specific file that contains non-matching blocks
			if r.isDataFile(basename, config.DataPrefix) {
				expectedDataType := strings.TrimPrefix(strings.TrimSuffix(basename, ".tf"), config.DataPrefix)
				if block.Type == "data" {
					if len(block.Labels) >= 1 && block.Labels[0] != expectedDataType {
						if err := runner.EmitIssue(
							r,
							fmt.Sprintf("Only '%s' data blocks should be declared in %s, found '%s' data block", expectedDataType, basename, block.Labels[0]),
							blockRange,
						); err != nil {
							return err
						}
					}
				} else {
					if err := runner.EmitIssue(
						r,
						fmt.Sprintf("Only '%s' data blocks should be declared in %s, found %s block", expectedDataType, basename, block.Type),
						blockRange,
					); err != nil {
						return err
					}
				}
			}

			// Check if this is an ephemeral-specific file that contains non-matching blocks
			if r.isEphemeralFile(basename, config.EphemeralPrefix) {
				expectedEphemeralType := strings.TrimPrefix(strings.TrimSuffix(basename, ".tf"), config.EphemeralPrefix)
				if block.Type == "ephemeral" {
					if len(block.Labels) >= 1 && block.Labels[0] != expectedEphemeralType {
						if err := runner.EmitIssue(
							r,
							fmt.Sprintf("Only '%s' ephemeral blocks should be declared in %s, found '%s' ephemeral block", expectedEphemeralType, basename, block.Labels[0]),
							blockRange,
						); err != nil {
							return err
						}
					}
				} else {
					if err := runner.EmitIssue(
						r,
						fmt.Sprintf("Only '%s' ephemeral blocks should be declared in %s, found %s block", expectedEphemeralType, basename, block.Type),
						blockRange,
					); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// isDataFile determines if a filename follows the data naming pattern with given prefix
func (r *TerraformBlockFileNameRule) isDataFile(basename string, dataPrefix string) bool {
	if !strings.HasSuffix(basename, ".tf") {
		return false
	}

	filename := strings.TrimSuffix(basename, ".tf")
	return strings.HasPrefix(filename, dataPrefix) && strings.Contains(filename, "_")
}

// isEphemeralFile determines if a filename follows the ephemeral naming pattern with given prefix
func (r *TerraformBlockFileNameRule) isEphemeralFile(basename string, ephemeralPrefix string) bool {
	if !strings.HasSuffix(basename, ".tf") {
		return false
	}

	filename := strings.TrimSuffix(basename, ".tf")
	return strings.HasPrefix(filename, ephemeralPrefix) && strings.Contains(filename, "_")
}

// isResourceFile determines if a filename follows the resource naming pattern
// Resource files should be named like "aws_instance.tf", "azurerm_virtual_machine.tf", etc.
// This is different from special files like "variables.tf", "outputs.tf", "locals.tf", "main.tf"
func (r *TerraformBlockFileNameRule) isResourceFile(basename string, dataPrefix string, ephemeralPrefix string) bool {
	if !strings.HasSuffix(basename, ".tf") {
		return false
	}

	filename := strings.TrimSuffix(basename, ".tf")

	// Exclude special files
	specialFiles := []string{"variables", "outputs", "locals", "main", "providers", "versions", "terraform", "imports", "moved", "removed", "checks"}
	for _, special := range specialFiles {
		if filename == special {
			return false
		}
	}

	// Exclude data files (they have different naming pattern)
	if strings.HasPrefix(filename, dataPrefix) {
		return false
	}

	// Exclude ephemeral files (they have different naming pattern)
	if strings.HasPrefix(filename, ephemeralPrefix) {
		return false
	}

	// Check if it contains underscores (indicating resource type pattern like aws_instance, azurerm_vm)
	return strings.Contains(filename, "_")
}
