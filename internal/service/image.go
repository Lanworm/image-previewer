package service

import (
	"context"
	"fmt"
	"image"
	"net/http"
	"time"

	lrucache "github.com/Lanworm/image-previewe/internal/cache"
	"github.com/Lanworm/image-previewe/internal/logger"
	"github.com/Lanworm/image-previewe/internal/storage"
	"github.com/nfnt/resize"
)

type ImageService struct {
	logger  *logger.Logger
	storage storage.Storage
	cache   lrucache.Cache
}

func NewImageService(
	logger *logger.Logger,
	storage storage.Storage,
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
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", imgParams.URL, nil)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req = req.WithContext(ctx)

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	fmt.Println("Downloaded from URL:", imgParams.URL)
	sourceImg, _, err := image.Decode(response.Body)
	if err != nil {
		return nil, err
	}

	s.cache.Set(lrucache.Key(imageID), sourceImg)
	newImg := resize.Resize(imgParams.Width, imgParams.Height, sourceImg, resize.Lanczos3)

	err = s.storage.Set(newImg, imageID)
	if err != nil {
		return nil, err
	}

	return newImg, nil
}
