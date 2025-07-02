package rules

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// BlockManager handles block extraction, movement, and deletion operations
type BlockManager struct {
	runner     tflint.Runner
	workingDir string // Optional working directory for test scenarios
}

// NewBlockManager creates a new BlockManager instance
func NewBlockManager(runner tflint.Runner) *BlockManager {
	return &BlockManager{runner: runner}
}

// NewBlockManagerWithWorkingDir creates a new BlockManager instance with working directory
func NewBlockManagerWithWorkingDir(runner tflint.Runner, workingDir string) *BlockManager {
	return &BlockManager{runner: runner, workingDir: workingDir}
}

// MoveBlock moves a block from its current location to a target file
func (bm *BlockManager) MoveBlock(block *hclext.Block, targetFileName string) error {
	// Extract block content
	blockContent, err := bm.ExtractBlockText(block)
	if err != nil {
		return fmt.Errorf("failed to extract block content: %w", err)
	}

	// Determine target file path
	var targetFile string
	if bm.workingDir != "" {
		// Use working directory for test scenarios
		targetFile = filepath.Join(bm.workingDir, targetFileName)
	} else {
		// Use source directory for normal operations
		sourceDir := filepath.Dir(block.DefRange.Filename)
		targetFile = filepath.Join(sourceDir, targetFileName)
	}

	// Append to target file
	if err := bm.appendToTargetFile(targetFile, blockContent); err != nil {
		return fmt.Errorf("failed to append to target file: %w", err)
	}

	// Remove from source file
	if err := bm.RemoveBlockFromFile(block); err != nil {
		return fmt.Errorf("failed to remove block from source file: %w", err)
	}

	return nil
}

// ExtractBlockText extracts the complete text content of a block
func (bm *BlockManager) ExtractBlockText(block *hclext.Block) (string, error) {
	files, err := bm.runner.GetFiles()
	if err != nil {
		return "", fmt.Errorf("failed to get files: %w", err)
	}

	sourceFile, exists := files[block.DefRange.Filename]
	if !exists {
		return "", fmt.Errorf("source file not found: %s", block.DefRange.Filename)
	}

	content := string(sourceFile.Bytes)
	lines := strings.Split(content, "\n")

	// Use DefRange for precise extraction
	startLine := block.DefRange.Start.Line - 1 // Convert to 0-based
	if startLine >= len(lines) {
		return "", fmt.Errorf("start line %d is beyond file length %d", startLine+1, len(lines))
	}

	// For import, moved, removed, and check blocks, include preceding comment lines
	if block.Type == "import" || block.Type == "moved" || block.Type == "removed" || block.Type == "check" {
		// Look backwards for comment lines
		originalStartLine := startLine
		for i := startLine - 1; i >= 0; i-- {
			trimmedLine := strings.TrimSpace(lines[i])
			if trimmedLine == "" {
				// Empty line - keep looking if we haven't found the comment yet
				if i > 0 && i == originalStartLine - 1 {
					continue
				}
				// Found content before, now hit empty line - stop
				if startLine < originalStartLine {
					break
				}
			} else if strings.HasPrefix(trimmedLine, "#") {
				// Comment line - include it
				startLine = i
			} else {
				// Non-comment, non-empty line - stop
				break
			}
		}
	}


	// Find the end of the block by counting braces with proper string handling
	braceCount := 0
	foundStart := false
	endLine := startLine
	inString := false
	inComment := false

	for lineNum := startLine; lineNum < len(lines); lineNum++ {
		line := lines[lineNum]

		for i, char := range line {
			// Handle string literals to avoid counting braces within strings
			if char == '"' && !inComment {
				// Check if it's escaped
				escaped := false
				if i > 0 && line[i-1] == '\\' {
					// Count consecutive backslashes
					backslashCount := 0
					for j := i - 1; j >= 0 && line[j] == '\\'; j-- {
						backslashCount++
					}
					escaped = (backslashCount%2 == 1)
				}
				if !escaped {
					inString = !inString
				}
				continue
			}

			// Handle single-line comments
			if !inString && i < len(line)-1 && line[i] == '/' && line[i+1] == '/' {
				inComment = true
				break // Skip rest of line
			}

			// Skip processing if we're in a string or comment
			if inString || inComment {
				continue
			}

			if char == '{' {
				foundStart = true
				braceCount++
			} else if char == '}' && foundStart {
				braceCount--
				if braceCount == 0 {
					endLine = lineNum
					// Extract the complete block including the closing brace line
					blockLines := lines[startLine : endLine+1]
					return strings.Join(blockLines, "\n"), nil
				}
			}
		}

		// Reset comment flag at end of line
		inComment = false
	}

	if !foundStart || braceCount != 0 {
		return "", fmt.Errorf("could not find matching braces for block at line %d", startLine+1)
	}

	return "", fmt.Errorf("incomplete block extraction")
}

