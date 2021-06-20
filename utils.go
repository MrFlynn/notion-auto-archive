package main

import (
	"os"
	"time"

	"github.com/dstotijn/go-notion"
)

// FindRefreshInterval finds the greatest common divisor of all ArchiveAfter fields in a slice of TaskBoards.
// This is useful for finding how often the program needs to wake up and initiate the archival process.
func FindRefreshInterval(boards []*TaskBoard) (interval time.Duration) {
	if len(boards) == 0 {
		return
	} else if len(boards) == 1 {
		interval = boards[0].archiveAfterParsed.Duration()
		return
	}

	gcd := func(x, y time.Duration) time.Duration {
		for y != 0 {
			t := x % y
			x = y
			y = t
		}

		return x
	}

	interval = gcd(boards[0].archiveAfterParsed.Duration(), boards[1].archiveAfterParsed.Duration())
	for i := 2; i < len(boards); i++ {
		interval = gcd(interval, boards[i].archiveAfterParsed.Duration())
	}

	return
}

// FilterOnLastEditedTime filters a list of tasks to find tasks that occured before some cutoff
// time.
func FilterOnLastEditedTime(beforeCutoff time.Time, tasks []notion.Page) (filteredTasks []notion.Page) {
	filteredTasks = make([]notion.Page, 0, len(tasks))

	for _, task := range tasks {
		if beforeCutoff.After(task.LastEditedTime) {
			filteredTasks = append(filteredTasks, task)
		}
	}

	return
}

func lookupEnvDefault(key, defaultValue string) (value string) {
	var ok bool

	value, ok = os.LookupEnv(key)
	if !ok {
		value = defaultValue
	}

	return
}
