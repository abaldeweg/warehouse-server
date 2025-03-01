package cover

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	Quality = 75
)

var Sizes = map[string]int{
	"l": 400,
	"m": 200,
	"s": 100,
}

func getPath() (string, error) {
	currentDir, _ := os.Getwd()
	uploadsDirPath := filepath.Join(currentDir, "uploads")

	if err := os.MkdirAll(uploadsDirPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create uploads directory")
	}

	return uploadsDirPath, nil
}