// RemoveBlockFromFile removes a block from its source file
func (bm *BlockManager) RemoveBlockFromFile(block *hclext.Block) error {
	sourceFile := block.DefRange.Filename

	content, err := os.ReadFile(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to read source file: %w", err)
	}

	// Extract the exact block content to match and remove
	blockContent, err := bm.ExtractBlockText(block)
	if err != nil {
		return fmt.Errorf("failed to extract block content for removal: %w", err)
	}

	originalContent := string(content)
	
	// Try to remove the exact block content
	newContent := bm.removeBlockByContent(originalContent, blockContent, block.Type)
	
	// If nothing was removed, try alternative methods
	if newContent == originalContent {
		// Fallback to line-based removal
		lines := strings.Split(originalContent, "\n")
		startLine, endLine, err := bm.findBlockLines(lines, block)
		if err != nil {
			// Block boundaries not found, might already be removed or modified
			return nil // Don't error out, just skip removal
		}

		// Verify we actually found a meaningful block to remove
		if startLine >= 0 && endLine >= startLine && startLine < len(lines) && endLine < len(lines) {
			// Include surrounding empty lines for cleaner removal
			deleteStart, deleteEnd := bm.expandDeletionRange(lines, startLine, endLine)

			// Construct new content without the block
			newLines := make([]string, 0, len(lines)-(deleteEnd-deleteStart))
			newLines = append(newLines, lines[:deleteStart]...)
			newLines = append(newLines, lines[deleteEnd:]...)
			newContent = strings.Join(newLines, "\n")
		}
	}

	// Write back to file only if content changed
	if newContent != originalContent {
		if err := os.WriteFile(sourceFile, []byte(newContent), 0644); err != nil {
			return fmt.Errorf("failed to write modified source file: %w", err)
		}
	}

	return nil
}

// CreateFixFunction creates a fix function for TFLint
func (bm *BlockManager) CreateFixFunction(block *hclext.Block, targetFileName string) func(f tflint.Fixer) error {
	return func(f tflint.Fixer) error {
		// Extract block content
		blockContent, err := bm.ExtractBlockText(block)
		if err != nil {
			return fmt.Errorf("failed to extract block content: %w", err)
		}

		// Check if this is an empty block that should be skipped
		if bm.isEmptyBlock(blockContent, block.Type) {
			// Just remove the empty block, don't move it
			if err := bm.RemoveBlockFromFile(block); err != nil {
				return fmt.Errorf("failed to remove empty block from source file: %w", err)
			}
			return nil
		}

		// Determine target file path
		var targetFile string
		if bm.workingDir != "" {
			// Use working directory for test scenarios
			targetFile = filepath.Join(bm.workingDir, targetFileName)
		} else {
			// Use source directory for normal operations
			sourceDir := filepath.Dir(block.DefRange.Filename)
			targetFile = filepath.Join(sourceDir, targetFileName)
		}

		// Append to target file
		if err := bm.appendToTargetFile(targetFile, blockContent); err != nil {
			return fmt.Errorf("failed to append to target file: %w", err)
		}

		// Remove block from source file
		if err := bm.RemoveBlockFromFile(block); err != nil {
			return fmt.Errorf("failed to remove block from source file: %w", err)
		}

		return nil
	}
}

