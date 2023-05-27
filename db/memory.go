package db

import (
	"context"
	"os"

	"github.com/oleiade/reflections"
	"github.com/yakiroren/dss-common/models"
)

type MemoryDataStore struct {
	storage map[string]models.FileMetadata
}

func NewMemoryDataStore() (*MemoryDataStore, error) {
	db := &MemoryDataStore{}
	db.storage = make(map[string]models.FileMetadata)

	return db, nil
}

func (db *MemoryDataStore) WriteFile(_ context.Context, file models.FileMetadata) (string, error) {
	db.storage[file.Path] = file
	return file.Path, nil
}

func (db *MemoryDataStore) UpdateField(ctx context.Context, path interface{}, field string, value interface{}) error {
	file, found := db.storage[path.(string)]
	if !found {
		return os.ErrNotExist
	}

	reflections.SetField(&file, field, value)

	return nil
}

func (db *MemoryDataStore) AppendFragment(_ context.Context, name string, fragment models.Fragment) error {
	file, found := db.storage[name]
	if !found {
		return os.ErrNotExist
	}
	file.Fragments = append(file.Fragments, fragment)

	db.storage[name] = file

	return nil
}

func (db *MemoryDataStore) ListFiles(ctx context.Context) ([]models.FileMetadata, error) {
	var results []models.FileMetadata

	for _, fileMetadata := range db.storage {
		results = append(results, fileMetadata)
	}

	return results, nil
}
