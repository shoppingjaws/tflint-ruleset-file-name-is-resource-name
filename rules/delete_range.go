package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
)

// calculateDeleteRange calculates the range to delete including the entire block and surrounding whitespace
func (bm *BlockMover) calculateDeleteRange(block *hclext.Block) (hcl.Range, error) {
	files, err := bm.runner.GetFiles()
	if err != nil {
		return hcl.Range{}, fmt.Errorf("failed to get files: %w", err)
	}

	sourceFile, exists := files[block.DefRange.Filename]
	if !exists {
		return hcl.Range{}, fmt.Errorf("source file not found: %s", block.DefRange.Filename)
	}

	content := string(sourceFile.Bytes)
	lines := strings.Split(content, "\n")

	// ブロックの開始位置
	startPos := block.DefRange.Start
	startLine := startPos.Line - 1 // 0-based indexing

	// より簡単なアプローチ: 行ベースで削除範囲を決定
	// 1. ブロック開始行を見つける
	// 2. 対応する閉じ括弧の行を見つける
	// 3. 可能な場合は前後の空行も含める

	// 開始行がブロック定義行
	blockStartLine := startLine

	// 閉じ括弧を含む行を見つける
	braceCount := 0
	blockEndLine := startLine
	foundStart := false

	for lineNum := startLine; lineNum < len(lines); lineNum++ {
		line := lines[lineNum]
		
		for _, char := range line {
			if char == '{' {
				foundStart = true
				braceCount++
			} else if char == '}' && foundStart {
				braceCount--
				if braceCount == 0 {
					blockEndLine = lineNum
					goto found
				}
			}
		}
	}

found:
	if !foundStart || braceCount != 0 {
		return hcl.Range{}, fmt.Errorf("could not find matching braces for block")
	}

	// 削除範囲を設定（行全体を削除）
	deleteStartLine := blockStartLine
	deleteEndLine := blockEndLine + 1 // 次の行の開始まで（つまり改行を含む）

	// 前の行が空行の場合は含める
	if deleteStartLine > 0 && strings.TrimSpace(lines[deleteStartLine-1]) == "" {
		deleteStartLine--
	}

	// 次の行が空行の場合は含める
	if deleteEndLine < len(lines) && strings.TrimSpace(lines[deleteEndLine]) == "" {
		deleteEndLine++
	}

	// 範囲を行・列で表現
	return hcl.Range{
		Filename: block.DefRange.Filename,
		Start: hcl.Pos{
			Line:   deleteStartLine + 1, // 1-based
			Column: 1,
		},
		End: hcl.Pos{
			Line:   deleteEndLine + 1, // 1-based
			Column: 1,
		},
	}, nil
}