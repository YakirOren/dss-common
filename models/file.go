package models

import (
	"io/fs"
	"os"
	"time"
)

type FileMetadata struct {
	Id             interface{} `bson:"_id"`
	FileName       string      `bson:"name"`
	FileSize       int64       `bson:"size"`
	CurrentSize    int64       `bson:"currentSize"`
	IsDirectory    bool        `bson:"isDirectory"`
	Path           string      `bson:"path"`
	Fragments      []Fragment
	TotalFragments int
	IsHidden       bool
}

func NewFileMetadata(name string, path string, size int64, isDir bool, fragments []Fragment) *FileMetadata {
	return &FileMetadata{
		FileName:    name,
		FileSize:    size,
		IsDirectory: isDir,
		Path:        path,
		Fragments:   fragments,
	}
}

func (f FileMetadata) Name() string {
	return f.FileName
}

func (f FileMetadata) Size() int64 {
	return f.FileSize
}

func (f FileMetadata) IsDir() bool {
	return f.IsDirectory
}

func (f FileMetadata) Mode() fs.FileMode {
	return os.ModePerm
}

func (f FileMetadata) ModTime() time.Time {
	return time.Now()
}

func (f FileMetadata) Sys() interface{} {
	return nil
}
