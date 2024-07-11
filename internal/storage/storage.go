package storage

import (
	"image"
)

type Storage interface {
	Set(item image.Image, id string) error
	Get(id string) (image.Image, error)
	Delete(id string) error
	GetFileList(folderPath string) ([]string, error)
}