// isRealFileSystem checks if we're dealing with actual files or test files
func (bm *BlockManager) isRealFileSystem(filename string) bool {
	// In test scenarios, filename will not exist on actual filesystem
	// In real scenarios, it will exist
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

// removeBlockByContent removes a block from content by matching the exact block text
func (bm *BlockManager) removeBlockByContent(content, blockContent, blockType string) string {
	// Normalize whitespace for better matching
	normalizedContent := strings.TrimSpace(content)
	normalizedBlock := strings.TrimSpace(blockContent)
	
	if normalizedBlock == "" {
		return content
	}
	
	// Try exact match first
	if strings.Contains(normalizedContent, normalizedBlock) {
		// Find the position and remove with surrounding whitespace
		beforeBlock := ""
		afterBlock := ""
		
		blockIndex := strings.Index(normalizedContent, normalizedBlock)
		if blockIndex >= 0 {
			beforeBlock = normalizedContent[:blockIndex]
			afterBlock = normalizedContent[blockIndex+len(normalizedBlock):]
			
			// Clean up extra newlines
			beforeBlock = strings.TrimRight(beforeBlock, "\n")
			afterBlock = strings.TrimLeft(afterBlock, "\n")
			
			if beforeBlock != "" && afterBlock != "" {
				return beforeBlock + "\n\n" + afterBlock
			} else if beforeBlock != "" {
				return beforeBlock + "\n"
			} else if afterBlock != "" {
				return afterBlock
			} else {
				return ""
			}
		}
	}
	
	// Try line-by-line matching for more flexible removal
	contentLines := strings.Split(content, "\n")
	blockLines := strings.Split(blockContent, "\n")
	
	if len(blockLines) == 0 {
		return content
	}
	
	// Find the starting position of the block
	for i := 0; i <= len(contentLines)-len(blockLines); i++ {
		match := true
		for j, blockLine := range blockLines {
			if i+j >= len(contentLines) || strings.TrimSpace(contentLines[i+j]) != strings.TrimSpace(blockLine) {
				match = false
				break
			}
		}
		
		if match {
			// Found the block, remove it
			newLines := make([]string, 0, len(contentLines)-len(blockLines))
			newLines = append(newLines, contentLines[:i]...)
			newLines = append(newLines, contentLines[i+len(blockLines):]...)
			
			// Clean up extra empty lines
			result := strings.Join(newLines, "\n")
			result = strings.ReplaceAll(result, "\n\n\n", "\n\n")
			return result
		}
	}
	
	return content
}

// isEmptyBlock checks if a block has no meaningful content
func (bm *BlockManager) isEmptyBlock(blockContent, blockType string) bool {
	trimmedContent := strings.TrimSpace(blockContent)
	if trimmedContent == "" {
		return true
	}

	lines := strings.Split(trimmedContent, "\n")
	hasRealContent := false
	
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		// Skip empty lines, braces, and block declaration lines
		if trimmedLine == "" || trimmedLine == "{" || trimmedLine == "}" {
			continue
		}
		
		// Skip block declaration lines
		if strings.HasPrefix(trimmedLine, blockType) {
			continue
		}
		
		// If we find any other content, it's not empty
		hasRealContent = true
		break
	}
	
	return !hasRealContent
}

// Helper methods

func (bm *BlockManager) calculateOffset(bytes []byte, line, col int) int {
	currentLine := 1 // HCL uses 1-based line numbers
	for i, b := range bytes {
		if currentLine == line {
			return i + col - 1 // Convert to 0-based column
		}
		if b == '\n' {
			currentLine++
		}
	}
	return len(bytes)
}

