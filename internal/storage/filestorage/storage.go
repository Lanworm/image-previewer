package filestorage

import (
	"fmt"
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
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (f FileStorage) GetFileList(folderPath string) ([]string, error) {
	// Проверяем существование папки, если нет - создаем
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.MkdirAll(folderPath, 0o755)
		if err != nil {
			return nil, err
		}
	}

	// Читаем содержимое папки
	fileInfos, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	// Создаем слайс для хранения имен файлов
	filenames := make([]string, 0, len(fileInfos))

	// Проходим по всем элементам в папке
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() { // Проверяем, не является ли текущий элемент папкой
			filenames = append(filenames, fileInfo.Name())
		}
	}

	return filenames, nil
}
