package parser

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/abaldeweg/warehouse-server/logs/entity"
)

// ReadLogEntries reads log entries from log files and returns them as a slice of LogEntry.
func ReadLogEntries(optionalFileName ...string) ([]entity.LogEntry, error) {
	var logFiles []string
	var err error

	if len(optionalFileName) > 0 {
		logFiles = optionalFileName
    fmt.Println("passed log file: ", logFiles)
	} else {
		logFiles, err = listAndFilterLogFiles("data/source/*")
    fmt.Println("read files from fs")
		if err != nil {
			return nil, err
		}
	}

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
	for _, file := range files {
		if strings.Contains(file, "access.log-") && strings.HasSuffix(file, ".gz") {
			logFiles = append(logFiles, file)
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
