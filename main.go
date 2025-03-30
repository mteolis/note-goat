package main

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/mteolis/note-goat/internal/goat"
	"github.com/sqweek/dialog"
)

var (
	appName = "NoteGoat"
	version = "v1.0.0"
	logger  *slog.Logger
)

func main() {
	start := time.Now()
	filename := "logs/" + appName + "_" + version + "_" + start.Format("20060102_150405") + ".log"

	err := os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		log.Fatalf("Error creating log directory: %+v", err)
	}

	logFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()

	handler := slog.NewTextHandler(logFile, nil)
	logger = slog.New(handler)

	logger.Info("%s %s executing...\n", appName, version)
	log.Printf("%s %s executing...\n", appName, version)

	excelPath := sqweekInputExcel()
	promptPath := sqweekInputPrompt()

	goat.InitGoat(logger, excelPath, promptPath)
	goat.AddAISummary()

	calcExecutionTime(start)
}

func sqweekInputExcel() string {
	filePath, err := dialog.File().
		Filter("Excel Files", "xlsx", "xls", "csv").
		Title("Select an Excel file to NoteGoat").
		Load()
	if err != nil {
		logger.Error("Error selecting file: %+v\n", "err", err)
		log.Fatalf("Error selecting file: %+v\n", err)
	}
	return filePath
}

func sqweekInputPrompt() string {
	filePath, err := dialog.File().
		Filter("Text Files", "txt").
		Title("Select a Text file as NoteGoat prompt").
		Load()
	if err != nil {
		logger.Error("Error selecting file: %+v\n", "err", err)
		log.Fatalf("Error selecting file: %+v\n", err)
	}
	return filePath
}

func calcExecutionTime(start time.Time) {
	duration := time.Since(start)

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	logger.Info("%s %s executed successfully in %dh:%dm:%ds\n",
		"appName", appName, "version", version, "hours", hours, "minutes", minutes, "seconds", seconds)
	log.Printf("%s %s executed successfully in %dh:%dm:%ds\n", appName, version, hours, minutes, seconds)
}
