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
