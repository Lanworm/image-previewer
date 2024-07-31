package httphandler

import (
	"bytes"
	"encoding/json"
	"image/jpeg"
	"net/http"
	"strconv"

	"github.com/Lanworm/image-previewer/internal/http/server/dto"
	"github.com/Lanworm/image-previewer/internal/logger"
	"github.com/Lanworm/image-previewer/internal/service"
)

type Handler struct {
	logger  *logger.Logger
	service *service.ImageService
}

func NewHandler(
	logger *logger.Logger,
	service *service.ImageService,
) *Handler {
	return &Handler{
		logger:  logger,
		service: service,
	}
}

func (h *Handler) ResizeHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	// Подготовка параметров изображения из запроса
	imgParams, err := service.PrepareImgParams(r)
	if err != nil {
		// Обработка ошибки и отправка ответа с кодом 500
		writeError(http.StatusInternalServerError, w, err.Error())
		h.logger.Error(err.Error())
		return
	}

	// Изменение размера изображения
	img, err := h.service.ResizeImg(imgParams, r)
	if err != nil {
		// Обработка ошибки и отправка ответа с кодом 500
		writeError(http.StatusInternalServerError, w, err.Error())
		h.logger.Error(err.Error())
		return
	}

	// Кодирование изображения в формат JPEG
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img, nil)
	if err != nil {
		// Обработка ошибки и отправка ответа с кодом 500
		writeError(http.StatusInternalServerError, w, err.Error())
		h.logger.Error(err.Error())
		return
	}

	// Установка заголовков и отправка изображения в ответе
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	w.Write(buf.Bytes())
}

// Функция для отправки ошибки в ответе.
func writeError(
	statusCode int,
	w http.ResponseWriter,
	msg string,
) {
	// Создание JSON с сообщением об ошибке
	js, err := json.Marshal(dto.Result{Message: msg})
	if err != nil {
		// Если возникла ошибка при сериализации, отправляем код 500
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		// Отправка ответа с указанным статусом и сообщением
		w.WriteHeader(statusCode)
		w.Write(js)
	}
}
