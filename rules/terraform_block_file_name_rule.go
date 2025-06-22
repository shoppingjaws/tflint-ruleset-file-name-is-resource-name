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

func (r *TerraformBlockFileNameRule) checkFile(runner tflint.Runner, filename string, file *hcl.File) error {
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
		if basename == "variables.tf" {
			if block.Type != "variable" {
				if block.Type == "output" {
					blockManager := NewBlockManager(runner)
					fixFunc := blockManager.CreateFixFunction(block, "outputs.tf")
					if err := runner.EmitIssueWithFix(
						r,
						fmt.Sprintf("Output block should be declared in outputs.tf, not in %s", basename),
						blockRange,
						fixFunc,
					); err != nil {
						return err
					}
				} else if block.Type == "locals" {
					blockManager := NewBlockManager(runner)
					fixFunc := blockManager.CreateFixFunction(block, "locals.tf")
					if err := runner.EmitIssueWithFix(
						r,
						fmt.Sprintf("Locals block should be declared in locals.tf, not in %s", basename),
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
				} else {
					if err := runner.EmitIssue(
						r,
						fmt.Sprintf("Only variable blocks should be declared in variables.tf, found %s block", block.Type),
						blockRange,
					); err != nil {
						return err
					}
				}
			}
		} else if basename == "outputs.tf" {
			if block.Type != "output" {
				if block.Type == "variable" {
					blockManager := NewBlockManager(runner)
					fixFunc := blockManager.CreateFixFunction(block, "variables.tf")
					if err := runner.EmitIssueWithFix(
						r,
						fmt.Sprintf("Variable block should be declared in variables.tf, not in %s", basename),
						blockRange,
						fixFunc,
					); err != nil {
						return err
					}
				} else if block.Type == "locals" {
					blockManager := NewBlockManager(runner)
					fixFunc := blockManager.CreateFixFunction(block, "locals.tf")
					if err := runner.EmitIssueWithFix(
						r,
						fmt.Sprintf("Locals block should be declared in locals.tf, not in %s", basename),
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
				} else {
					if err := runner.EmitIssue(
						r,
						fmt.Sprintf("Only output blocks should be declared in outputs.tf, found %s block", block.Type),
						blockRange,
					); err != nil {
						return err
					}
				}
			}
		} else if basename == "locals.tf" {
			if block.Type != "locals" {
				if block.Type == "variable" {
					blockManager := NewBlockManager(runner)
					fixFunc := blockManager.CreateFixFunction(block, "variables.tf")
					if err := runner.EmitIssueWithFix(
						r,
						fmt.Sprintf("Variable block should be declared in variables.tf, not in %s", basename),
						blockRange,
						fixFunc,
					); err != nil {
						return err
					}
				} else if block.Type == "output" {
					blockManager := NewBlockManager(runner)
					fixFunc := blockManager.CreateFixFunction(block, "outputs.tf")
					if err := runner.EmitIssueWithFix(
						r,
						fmt.Sprintf("Output block should be declared in outputs.tf, not in %s", basename),
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
				} else {
					if err := runner.EmitIssue(
						r,
						fmt.Sprintf("Only locals blocks should be declared in locals.tf, found %s block", block.Type),
						blockRange,
					); err != nil {
						return err
					}
				}
			}
		} else {
			// Handle non-special files
			switch block.Type {
			case "variable":
				blockManager := NewBlockManager(runner)
				fixFunc := blockManager.CreateFixFunction(block, "variables.tf")

				if err := runner.EmitIssueWithFix(
					r,
					fmt.Sprintf("Variable block should be declared in variables.tf, not in %s", basename),
					blockRange,
					fixFunc,
				); err != nil {
					return err
				}
			case "output":
				blockManager := NewBlockManager(runner)
				fixFunc := blockManager.CreateFixFunction(block, "outputs.tf")

				if err := runner.EmitIssueWithFix(
					r,
					fmt.Sprintf("Output block should be declared in outputs.tf, not in %s", basename),
					blockRange,
					fixFunc,
				); err != nil {
					return err
				}
			case "locals":
				blockManager := NewBlockManager(runner)
				fixFunc := blockManager.CreateFixFunction(block, "locals.tf")

				if err := runner.EmitIssueWithFix(
					r,
					fmt.Sprintf("Locals block should be declared in locals.tf, not in %s", basename),
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
			}

			// Check if this is a resource-specific file that contains non-matching blocks
			if r.isResourceFile(basename) {
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
		}
	}

	return nil
}

// isResourceFile determines if a filename follows the resource naming pattern
// Resource files should be named like "aws_instance.tf", "azurerm_virtual_machine.tf", etc.
// This is different from special files like "variables.tf", "outputs.tf", "locals.tf", "main.tf"
func (r *TerraformBlockFileNameRule) isResourceFile(basename string) bool {
	if !strings.HasSuffix(basename, ".tf") {
		return false
	}

	filename := strings.TrimSuffix(basename, ".tf")

	// Exclude special files
	specialFiles := []string{"variables", "outputs", "locals", "main", "providers", "versions", "terraform"}
	for _, special := range specialFiles {
		if filename == special {
			return false
		}
	}

	// Check if it contains underscores (indicating resource type pattern like aws_instance, azurerm_vm)
	return strings.Contains(filename, "_")
}
