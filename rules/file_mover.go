package rules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	// ソースファイルから直接テキストを抽出
	files, err := bm.runner.GetFiles()
	if err != nil {
		return "", fmt.Errorf("failed to get files: %w", err)
	}

	sourceFile, exists := files[block.DefRange.Filename]
	if !exists {
		return "", fmt.Errorf("source file not found: %s", block.DefRange.Filename)
	}

	// HCLブロックの完全な範囲を取得（TypeRangeやBodyRangeを考慮）
	// DefRangeはブロック定義の開始のみを表すため、ブロック全体を取得する必要がある
	
	// 開始位置はDefRangeから
	startPos := block.DefRange.Start
	
	// 終了位置を見つけるため、BodyRangeまたは手動で探す
	content := string(sourceFile.Bytes)
	
	// ブロックが始まる行から、対応する閉じ括弧を見つける
	startLine := startPos.Line - 1
	startCol := startPos.Column - 1
	
	// 文字ベースで処理
	bytes := []byte(content)
	startOffset := 0
	
	// 開始位置を計算
	currentLine := 0
	for i, b := range bytes {
		if currentLine == startLine {
			startOffset = i + startCol
			break
		}
		if b == '\n' {
			currentLine++
		}
	}
	
	// 開始位置から閉じ括弧を探す
	braceCount := 0
	inBlock := false
	endOffset := startOffset
	
	for i := startOffset; i < len(bytes); i++ {
		char := bytes[i]
		
		if char == '{' {
			inBlock = true
			braceCount++
		} else if char == '}' && inBlock {
			braceCount--
			if braceCount == 0 {
				endOffset = i + 1
				break
			}
		}
	}
	
	if endOffset <= startOffset {
		return "", fmt.Errorf("could not find end of block")
	}
	
	return string(bytes[startOffset:endOffset]), nil
}

// copyBlockBody is no longer used since we extract text directly
// func (bm *BlockMover) copyBlockBody(source *hclext.BodyContent, target *hclwrite.Body) error {
//	// This function is deprecated in favor of direct text extraction
//	return nil
// }

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
		// シンプルなアプローチ: TFLint Plugin SDKの制限を回避
		// ファイル全体を読み取り、ブロックを除去して書き換える
		
		// ブロックの内容を抽出
		blockContent, err := bm.extractBlockContent(block)
		if err != nil {
			return fmt.Errorf("failed to extract block content: %w", err)
		}

		// ターゲットファイルに移動
		sourceFile := block.DefRange.Filename
		sourceDir := filepath.Dir(sourceFile)
		targetFile := filepath.Join(sourceDir, targetFileName)

		if err := bm.appendToTargetFile(targetFile, blockContent); err != nil {
			return fmt.Errorf("failed to append to target file: %w", err)
		}

		// ソースファイルからブロックを削除
		if err := bm.removeBlockFromSourceFile(block); err != nil {
			return fmt.Errorf("failed to remove block from source file: %w", err)
		}

		return nil
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