package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Lanworm/image-previewer/internal/validation"
	"github.com/gorilla/mux"
)

var (
	ErrInvalidFormatOfArguments           = errors.New("wrong format of arguments")
	ErrInvalidArgumentTypeOfWidthOrHeight = errors.New("invalid argument type of width or height")
	ErrInvalidURL                         = errors.New("invalid URL")
)

func PrepareImgParams(r *http.Request) (imgParams *ImgParams, err error) {
	vars := mux.Vars(r)
	width := vars["width"]
	height := vars["height"]
	imageURL := vars["url"]

	// Удаляем "https/" из URL, если присутствует
	imageURL = strings.ReplaceAll(imageURL, "http:/", "")
	imageURL = strings.ReplaceAll(imageURL, "https:/", "")

	// Добавляем 'https://' в URL, если отсутствует
	if !strings.HasPrefix(imageURL, "http://") && !strings.HasPrefix(imageURL, "https://") {
		imageURL = "http://" + imageURL
	}

	// Удаляем лишний символ '/' в конце URL изображения
	imageURL = strings.TrimSuffix(imageURL, "/")

	// Проверяем валидность URl
	_, err = url.ParseRequestURI(imageURL)
	if err != nil {
		return nil, ErrInvalidURL
	}
	// Создаем новую структуру с параметрами
	params, err := NewImgParams(width, height, imageURL)
	if err != nil {
		return nil, ErrInvalidFormatOfArguments
	}
	return params, nil
}

func NewImgParams(width string, height string, url string) (*ImgParams, error) {
	w, errw := strconv.Atoi(width)
	h, errh := strconv.Atoi(height)
	//
	if errw != nil || errh != nil {
		return nil, ErrInvalidArgumentTypeOfWidthOrHeight
	}

	p := &ImgParams{
		Width:  w,
		Height: h,
		URL:    url,
	}

	if err := validation.Validate(p); err != nil {
		return nil, err
	}
	return p, nil
}

func getURLHash(url string) string {
	hasher := sha256.New()
	hasher.Write([]byte(url))
	hashInBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashInBytes)
}
