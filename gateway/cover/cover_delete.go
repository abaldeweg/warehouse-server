package cover

import (
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// DeleteCover removes the cover images associated with the given book ID.
func DeleteCover(bookID uuid.UUID) {
	path, _ := getPath()
	sizes := []string{"l", "m", "s"}
	for _, size := range sizes {
		filename := filepath.Join(path, bookID.String()+"-"+size+".jpg")
		if err := os.Remove(filename); err == nil {
			log.Printf("failed to delete file: %s", filename)
		}
	}
}
