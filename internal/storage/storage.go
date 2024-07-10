package storage

import (
	"image"
)

type IStorage interface {
	Set(item image.Image, id string) (bool, error)
	Get(id string) (image.Image, error)
	Delete(id string) error
}
