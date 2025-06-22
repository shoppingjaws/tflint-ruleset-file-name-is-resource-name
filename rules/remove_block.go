package rules

import (
	"fmt"
	"os"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
)

// removeBlockFromSourceFile removes a block from the source file directly
func (bm *BlockMover) removeBlockFromSourceFile(block *hclext.Block) error {
	sourceFile := block.DefRange.Filename

	// ファイルを読み取り
	content, err := os.ReadFile(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	
	// ブロックの開始位置
	startPos := block.DefRange.Start
	startLine := startPos.Line - 1 // 0-based indexing

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
		return fmt.Errorf("could not find matching braces for block")
	}

	// 削除範囲を設定（行全体を削除）
	deleteStartLine := startLine
	deleteEndLine := blockEndLine + 1 // 次の行まで（改行を含む）

	// 前の行が空行の場合は含める
	if deleteStartLine > 0 && strings.TrimSpace(lines[deleteStartLine-1]) == "" {
		deleteStartLine--
	}

	// 次の行が空行の場合は含める
	if deleteEndLine < len(lines) && strings.TrimSpace(lines[deleteEndLine]) == "" {
		deleteEndLine++
	}

	// 範囲外チェック
	if deleteEndLine > len(lines) {
		deleteEndLine = len(lines)
	}

	// 新しいコンテンツを構築（削除範囲を除く）
	newLines := make([]string, 0, len(lines)-(deleteEndLine-deleteStartLine))
	newLines = append(newLines, lines[:deleteStartLine]...)
	newLines = append(newLines, lines[deleteEndLine:]...)

	// ファイルに書き戻し
	newContent := strings.Join(newLines, "\n")
	if err := os.WriteFile(sourceFile, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write modified source file: %w", err)
	}

	return nil
}