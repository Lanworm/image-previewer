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
	imgParams, err := service.PrepareImgParams(r)
	if err != nil {
		writeError(http.StatusInternalServerError, w, err.Error())
		h.logger.Error(err.Error())
		return
	}

	img, err := h.service.ResizeImg(imgParams, r)
	if err != nil {
		writeError(http.StatusInternalServerError, w, err.Error())
		h.logger.Error(err.Error())
		return
	}

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, img, nil)
	if err != nil {
		writeError(http.StatusInternalServerError, w, err.Error())
		h.logger.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	w.Write(buf.Bytes())
}

func writeError(
	statusCode int,
	w http.ResponseWriter,
	msg string,
) {
	js, err := json.Marshal(dto.Result{Message: msg})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(statusCode)
		w.Write(js)
	}
}
