package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/hako/durafmt"
	"gopkg.in/yaml.v2"
)

// TaskBoard contains all settings for a single task board in Notion.
type TaskBoard struct {
	ID           string `yaml:"id"`
	ArchiveAfter string `yaml:"archiveAfter"`
	Selectors    struct {
		ColumnName   string `yaml:"columnName"`
		SourceColumn string `yaml:"sourceColumn"`
		TargetColumn string `yaml:"targetColumn"`
	} `yaml:"selectors"`

	archiveAfterParsed *durafmt.Durafmt
}

// Configuration is a struct represenation of the program configuration file.
type Configuration struct {
	APIKey string       `yaml:"apiKey"`
	Boards []*TaskBoard `yaml:"boards,flow"`
}

// LoadConfigurationFile takes the path to a configuration file, parses it, and sets default
// values.
func LoadConfigurationFile(filename string) (config *Configuration, err error) {
	var contents []byte

	contents, err = os.ReadFile(filename)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(contents, &config)
	if err != nil {
		return
	}

	if len(config.Boards) < 1 {
		err = errors.New("at least one task board must be specified")
	}

	for i, board := range config.Boards {
		if board.ID == "" {
			err = fmt.Errorf("empty ID for board %d", i+1)
			return
		}

		if board.ArchiveAfter == "" {
			board.ArchiveAfter = "24h"
		}

		if board.Selectors.ColumnName == "" {
			board.Selectors.ColumnName = "Status"
		}

		if board.Selectors.SourceColumn == "" {
			board.Selectors.SourceColumn = "Completed"
		}

		if board.Selectors.TargetColumn == "" {
			board.Selectors.TargetColumn = "Archived"
		}

		board.archiveAfterParsed, err = durafmt.ParseString(board.ArchiveAfter)
		if err != nil {
			return
		}
	}

	return
}
