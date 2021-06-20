package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	configFilePath string
	apiKey         string
)

func init() {
	flag.StringVar(
		&configFilePath, "config", "config.yml", "Path to yaml-formatted configuration file",
	)
	flag.StringVar(
		&apiKey, "key", lookupEnvDefault("NOTION_AUTO_ARCHIVE_API_KEY", ""),
		"Notion.so API key (can also be set using NOTION_AUTO_ARCHIVE_API_KEY environment variable)",
	)

	flag.Parse()
}

func runBoard(board *TaskBoard, handler *NotionHandler) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	tasks, err := handler.GetTasksOnBoard(ctx, board)
	if err != nil {
		log.WithError(err).WithField("id", board.ID).Error("Could not get tasks from board")
		return
	}

	if len(tasks) < 1 {
		log.WithField("id", board.ID).Debug("Nothing to move. Exiting execution...")
		return
	}

	cutoff := time.Now().Add(-board.archiveAfterParsed.Duration())
	tasks = FilterOnLastEditedTime(cutoff, tasks)

	err = handler.MoveTasks(ctx, board, tasks...)
	if err != nil {
		log.WithError(err).WithField("id", board.ID).Errorf(
			"Could not move tasks from:%s to:%s", board.Selectors.SourceColumn, board.Selectors.TargetColumn,
		)
		return
	}
}

func run(config *Configuration) {
	sleepInterval := FindRefreshInterval(config.Boards)
	handler := NewNotionHandler(config.APIKey)

	for {
		for i, board := range config.Boards {
			log.Infof("Running archive task on board %d", i+1)
			runBoard(board, handler)
		}

		log.Debugf("Going to sleep for %v", sleepInterval)
		time.Sleep(sleepInterval)
	}
}

func main() {
	config, err := LoadConfigurationFile(configFilePath)
	if err != nil {
		log.WithError(err).Fatal("Could not load configuration file from disk")
	}

	log.WithField("configFile", configFilePath).Info("Application has successfully loaded configuration file")

	if apiKey != "" {
		log.Info("Preferring CLI/Environment set API key")
		config.APIKey = apiKey
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go run(config)
	log.Info("Application is now running")

	<-sig
	fmt.Printf("\n")

	log.Info("Application shutting down. Goodbye...")
}
