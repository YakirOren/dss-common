package db

import (
	"context"

	"github.com/yakiroren/dss-common/models"
)

type DataStore interface {
	WriteFile(ctx context.Context, file models.FileMetadata) (string, error)
	AppendFragment(ctx context.Context, path string, fragment models.Fragment) error
	GetMetadataByPath(ctx context.Context, path string) (*models.FileMetadata, bool)
	ListFiles(ctx context.Context, path string) ([]models.FileMetadata, error)
	UpdateField(ctx context.Context, id string, field string, value interface{}) error
	GetMetadataByID(ctx context.Context, id string) (*models.FileMetadata, bool)
	Delete(ctx context.Context, id string) bool
}
