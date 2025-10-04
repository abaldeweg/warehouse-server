package cover

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
)

func TestDeleteCover(t *testing.T) {
	bookID := uuid.New()
	path, _ := getPath()
	sizes := []string{"l", "m", "s"}

	for _, size := range sizes {
		filename := filepath.Join(path, bookID.String()+"-"+size+".jpg")
		if err := os.WriteFile(filename, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}

	DeleteCover(bookID)

	for _, size := range sizes {
		filename := filepath.Join(path, bookID.String()+"-"+size+".jpg")
		if _, err := os.Stat(filename); !os.IsNotExist(err) {
			t.Errorf("file %s was not deleted", filename)
		}
	}
}
