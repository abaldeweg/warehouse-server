package parser

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/abaldeweg/warehouse-server/logs/entity"
)

// ReadLogEntries reads log entries from log files and returns them as a slice of LogEntry.
func ReadLogEntries() ([]entity.LogEntry, error) {
	var logFiles []string
	var err error

	logFiles, err = listAndFilterLogFiles("data/source/*")
	if err != nil {
		return nil, err
	}
	fmt.Println("found files: ", logFiles)

	var entries []entity.LogEntry
	for _, logFile := range logFiles {
		fileEntries, err := parseLogFile(logFile)
		if err != nil {
			return nil, err
		}
		entries = append(entries, fileEntries...)
	}

	return entries, nil
}

// listAndFilterLogFiles lists all files matching the given pattern and filters them to include only log files with "access.log-" in their name and ".gz" extension.
func listAndFilterLogFiles(pattern string) ([]string, error) {
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	var logFiles []string
	threshold := time.Now().AddDate(0, 0, -5).Format("20060102")
	thresholdInt, err := strconv.Atoi(threshold)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if strings.Contains(file, "access.log-") && strings.HasSuffix(file, ".gz") {
			date := strings.TrimSuffix(strings.Split(file, "access.log-")[1], ".gz")
			dateInt, err := strconv.Atoi(date)
			if err != nil {
				return nil, err
			}
			fmt.Println("Extracted date:", dateInt)
			fmt.Println("Threshold:", thresholdInt)
			if dateInt >= thresholdInt {
				logFiles = append(logFiles, file)
			}
		}
	}
	return logFiles, nil
}

// parseLogFile parses a log file and returns a slice of LogEntry.
func parseLogFile(logFile string) ([]entity.LogEntry, error) {
	var entries []entity.LogEntry

	f, err := os.Open(logFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	scanner := bufio.NewScanner(gz)
	for scanner.Scan() {
		var entry entity.LogEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}
