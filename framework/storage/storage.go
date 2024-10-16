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
}

// Storage represents a storage object with configurable parameters.
type Storage struct {
	Type     StorageType
	Resource string // e.g. BucketName, DB Name
	Location string // e.g. Directory
	Name     string // e.g. Filename, Table Name
}

// NewStorage creates a new Storage instance with default values.
func NewStorage(resource, location, name string) *Storage {
	return &Storage{
		Type:     StorageTypeFilesystem,
		Resource: resource,
		Location: location,
		Name:     name,
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

// getStorage returns the appropriate storage implementation based on the provided type.
func (s *Storage) getStorage() persistence {
	switch s.Type {
	case StorageTypeCloud:
		return &cloudStorage{bucketName: s.Resource, directory: s.Location, name: s.Name}
	case StorageTypeFilesystem:
		fallthrough
	default:
		return &filesystemStorage{basePath: s.Location, name: s.Name}
	}
}