func (bm *BlockManager) findBlockEnd(bytes []byte, startOffset int) int {
	braceCount := 0
	foundFirstBrace := false
	inString := false
	inComment := false

	for i := startOffset; i < len(bytes); i++ {
		char := bytes[i]

		// Handle string literals to avoid counting braces within strings
		if char == '"' && !inComment {
			inString = !inString
			continue
		}

		// Handle comments
		if !inString && i < len(bytes)-1 && bytes[i] == '/' && bytes[i+1] == '/' {
			inComment = true
			continue
		}

		if inComment && char == '\n' {
			inComment = false
			continue
		}

		// Skip processing if we're in a string or comment
		if inString || inComment {
			continue
		}

		if char == '{' {
			foundFirstBrace = true
			braceCount++
		} else if char == '}' && foundFirstBrace {
			braceCount--
			if braceCount == 0 {
				// Find the end of the line after the closing brace, including the newline
				for j := i + 1; j < len(bytes); j++ {
					if bytes[j] == '\n' {
						return j + 1
					}
				}
				return i + 1
			}
		}
	}
	return startOffset
}

func (bm *BlockManager) findBlockLines(lines []string, block *hclext.Block) (int, int, error) {
	startPos := block.DefRange.Start
	startLine := startPos.Line - 1 // 0-based indexing

	if startLine < 0 || startLine >= len(lines) {
		return 0, 0, fmt.Errorf("start line %d is out of bounds", startLine+1)
	}

	// Verify the start line contains the expected block type
	expectedBlockStart := false
	startLineContent := strings.TrimSpace(lines[startLine])
	if strings.HasPrefix(startLineContent, block.Type) {
		expectedBlockStart = true
	}

	if !expectedBlockStart {
		// Try to find the block by searching around the expected line
		found := false
		for offset := -2; offset <= 2; offset++ {
			checkLine := startLine + offset
			if checkLine >= 0 && checkLine < len(lines) {
				lineContent := strings.TrimSpace(lines[checkLine])
				if strings.HasPrefix(lineContent, block.Type) {
					startLine = checkLine
					found = true
					break
				}
			}
		}
		if !found {
			return 0, 0, fmt.Errorf("could not find block starting with '%s' near line %d", block.Type, startLine+1)
		}
	}

	// For import, moved, removed, and check blocks, include preceding comment lines
	if block.Type == "import" || block.Type == "moved" || block.Type == "removed" || block.Type == "check" {
		// Look backwards for comment lines
		originalStartLine := startLine
		for i := startLine - 1; i >= 0; i-- {
			trimmedLine := strings.TrimSpace(lines[i])
			if trimmedLine == "" {
				// Empty line - keep looking if we're right before the block
				if i == originalStartLine - 1 {
					continue
				}
				// Found content before, now hit empty line - stop
				if startLine < originalStartLine {
					break
				}
			} else if strings.HasPrefix(trimmedLine, "#") {
				// Comment line - include it
				startLine = i
			} else {
				// Non-comment, non-empty line - stop
				break
			}
		}
	}


	// Find the end line by counting braces with proper string handling
	braceCount := 0
	blockEndLine := startLine
	foundStart := false
	inString := false
	inComment := false

	for lineNum := startLine; lineNum < len(lines); lineNum++ {
		line := lines[lineNum]

		for i, char := range line {
			// Handle string literals to avoid counting braces within strings
			if char == '"' && !inComment {
				// Check if it's escaped
				escaped := false
				if i > 0 && line[i-1] == '\\' {
					// Count consecutive backslashes
					backslashCount := 0
					for j := i - 1; j >= 0 && line[j] == '\\'; j-- {
						backslashCount++
					}
					escaped = (backslashCount%2 == 1)
				}
				if !escaped {
					inString = !inString
				}
				continue
			}

			// Handle single-line comments
			if !inString && i < len(line)-1 && line[i] == '/' && line[i+1] == '/' {
				inComment = true
				break // Skip rest of line
			}

			// Skip processing if we're in a string or comment
			if inString || inComment {
				continue
			}

			if char == '{' {
				foundStart = true
				braceCount++
			} else if char == '}' && foundStart {
				braceCount--
				if braceCount == 0 {
					blockEndLine = lineNum
					return startLine, blockEndLine, nil
				}
			}
		}

		// Reset comment flag at end of line
		inComment = false
	}

	if !foundStart || braceCount != 0 {
		return 0, 0, fmt.Errorf("could not find matching braces for block")
	}

	return startLine, blockEndLine, nil
}

