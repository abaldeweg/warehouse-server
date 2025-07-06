package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"cloud.google.com/go/storage"
)

// CloudStorage implements storage for Google Cloud Storage.
type CloudStorage struct {
	BucketName string
	Path       string
	FileName   string
}

// save uploads data to Google Cloud Storage.
func (s *CloudStorage) save(data []byte) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("creating cloud storage client: %w", err)
	}
	defer client.Close()

	bucket := client.Bucket(s.BucketName)
	object := bucket.Object(filepath.Join(s.Path, s.FileName))
	wc := object.NewWriter(ctx)

	if _, err = wc.Write(data); err != nil {
		return fmt.Errorf("writing to cloud storage: %w", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("closing cloud storage writer: %w", err)
	}

	return nil
}

// load downloads data from Cloud.
func (s *CloudStorage) load() ([]byte, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		return nil, fmt.Errorf("creating cloud storage client: %w", err)
	}
	defer client.Close()

	bucket := client.Bucket(s.BucketName)
	object := bucket.Object(filepath.Join(s.Path, s.FileName))
	reader, err := object.NewReader(ctx)

	if err != nil {
		return nil, fmt.Errorf("creating cloud storage reader: %w", err)
	}
	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("reading from cloud storage: %w", err)
	}

	return data, nil
}

// remove deletes the object from Cloud.
func (s *CloudStorage) remove() error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("creating cloud storage client: %w", err)
	}
	defer client.Close()

	bucket := client.Bucket(s.BucketName)
	object := bucket.Object(filepath.Join(s.Path, s.FileName))

	return object.Delete(ctx)
}

// exists checks if the object exists in Cloud Storage.
func (s *CloudStorage) exists() (bool, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return false, fmt.Errorf("creating cloud storage client: %w", err)
	}
	defer client.Close()

	bucket := client.Bucket(s.BucketName)
	object := bucket.Object(filepath.Join(s.Path, s.FileName))
	_, err = object.Attrs(ctx)
	if err == storage.ErrObjectNotExist {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
