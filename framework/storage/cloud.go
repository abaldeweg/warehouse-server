package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"cloud.google.com/go/storage"
)

// cloudStorage implements storage for Google Cloud Storage.
type cloudStorage struct {
	bucketName string
	directory  string
	name       string
}

// save uploads data to Google Cloud Storage.
func (s *cloudStorage) save(data []byte) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("creating cloud storage client: %w", err)
	}
	defer client.Close()

	bucket := client.Bucket(s.bucketName)
	object := bucket.Object(filepath.Join(s.directory, s.name))
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
func (s *cloudStorage) load() ([]byte, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		return nil, fmt.Errorf("creating cloud storage client: %w", err)
	}
	defer client.Close()

	bucket := client.Bucket(s.bucketName)
	object := bucket.Object(filepath.Join(s.directory, s.name))
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
func (s *cloudStorage) remove() error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("creating cloud storage client: %w", err)
	}
	defer client.Close()

	bucket := client.Bucket(s.bucketName)
	object := bucket.Object(filepath.Join(s.directory, s.name))

	return object.Delete(ctx)
}
