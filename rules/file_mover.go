package rules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type BlockMover struct {
	runner tflint.Runner
}

func NewBlockMover(runner tflint.Runner) *BlockMover {
	return &BlockMover{runner: runner}
}

func (bm *BlockMover) MoveBlockToFile(blockType string, block *hclext.Block, targetFileName string) error {
	sourceFile := block.DefRange.Filename
	sourceDir := filepath.Dir(sourceFile)
	targetFile := filepath.Join(sourceDir, targetFileName)

	return bm.moveBlock(block, sourceFile, targetFile)
}

func (bm *BlockMover) moveBlock(block *hclext.Block, sourceFile, targetFile string) error {
	blockContent, err := bm.extractBlockContent(block)
	if err != nil {
		return fmt.Errorf("failed to extract block content: %w", err)
	}

	if err := bm.appendToTargetFile(targetFile, blockContent); err != nil {
		return fmt.Errorf("failed to append to target file: %w", err)
	}

	return nil
}

func (bm *BlockMover) extractBlockContent(block *hclext.Block) (string, error) {
	f := hclwrite.NewEmptyFile()
	body := f.Body()

	newBlock := body.AppendNewBlock(block.Type, block.Labels)
	newBlockBody := newBlock.Body()

	for name, attr := range block.Body.Attributes {
		if attr.Expr != nil {
			val, diags := attr.Expr.Value(nil)
			if diags.HasErrors() {
				return "", fmt.Errorf("failed to evaluate attribute %s: %s", name, diags.Error())
			}
			newBlockBody.SetAttributeValue(name, val)
		}
	}

	for _, nestedBlock := range block.Body.Blocks {
		nestedNewBlock := newBlockBody.AppendNewBlock(nestedBlock.Type, nestedBlock.Labels)
		if err := bm.copyBlockBody(nestedBlock.Body, nestedNewBlock.Body()); err != nil {
			return "", fmt.Errorf("failed to copy nested block: %w", err)
		}
	}

	return string(f.Bytes()), nil
}

func (bm *BlockMover) copyBlockBody(source *hclext.BodyContent, target *hclwrite.Body) error {
	for name, attr := range source.Attributes {
		if attr.Expr != nil {
			val, diags := attr.Expr.Value(nil)
			if diags.HasErrors() {
				return fmt.Errorf("failed to evaluate attribute %s: %s", name, diags.Error())
			}
			target.SetAttributeValue(name, val)
		}
	}

	for _, block := range source.Blocks {
		newBlock := target.AppendNewBlock(block.Type, block.Labels)
		if err := bm.copyBlockBody(block.Body, newBlock.Body()); err != nil {
			return fmt.Errorf("failed to copy nested block: %w", err)
		}
	}

	return nil
}

func (bm *BlockMover) appendToTargetFile(targetFile, content string) error {
	var existingContent []byte

	if _, statErr := os.Stat(targetFile); statErr == nil {
		var err error
		existingContent, err = os.ReadFile(targetFile)
		if err != nil {
			return fmt.Errorf("failed to read existing file: %w", err)
		}
	} else if !os.IsNotExist(statErr) {
		return fmt.Errorf("failed to check file existence: %w", statErr)
	}

	var newContent string
	if len(existingContent) > 0 {
		existingStr := string(existingContent)
		if !strings.HasSuffix(existingStr, "\n") {
			existingStr += "\n"
		}
		newContent = existingStr + "\n" + content
	} else {
		newContent = content
	}

	if err := os.WriteFile(targetFile, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write target file: %w", err)
	}

	return nil
}

func (bm *BlockMover) CreateFixFunction(blockType string, block *hclext.Block, targetFileName string) func(f tflint.Fixer) error {
	return func(f tflint.Fixer) error {
		if err := f.Remove(block.DefRange); err != nil {
			return err
		}

		return bm.MoveBlockToFile(blockType, block, targetFileName)
	}
}

func (bm *BlockMover) CreateFixFunctionForTestDir(blockType string, block *hclext.Block, targetFileName, testDir string) func(f tflint.Fixer) error {
	return func(f tflint.Fixer) error {
		if err := f.Remove(block.DefRange); err != nil {
			return err
		}

		targetFile := filepath.Join(testDir, targetFileName)

		return bm.moveBlock(block, block.DefRange.Filename, targetFile)
	}
}