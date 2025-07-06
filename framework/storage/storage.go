package storage

// StorageType represents the type of storage to use.
type StorageType string

// Constants for different storage types.
const (
	StorageTypeFilesystem StorageType = "filesystem"
	StorageTypeCloud      StorageType = "cloud"
)

// persistence defines the interface for saving and loading data.
type persistence interface {
	save(data []byte) error
	load() ([]byte, error)
	remove() error
	exists() (bool, error)
}

// Storage represents a storage object with configurable parameters.
type Storage struct {
	Type       StorageType
	FileSystem FilesystemStorage
	Cloud      CloudStorage
}

// NewStorage creates a new Storage instance with default values.
func NewStorage() *Storage {
	return &Storage{
		Type: StorageTypeFilesystem,
		FileSystem: FilesystemStorage{
			Path:     ".",
			FileName: "data.json",
		},
		Cloud: CloudStorage{
			BucketName: "storage-bucket",
			Path:       ".",
			FileName:   "data.json",
		},
	}
}

// Save stores data using the configured storage mechanism.
func (s *Storage) Save(content []byte) error {
	storage := s.getStorage()
	return storage.save(content)
}

// Load retrieves data using the configured storage mechanism.
func (s *Storage) Load() ([]byte, error) {
	storage := s.getStorage()
	return storage.load()
}

// Remove deletes data using the configured storage mechanism.
func (s *Storage) Remove() error {
	storage := s.getStorage()
	return storage.remove()
}

// Exists checks if the data exists in the configured storage mechanism.
func (s *Storage) Exists() (bool, error) {
	storage := s.getStorage()
	return storage.exists()
}

// getStorage returns the appropriate storage implementation based on the provided type.
func (s *Storage) getStorage() persistence {
	switch s.Type {
	case StorageTypeCloud:
		return &s.Cloud
	case StorageTypeFilesystem:
		fallthrough
	default:
		return &s.FileSystem
	}
}
