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

func NewFileStorage(path string) *FileStorage {
	return &FileStorage{storagePath: path}
}

func (f FileStorage) Set(item image.Image, id string) error {
	if err := os.MkdirAll(f.storagePath, os.ModePerm); err != nil {
		return err
	}

	outputFile, err := os.Create(filepath.Join(f.storagePath, id))
	if err != nil {
		return err
	}

	defer outputFile.Close()
	err = jpeg.Encode(outputFile, item, nil)
	if err != nil {
		return err
	}
	return nil
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

func (f FileStorage) Delete(id string) error {
	err := os.Remove(filepath.Join(f.storagePath, id))
	if err != nil {
		return err
	}
	return nil
}

func (f FileStorage) GetFileList(folderPath string) ([]string, error) {
	file, err := os.Open(folderPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfos, err := file.Readdir(-1)
	if err != nil {
		return nil, err
	}

	filenames := make([]string, 0, len(fileInfos))

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}
		filenames = append(filenames, fileInfo.Name())
	}

	return filenames, nil
}
