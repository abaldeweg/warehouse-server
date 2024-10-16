package home

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/abaldeweg/warehouse-server/blog/content"
	"gopkg.in/yaml.v3"
)

// GetHome returns the content index as a JSON string.
func GetHome() (string, error) {
	indexPath := filepath.Join(content.GetContentRoot(), "home.yaml")

	entries, err := loadHome(indexPath)
	if err != nil {
		return "", fmt.Errorf("error getting index: %w", err)
	}

	jsonData, err := json.Marshal(entries)
	if err != nil {
		return "", fmt.Errorf("error marshalling index to JSON: %w", err)
	}

	return string(jsonData), nil
}

// GetNewArticles returns the number of articles newer than the specified number of days.
func GetNewArticles(days int) (int, error) {
	indexPath := filepath.Join(content.GetContentRoot(), "home.yaml")

	entries, err := loadHome(indexPath)
	if err != nil {
		return 0, fmt.Errorf("error getting index: %w", err)
	}

	cutoffDate := time.Now().AddDate(0, 0, -days)

	newArticleCount := 0
	for _, entry := range entries {
		if entry.Date.After(cutoffDate) {
			newArticleCount++
		}
	}

	return newArticleCount, nil
}

// loadHome reads the content index file and returns its content as []IndexEntry.
// If the file doesn't exist, it creates a default one.
func loadHome(indexPath string) ([]content.Home, error) {
	if err := createHomeFileIfNotExists(indexPath); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, fmt.Errorf("error reading index file: %w", err)
	}

	var entries []content.Home
	if err := yaml.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("error unmarshalling index file: %w", err)
	}

	return entries, nil
}

// createHomeFileIfNotExists creates the index file if it doesn't exist.
func createHomeFileIfNotExists(indexPath string) error {
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		exampleContent := []byte("- slug: example\n  date: 2024-09-24\n  summary: This is an example article.\n")
		if err := os.WriteFile(indexPath, exampleContent, 0644); err != nil {
			return fmt.Errorf("error creating index file: %w", err)
		}
	}

	return nil
}
