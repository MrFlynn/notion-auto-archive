package main

import (
	"testing"
	"time"

	"github.com/dstotijn/go-notion"
	"github.com/hako/durafmt"
)

func TestFindRefreshIntervalBase(t *testing.T) {
	interval := FindRefreshInterval([]*TaskBoard{})
	if interval != 0 {
		t.Errorf("Expected 0 duration, unitless time, but got %v", interval)
	}
}

func TestFindRefreshIntervalSingle(t *testing.T) {
	boards := []*TaskBoard{
		{
			archiveAfterParsed: durafmt.Parse(10 * time.Second),
		},
	}

	interval := FindRefreshInterval(boards)
	if interval != 10*time.Second {
		t.Errorf("Expected to get %v, but got %v", 10*time.Second, interval)
	}
}

func TestFindRefreshIntervalMulti(t *testing.T) {
	boards := []*TaskBoard{
		{
			archiveAfterParsed: durafmt.Parse(10 * time.Second),
		},
		{
			archiveAfterParsed: durafmt.Parse(20 * time.Second),
		},
		{
			archiveAfterParsed: durafmt.Parse(15 * time.Second),
		},
	}

	interval := FindRefreshInterval(boards)
	if interval != 5*time.Second {
		t.Errorf("Expected to get %v, but got %v", 5*time.Second, interval)
	}
}

func TestFilterOnLastEditedTime(t *testing.T) {
	tasks := []notion.Page{
		{
			LastEditedTime: time.Now().Add(-60 * time.Minute),
			ID:             "1",
		},
		{
			LastEditedTime: time.Now().Add(-90 * time.Minute),
			ID:             "2",
		},
		{
			LastEditedTime: time.Now().Add(-30 * time.Minute),
			ID:             "3",
		},
	}

	filtered := FilterOnLastEditedTime(time.Now().Add(-70*time.Minute), tasks)
	if len(filtered) != 1 {
		t.Errorf("Expected to get 1 tasks, got %d", len(filtered))
	}

	if filtered[0].ID != "2" {
		t.Errorf("Expected ID to be 2, got %s", filtered[0].ID)
	}
}
