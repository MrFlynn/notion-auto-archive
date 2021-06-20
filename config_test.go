package main

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var (
	validConfig = `
apiKey: test_key
boards:
  - id: test_id
    archiveAfter: 24h
    selectors:
      columnName: Status
      sourceColumn: Completed
      targetColumn: Archived 
`
	validConfigDefaults = `
apiKey: test_key
boards:
  - id: test_id
`
	invalidConfigMissingBoards = `
apiKey: test_key
`
	invalidConfigMissingBoardID = `
apiKey: test_key
boards:
  - archiveAfter: 24h
    selectors:
      columnName: Status
      sourceColumn: Completed
      targetColumn: Archived 
`

	expectedConfigurationFile = &Configuration{
		APIKey: "test_key",
		Boards: []*TaskBoard{
			{
				ID:           "test_id",
				ArchiveAfter: "24h",
				Selectors: struct {
					ColumnName   string "yaml:\"columnName\""
					SourceColumn string "yaml:\"sourceColumn\""
					TargetColumn string "yaml:\"targetColumn\""
				}{
					ColumnName:   "Status",
					SourceColumn: "Completed",
					TargetColumn: "Archived",
				},
			},
		},
	}

	compareOptions = cmpopts.IgnoreFields(TaskBoard{}, "archiveAfterParsed")
)

func temporaryFileTestWrapper(t *testing.T, contents string, testFunc func(t *testing.T, filename string)) {
	tempFile, err := ioutil.TempFile(os.TempDir(), "notion-auto-archive-")
	if err != nil {
		t.Errorf("Got unexpected tempfile error: %s", err)
	}

	defer os.Remove(tempFile.Name())

	if _, err := tempFile.WriteString(contents); err != nil {
		t.Errorf("Got error writing data to tempfile: %s", err)
	}

	if err := tempFile.Close(); err != nil {
		t.Errorf("Got error trying to close tempfile: %s", err)
	}

	testFunc(t, tempFile.Name())
}

func TestValidConfig(t *testing.T) {
	temporaryFileTestWrapper(t, validConfig, func(t *testing.T, filename string) {
		config, err := LoadConfigurationFile(filename)
		if err != nil {
			t.Errorf("Got unexpected error: %s", err)
		}

		if !cmp.Equal(config, expectedConfigurationFile, compareOptions) {
			t.Errorf("Got unexpected diff:\n%s", cmp.Diff(config, expectedConfigurationFile, compareOptions))
		}

		if parsed := config.Boards[0].archiveAfterParsed.Duration(); parsed != 24*time.Hour {
			t.Errorf("Expected parsed time to field to be %v, got %v", 24*time.Hour, parsed)
		}
	})
}

func TestValidConfigWithDefaults(t *testing.T) {
	temporaryFileTestWrapper(t, validConfigDefaults, func(t *testing.T, filename string) {
		config, err := LoadConfigurationFile(filename)
		if err != nil {
			t.Errorf("Got unexpected error: %s", err)
		}

		if !cmp.Equal(config, expectedConfigurationFile, compareOptions) {
			t.Errorf("Got unexpected diff:\n%s", cmp.Diff(config, expectedConfigurationFile, compareOptions))
		}

		if parsed := config.Boards[0].archiveAfterParsed.Duration(); parsed != 24*time.Hour {
			t.Errorf("Expected parsed time to field to be %v, got %v", 24*time.Hour, parsed)
		}
	})
}

func TestMissingTaskBoards(t *testing.T) {
	temporaryFileTestWrapper(t, invalidConfigMissingBoards, func(t *testing.T, filename string) {
		_, err := LoadConfigurationFile(filename)
		if err == nil {
			t.Error("Expected error, got nil")
		}

		if err.Error() != "at least one task board must be specified" {
			t.Errorf("Expected error to be 'at least one task board must be specified', but got '%s'", err)
		}
	})
}

func TestMissingTaskBoardID(t *testing.T) {
	temporaryFileTestWrapper(t, invalidConfigMissingBoardID, func(t *testing.T, filename string) {
		_, err := LoadConfigurationFile(filename)
		if err == nil {
			t.Error("Expected error, got nil")
		}

		if err.Error() != "empty ID for board 1" {
			t.Errorf("Expected error to be 'empty ID for board 1', but got '%s'", err)
		}
	})
}
