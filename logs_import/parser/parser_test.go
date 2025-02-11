package parser

import (
	"compress/gzip"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestListAndFilterLogFiles(t *testing.T) {
	os.MkdirAll("data/logs", os.ModePerm)
	defer os.RemoveAll("data")

	testFiles := []string{
		"data/logs/access.log-" + time.Now().Format("20060102") + ".gz",
		"data/logs/access.log-" + time.Now().AddDate(0, 0, -1).Format("20060102") + ".gz",
		"data/logs/access.log-20230101.gz",
	}
	for _, file := range testFiles {
		f, _ := os.Create(file)
		f.Close()
	}

	files, err := listAndFilterLogFiles("data/logs/access.log-*.gz")

	assert.NoError(t, err)
	assert.Equal(t, 2, len(files))
}

func TestParseLogFile(t *testing.T) {
	os.MkdirAll("data/logs", os.ModePerm)
	defer os.RemoveAll("data")

	logFile := "data/logs/access.log-" + time.Now().Format("20060102") + ".gz"
	f, _ := os.Create(logFile)
	defer f.Close()

	gz := gzip.NewWriter(f)
	gz.Write([]byte(`{"msg": "Test"}`))
	gz.Close()

	entries, err := parseLogFile(logFile)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(entries))
	assert.Equal(t, "Test", entries[0].Msg)
}
