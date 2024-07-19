package httphandler

import (
	"bytes"
	"image/jpeg"
	"net/http"
	"strconv"

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
		http.Error(w, err.Error(), http.StatusBadGateway)
		h.logger.Error(err.Error())
		return
	}

	img, err := h.service.ResizeImg(imgParams, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		h.logger.Error(err.Error())
		return
	}

	buf := new(bytes.Buffer)
	encodeErr := jpeg.Encode(buf, img, nil)

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	if encodeErr != nil {
		http.Error(w, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
		h.logger.Error(encodeErr.Error())
		return
	}
	w.Write(buf.Bytes())
}
