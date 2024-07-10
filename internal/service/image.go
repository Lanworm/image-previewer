package service

import (
	"fmt"
	lrucache "github.com/Lanworm/image-previewe/internal/cache"
	"github.com/Lanworm/image-previewe/internal/logger"
	"github.com/Lanworm/image-previewe/internal/storage"
	"github.com/nfnt/resize"
	"image"
	"io"
	"net/http"
)

type ImageService struct {
	logger  *logger.Logger
	storage storage.IStorage
	cache   lrucache.Cache
}

func NewImageService(
	logger *logger.Logger,
	storage storage.IStorage,
	cache lrucache.Cache,
) *ImageService {
	return &ImageService{
		logger:  logger,
		storage: storage,
		cache:   cache,
	}
}

type ImgParams struct {
	Width  uint `validate:"required,gt=0,lte=9999"`
	Height uint `validate:"required,gt=0,lte=9999"`
	URL    string
}

func (s *ImageService) ResizeImg(imgParams *ImgParams) (img image.Image, err error) {

	imageID := getURLHash(imgParams.URL)

	cachedImg, ok := s.cache.Get(lrucache.Key(imageID))
	if ok {
		fmt.Println("received from cache: ", imageID)
		return cachedImg, nil
	} else {
		response, err := http.Get(imgParams.URL)
		if err != nil {
			return nil, err
		}
		fmt.Println("downloaded from url: ", imgParams.URL)

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				s.logger.Error(err.Error())
				return
			}
		}(response.Body)

		sourceImg, _, err := image.Decode(response.Body)
		if err != nil {
			return nil, err
		}
		s.cache.Set(lrucache.Key(imageID), sourceImg)
		newImg := resize.Resize(imgParams.Width, imgParams.Height, sourceImg, resize.Lanczos3)

		_, err = s.storage.Set(newImg, imageID)
		if err != nil {
			return nil, err
		}

		return newImg, nil
	}
}
