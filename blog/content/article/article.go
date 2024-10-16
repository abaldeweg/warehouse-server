package article

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/abaldeweg/warehouse-server/blog/content"
)

// GetArticle retrieves an article from the filesystem safely.
func GetArticle(path string) (string, error) {
	fullPath := filepath.Join(content.GetContentRoot(), filepath.Clean("/"+path))

	if !strings.HasPrefix(fullPath, content.GetContentRoot()) {
		return "", fmt.Errorf("error: path traversal attempt detected")
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("error reading article: %w", err)
	}

	return string(data), nil
}