func (bm *BlockManager) expandDeletionRange(lines []string, startLine, endLine int) (int, int) {
	deleteStart := startLine
	deleteEnd := endLine + 1 // Include the closing brace line

	// Include all preceding empty lines
	for deleteStart > 0 && strings.TrimSpace(lines[deleteStart-1]) == "" {
		deleteStart--
	}

	// Include all following empty lines
	for deleteEnd < len(lines) && strings.TrimSpace(lines[deleteEnd]) == "" {
		deleteEnd++
	}

	// Ensure we don't go out of bounds
	if deleteEnd > len(lines) {
		deleteEnd = len(lines)
	}

	return deleteStart, deleteEnd
}

func (bm *BlockManager) appendToTargetFile(targetFile, content string) error {
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

	// Skip adding empty or invalid blocks first
	trimmedContent := strings.TrimSpace(content)
	if trimmedContent == "" || strings.Count(trimmedContent, "{") != strings.Count(trimmedContent, "}") {
		// Content is empty or braces don't match, skip adding
		return nil
	}

	// Check if the block is effectively empty (only has block declaration with empty body)
	lines := strings.Split(trimmedContent, "\n")
	nonEmptyLines := 0
	hasRealContent := false
	blockDeclarationLine := ""
	
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if trimmedLine != "" && trimmedLine != "{" && trimmedLine != "}" {
			nonEmptyLines++
			// Identify the block declaration line
			if strings.HasPrefix(trimmedLine, "locals") ||
				strings.HasPrefix(trimmedLine, "output") ||
				strings.HasPrefix(trimmedLine, "variable") ||
				strings.HasPrefix(trimmedLine, "resource") {
				blockDeclarationLine = trimmedLine
			} else {
				// This is actual content beyond the block declaration
				hasRealContent = true
			}
		}
	}
	
	// If we only have a block declaration line and no real content, skip it
	if nonEmptyLines <= 1 || !hasRealContent {
		return nil
	}
	
	// Special case: Check for empty locals block pattern "locals {"
	if strings.HasPrefix(blockDeclarationLine, "locals") && !hasRealContent {
		return nil
	}

	// Check for duplicates to prevent adding the same block multiple times
	if len(existingContent) > 0 {
		existingStr := string(existingContent)
		// More precise duplicate check - look for the actual block content
		contentLines := strings.Split(trimmedContent, "\n")
		if len(contentLines) > 2 {
			// Check if the core content (excluding first and last line) already exists
			coreContent := strings.Join(contentLines[1:len(contentLines)-1], "\n")
			if strings.Contains(existingStr, strings.TrimSpace(coreContent)) {
				// Block already exists, skip adding
				return nil
			}
		}
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

// removeBlockFromRealFile removes a block from the actual file system
func (bm *BlockManager) removeBlockFromRealFile(block *hclext.Block) error {
	// First remove the specific block
	if err := bm.RemoveBlockFromFile(block); err != nil {
		return err
	}

	// Then clean up any empty blocks of the same type that might be left
	return bm.cleanupEmptyBlocksInFile(block.DefRange.Filename, block.Type)
}

// cleanupEmptyBlocksInFile removes empty blocks from a file
func (bm *BlockManager) cleanupEmptyBlocksInFile(filename, blockType string) error {
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file for cleanup: %w", err)
	}

	contentStr := string(content)
	original := contentStr

	// Define patterns for empty blocks based on block type
	switch blockType {
	case "locals":
		// Remove empty locals blocks: locals { }
		contentStr = removeEmptyBlockPattern(contentStr, `(?m)^\s*locals\s*\{\s*\}\s*\n?`)
	case "output":
		// Remove empty output blocks: output "name" { }
		contentStr = removeEmptyBlockPattern(contentStr, `(?m)^\s*output\s*"[^"]*"\s*\{\s*\}\s*\n?`)
	case "variable":
		// Remove empty variable blocks: variable "name" { }
		contentStr = removeEmptyBlockPattern(contentStr, `(?m)^\s*variable\s*"[^"]*"\s*\{\s*\}\s*\n?`)
	}

	// Only write back if content changed
	if contentStr != original {
		if err := os.WriteFile(filename, []byte(contentStr), 0644); err != nil {
			return fmt.Errorf("failed to write cleaned file: %w", err)
		}
	}

	return nil
}

