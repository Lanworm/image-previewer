package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/Lanworm/OTUS_GO/final_project/internal/validation"
)

var (
	ErrInvalidNumberOfArguments           = errors.New("wrong number of arguments")
	ErrInvalidFormatOfArguments           = errors.New("wrong format of arguments")
	ErrInvalidArgumentTypeOfWidthOrHeight = errors.New("invalid argument type of width or height")
	ErrInvalidURL                         = errors.New("invalid URL")
)

func PrepareImgParams(u *url.URL) (imgParams *ImgParams, err error) {
	parts := strings.Split(u.String(), "/")

	if len(parts) < 4 {
		return nil, ErrInvalidNumberOfArguments
	}

	width := parts[2]
	height := parts[3]

	imageURLParts := parts[4:]
	imageURL := strings.Join(imageURLParts, "/")

	// Удаляем "https/" из URL, если присутствует
	imageURL = strings.Replace(imageURL, "https:/", "", -1)

	// Добавляем 'https://' в URL, если отсутствует
	if !strings.HasPrefix(imageURL, "http://") && !strings.HasPrefix(imageURL, "https://") {
		imageURL = "https://" + imageURL
	}

	// Удаляем лишний символ '/' в конце URL изображения
	if strings.HasSuffix(imageURL, "/") {
		imageURL = imageURL[:len(imageURL)-1]
	}

	// Проверяем валидность URl
	_, err = url.ParseRequestURI(imageURL)
	if err != nil {
		return nil, ErrInvalidURL
	}

	params, err := NewImgParams(width, height, imageURL)
	if err != nil {
		return nil, ErrInvalidFormatOfArguments
	}
	return params, nil
}

func NewImgParams(width string, height string, url string) (*ImgParams, error) {
	w, errw := strconv.Atoi(width)
	h, errh := strconv.Atoi(height)

	if errw != nil || errh != nil {
		return nil, ErrInvalidArgumentTypeOfWidthOrHeight
	}

	p := &ImgParams{
		Width:  uint(w),
		Height: uint(h),
		URL:    url,
	}

	err := validation.Validate(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func getURLHash(url string) string {
	hash := md5.Sum([]byte(url))
	return hex.EncodeToString(hash[:])
}
