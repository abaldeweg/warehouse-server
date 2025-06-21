package storage

import (
	"os"
	"path/filepath"
)

// filesystemStorage implements storage for the filesystem.
type filesystemStorage struct {
	basePath string
	name     string
}

// save writes data to a file.
func (s *filesystemStorage) save(data []byte) error {
	fullPath := filepath.Join(s.basePath, s.name)

	return os.WriteFile(fullPath, data, 0644)
}

// load reads data from a file.
func (s *filesystemStorage) load() ([]byte, error) {
	fullPath := filepath.Join(s.basePath, s.name)

	if _, err := os.Stat(fullPath); err != nil {
		if os.IsNotExist(err) {
			return []byte("[]"), nil
		}
		return nil, err
	}

	return os.ReadFile(fullPath)
}

// remove deletes the file from the filesystem.
func (s *filesystemStorage) remove() error {
	fullPath := filepath.Join(s.basePath, s.name)

	return os.Remove(fullPath)
}

// exists checks if the file exists in the filesystem.
func (s *filesystemStorage) exists() (bool, error) {
	fullPath := filepath.Join(s.basePath, s.name)
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
