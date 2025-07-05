package rules

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_BlockManager_AppendToFile(t *testing.T) {
	tests := []struct {
		name            string
		existingContent string
		newContent      string
		expectedResult  string
	}{
		{
			name:            "create new file",
			existingContent: "",
			newContent:      "variable \"example\" {\n  type = string\n}",
			expectedResult:  "variable \"example\" {\n  type = string\n}",
		},
		{
			name:            "append to existing file",
			existingContent: "variable \"existing\" {\n  type = string\n}",
			newContent:      "variable \"new\" {\n  type = number\n}",
			expectedResult:  "variable \"existing\" {\n  type = string\n}\n\nvariable \"new\" {\n  type = number\n}",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Use project testdata directory for easy verification
			testDir := filepath.Join("testdata", "block_manager", t.Name())
			if err := os.MkdirAll(testDir, 0755); err != nil {
				t.Fatalf("Failed to create test directory: %s", err)
			}

			// Only clean up if KEEP_TEST_FILES is not set
			if os.Getenv("KEEP_TEST_FILES") == "" {
				defer os.RemoveAll(filepath.Join("testdata", "block_manager"))
			} else {
				t.Logf("Test files preserved in: %s", testDir)
			}

			targetFile := filepath.Join(testDir, "variables.tf")

			if test.existingContent != "" {
				if err := os.WriteFile(targetFile, []byte(test.existingContent), 0644); err != nil {
					t.Fatalf("Failed to write existing content: %s", err)
				}
			}

			blockManager := NewBlockManagerWithWorkingDir(nil, testDir)
			if err := blockManager.appendToTargetFile(targetFile, test.newContent); err != nil {
				t.Fatalf("Failed to append to target file: %s", err)
			}

			actualContent, err := os.ReadFile(targetFile)
			if err != nil {
				t.Fatalf("Failed to read result file: %s", err)
			}

			if string(actualContent) != test.expectedResult {
				t.Errorf("Content mismatch\nExpected: %q\nActual: %q", test.expectedResult, string(actualContent))
			}
		})
	}
}

func Test_BlockManager_DeleteFileIfEmpty(t *testing.T) {
	tests := []struct {
		name           string
		fileContent    string
		shouldDelete   bool
		expectedExists bool
	}{
		{
			name:           "delete file with only whitespace",
			fileContent:    "   \n\n  \t\n",
			shouldDelete:   true,
			expectedExists: false,
		},
		{
			name:           "delete file with only comments",
			fileContent:    "# This is a comment\n// Another comment\n\n",
			shouldDelete:   true,
			expectedExists: false,
		},
		{
			name:           "delete file with whitespace and comments",
			fileContent:    "   \n# Comment\n  \n// Another comment\n\t\n",
			shouldDelete:   true,
			expectedExists: false,
		},
		{
			name:           "keep file with actual content",
			fileContent:    "variable \"example\" {\n  type = string\n}",
			shouldDelete:   false,
			expectedExists: true,
		},
		{
			name:           "keep file with content and comments",
			fileContent:    "# This is a comment\nvariable \"example\" {\n  type = string\n}",
			shouldDelete:   false,
			expectedExists: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Use project testdata directory for easy verification
			testDir := filepath.Join("testdata", "block_manager", t.Name())
			if err := os.MkdirAll(testDir, 0755); err != nil {
				t.Fatalf("Failed to create test directory: %s", err)
			}

			// Only clean up if KEEP_TEST_FILES is not set
			if os.Getenv("KEEP_TEST_FILES") == "" {
				defer os.RemoveAll(filepath.Join("testdata", "block_manager"))
			} else {
				t.Logf("Test files preserved in: %s", testDir)
			}

			testFile := filepath.Join(testDir, "test.tf")

			// Write test content
			if err := os.WriteFile(testFile, []byte(test.fileContent), 0644); err != nil {
				t.Fatalf("Failed to write test file: %s", err)
			}

			blockManager := NewBlockManagerWithWorkingDir(nil, testDir)
			if err := blockManager.deleteFileIfEmpty(testFile); err != nil {
				t.Fatalf("Failed to check/delete empty file: %s", err)
			}

			// Check if file exists
			_, err := os.Stat(testFile)
			fileExists := !os.IsNotExist(err)

			if fileExists != test.expectedExists {
				t.Errorf("File existence mismatch. Expected exists: %v, Actual exists: %v", test.expectedExists, fileExists)
			}
		})
	}
}

func Test_BlockManager_IsFileEmpty(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "empty string",
			content:  "",
			expected: true,
		},
		{
			name:     "only whitespace",
			content:  "   \n\n  \t\n",
			expected: true,
		},
		{
			name:     "only hash comments",
			content:  "# This is a comment\n# Another comment\n",
			expected: true,
		},
		{
			name:     "only slash comments",
			content:  "// This is a comment\n// Another comment\n",
			expected: true,
		},
		{
			name:     "mixed comments and whitespace",
			content:  "   \n# Comment\n  \n// Another comment\n\t\n",
			expected: true,
		},
		{
			name:     "contains actual content",
			content:  "variable \"example\" {\n  type = string\n}",
			expected: false,
		},
		{
			name:     "content with comments",
			content:  "# This is a comment\nvariable \"example\" {\n  type = string\n}",
			expected: false,
		},
		{
			name:     "single meaningful line",
			content:  "terraform {\n}",
			expected: false,
		},
	}

	blockManager := NewBlockManager(nil)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := blockManager.isFileEmpty(test.content)
			if actual != test.expected {
				t.Errorf("isFileEmpty(%q) = %v, expected %v", test.content, actual, test.expected)
			}
		})
	}
}
