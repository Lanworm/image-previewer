package service

import (
	"context"
	"errors"
	"fmt"
	"image"
	"net/http"
	"strings"
	"time"

	lrucache "github.com/Lanworm/image-previewer/internal/cache"
	"github.com/Lanworm/image-previewer/internal/logger"
	"github.com/Lanworm/image-previewer/internal/storage"
	"github.com/nfnt/resize"
)

type ImageService struct {
	logger       *logger.Logger
	storage      storage.Storage
	cache        lrucache.Cache
	maxImageSize int
}

func NewImageService(
	logger *logger.Logger,
	storage storage.Storage,
	cache lrucache.Cache,
	maxImageSize int,
) *ImageService {
	return &ImageService{
		logger:       logger,
		storage:      storage,
		cache:        cache,
		maxImageSize: maxImageSize,
	}
}

type ImgParams struct {
	Width  int `validate:"required,gt=0,lte=9999"`
	Height int `validate:"required,gt=0,lte=9999"`
	URL    string
}

var (
	ErrImageNotFound  = errors.New("image not found on remote server")
	ErrTargetNotImage = errors.New("requested URL does not point to an image")
	ErrImageSize      = errors.New("image size exceeds the limit")
	ErrOutOfBounds    = errors.New("image is out of bounds")
)

func (s *ImageService) ResizeImg(imgParams *ImgParams) (img image.Image, err error) {
	imageID := getURLHash(imgParams.URL)

	// Проверяем наличие изображения в кэше
	sourceImg, ok := s.cache.Get(lrucache.Key(imageID))

	// Если изображение найдено в кэше, отдаем его
	if ok {
		fmt.Println("received from cache: ", imageID)
		cachedImg := resize.Resize(uint(imgParams.Width), uint(imgParams.Height), sourceImg, resize.Lanczos3)
		return cachedImg, nil
	}

	// Если изображение не найдено в кэше, загружаем его
	client := &http.Client{}
	req, err := http.NewRequest("GET", imgParams.URL, nil)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req = req.WithContext(ctx)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		return nil, ErrImageNotFound
	}

	// Проверяем тип контента
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "image") {
		return nil, ErrTargetNotImage
	}

	// Проверяем размер изображения
	if resp.ContentLength > int64(s.maxImageSize*1024) {
		return nil, ErrImageSize
	}

	fmt.Println("Downloaded from URL:", imgParams.URL)
	// Читаем изображение
	sourceImg, _, err = image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	// Проверка размеров изображения и параметров ширины и высоты
	bounds := sourceImg.Bounds()
	imgWidth := bounds.Max.X
	imgHeight := bounds.Max.Y

	if imgWidth <= imgParams.Width && imgHeight <= imgParams.Width {
		return nil, ErrOutOfBounds
	}

	s.cache.Set(lrucache.Key(imageID), sourceImg)
	newImg := resize.Resize(uint(imgParams.Width), uint(imgParams.Height), sourceImg, resize.Lanczos3)

	err = s.storage.Set(newImg, imageID)
	if err != nil {
		return nil, err
	}

	return newImg, nil
}
