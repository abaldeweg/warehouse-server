package storage

import (
	"os"
	"path/filepath"
)

// FilesystemStorage implements storage for the filesystem.
type FilesystemStorage struct {
	Path     string
	FileName string
}

// save writes data to a file.
func (s *FilesystemStorage) save(data []byte) error {
	fullPath := filepath.Join(s.Path, s.FileName)

	return os.WriteFile(fullPath, data, 0644)
}

// load reads data from a file.
func (s *FilesystemStorage) load() ([]byte, error) {
	fullPath := filepath.Join(s.Path, s.FileName)

	if _, err := os.Stat(fullPath); err != nil {
		if os.IsNotExist(err) {
			return []byte("[]"), nil
		}
		return nil, err
	}

	return os.ReadFile(fullPath)
}

// remove deletes the file from the filesystem.
func (s *FilesystemStorage) remove() error {
	fullPath := filepath.Join(s.Path, s.FileName)

	return os.Remove(fullPath)
}

// exists checks if the file exists in the filesystem.
func (s *FilesystemStorage) exists() (bool, error) {
	fullPath := filepath.Join(s.Path, s.FileName)
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
