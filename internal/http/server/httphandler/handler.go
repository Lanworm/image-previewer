package httphandler

import (
	"bytes"
	"encoding/json"
	"image/jpeg"
	"net/http"
	"strconv"

	"github.com/Lanworm/OTUS_GO/final_project/internal/http/server/dto"
	"github.com/Lanworm/OTUS_GO/final_project/internal/logger"
	"github.com/Lanworm/OTUS_GO/final_project/internal/service"
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
	imgParams, err := service.PrepareImgParams(r.URL)
	if err != nil {
		writeError(http.StatusInternalServerError, w, err.Error())
		h.logger.Error(err.Error())
		return
	}
	img, err := h.service.ResizeImg(imgParams)
	if err != nil {
		writeError(http.StatusInternalServerError, w, err.Error())
		h.logger.Error(err.Error())
		return
	}

	buf := new(bytes.Buffer)
	encodeErr := jpeg.Encode(buf, img, nil)
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))
	if encodeErr != nil {
		writeError(http.StatusInternalServerError, w, encodeErr.Error())
		h.logger.Error(encodeErr.Error())
		return
	}
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