// removeEmptyBlockPattern removes blocks matching the given regex pattern
func removeEmptyBlockPattern(content, pattern string) string {
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(content, "")
}

// calculateFullBlockRange calculates the complete range of a block including its body
func (bm *BlockManager) calculateFullBlockRange(block *hclext.Block) (hcl.Range, error) {
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

	// Use DefRange for precise extraction
	startLine := block.DefRange.Start.Line - 1 // Convert to 0-based
	if startLine >= len(lines) {
		return hcl.Range{}, fmt.Errorf("start line %d is beyond file length %d", startLine+1, len(lines))
	}

	// Find the end of the block by counting braces with proper string handling
	braceCount := 0
	foundStart := false
	endLine := startLine
	inString := false
	inComment := false

	for lineNum := startLine; lineNum < len(lines); lineNum++ {
		line := lines[lineNum]

		for i, char := range line {
			// Handle string literals to avoid counting braces within strings
			if char == '"' && !inComment {
				// Check if it's escaped
				escaped := false
				if i > 0 && line[i-1] == '\\' {
					// Count consecutive backslashes
					backslashCount := 0
					for j := i - 1; j >= 0 && line[j] == '\\'; j-- {
						backslashCount++
					}
					escaped = (backslashCount%2 == 1)
				}
				if !escaped {
					inString = !inString
				}
				continue
			}

			// Handle single-line comments
			if !inString && i < len(lines)-1 && line[i] == '/' && line[i+1] == '/' {
				inComment = true
				break // Skip rest of line
			}

			// Skip processing if we're in a string or comment
			if inString || inComment {
				continue
			}

			if char == '{' {
				foundStart = true
				braceCount++
			} else if char == '}' && foundStart {
				braceCount--
				if braceCount == 0 {
					endLine = lineNum
					break
				}
			}
		}

		// Reset comment flag at end of line
		inComment = false

		// If we found the end of the block, break out of the outer loop
		if foundStart && braceCount == 0 {
			break
		}
	}

	if !foundStart || braceCount != 0 {
		return hcl.Range{}, fmt.Errorf("could not find matching braces for block at line %d", startLine+1)
	}

	// Find the end of the closing brace line
	endOffset := bm.calculateOffset(sourceFile.Bytes, endLine+1, 1) // Start of next line
	if endLine+1 >= len(lines) {
		// If this is the last line, include the entire content
		endOffset = len(sourceFile.Bytes)
	}

	return hcl.Range{
		Filename: block.DefRange.Filename,
		Start: block.DefRange.Start,
		End: hcl.Pos{
			Line:   endLine + 1, // Convert back to 1-based
			Column: len(lines[endLine]) + 1,
			Byte:   endOffset,
		},
	}, nil
}
