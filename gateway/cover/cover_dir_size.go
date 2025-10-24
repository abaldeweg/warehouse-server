package cover

import (
	"os"
	"path/filepath"
)

// GetSize calculates the total size of all files in the specified directory.
func GetSize() (int64, error) {
	dir, err := getPath()
	if err != nil {
		return 0, err
	}

	var total int64 = 0
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() {
			total += info.Size()
		}
		return nil
	}); err != nil {
		return total, err
	}

	return total, nil
}
