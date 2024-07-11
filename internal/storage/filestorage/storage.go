package filestorage

import (
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
)

type FileStorage struct {
	storagePath string
}

func (f FileStorage) Delete(id string) error {
	// TODO implement me
	panic(id)
}

func (f FileStorage) Set(item image.Image, id string) (bool, error) {
	if err := os.MkdirAll(f.storagePath, os.ModePerm); err != nil {
		return false, err
	}

	outputFile, cfErr := os.Create(filepath.Join(f.storagePath, id))
	if cfErr != nil {
		return false, cfErr
	}

	defer outputFile.Close()
	EncodeErr := jpeg.Encode(outputFile, item, nil)
	if EncodeErr != nil {
		return false, EncodeErr
	}
	return true, nil
}

func (f FileStorage) Get(id string) (image.Image, error) {
	file, err := os.Open(filepath.Join(f.storagePath, id))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func NewFileStorage(path string) *FileStorage {
	return &FileStorage{storagePath: path}
}
