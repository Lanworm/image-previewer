package service

import (
	"errors"
	"fmt"
	"image"
	"net/http"
	"strconv"
	"strings"
	"time"

	lrucache "github.com/Lanworm/image-previewer/internal/cache"
	"github.com/Lanworm/image-previewer/internal/http/client"
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
)

func (s *ImageService) ResizeImg(imgParams *ImgParams, r *http.Request) (img image.Image, err error) {
	imageID := getURLHash("resize" + strconv.Itoa(imgParams.Width) + strconv.Itoa(imgParams.Height) + imgParams.URL)

	// Проверяем наличие изображения в кэше
	cachedImg, ok := s.cache.Get(lrucache.Key(imageID))

	// Если изображение найдено в кэше, отдаем его
	if ok {
		fmt.Println("received from cache: ", imageID)
		return cachedImg, nil
	}

	// Если изображение не найдено в кэше, загружаем его
	sourceImg, err := s.getImage(imgParams.URL, r)
	if err != nil {
		return nil, err
	}
	// Изменяем размер
	resizedImg := resize.Resize(uint(imgParams.Width), uint(imgParams.Height), sourceImg, resize.Lanczos3)

	// Кладем измененное изображение в кеш
	s.cache.Set(lrucache.Key(imageID), resizedImg)

	// Записываем измененное изображение в хранилище
	err = s.storage.Set(resizedImg, imageID)
	if err != nil {
		return nil, err
	}

	return resizedImg, nil
}

func (s *ImageService) getImage(imgURL string, r *http.Request) (image.Image, error) {
	HTTPClient := client.NewHTTPClient(10 * time.Second)

	resp, err := HTTPClient.DoRequest("GET", imgURL, nil, r.Header)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// Проверяем статус ответа
	if resp.StatusCode == 404 {
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

	fmt.Println("Downloaded from URL:", imgURL)
	// Читаем изображение
	sourceImg, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	return sourceImg, nil
}
