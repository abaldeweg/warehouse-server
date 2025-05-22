package cover

import (
	"encoding/base64"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

func ShowCover(size string, bookID uuid.UUID) string {
	path, _ := getPath()
	filename := filepath.Join(path, bookID.String()+"-"+size+".jpg")

	if _, err := os.Stat(filename); err != nil {
		filename = filepath.Join(path, "none.jpg")
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return ""
	}
	return "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(data)
}
