package httpserver

import (
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"floway-backend/internal/storage"
)

const maxUploadSize = 10 << 20 // 10 MB

var allowedImageExtensions = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

type uploadHandler struct {
	storage *storage.Client
	admin   func(http.Handler) http.Handler
}

func newUploadHandler(storage *storage.Client, admin func(http.Handler) http.Handler) *uploadHandler {
	return &uploadHandler{storage: storage, admin: admin}
}

func (h *uploadHandler) adminRoutes(r chi.Router) {
	r.With(h.admin).Post("/", h.upload)
}

func (h *uploadHandler) upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	defer func() { _ = file.Close() }()

	contentType := header.Header.Get("Content-Type")
	ext, ok := allowedImageExtensions[contentType]
	if !ok {
		writeError(w, http.StatusBadRequest, errors.New("unsupported image type, allowed: jpeg, png, webp, gif"))
		return
	}

	key := uuid.NewString() + ext
	if err := h.storage.Upload(r.Context(), key, file, header.Size, contentType); err != nil {
		writeInternalError(w, r, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"url": "/uploads/" + key})
}

func (h *uploadHandler) serve(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	obj, err := h.storage.Download(r.Context(), key)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}
	defer func() { _ = obj.Close() }()

	info, err := obj.Stat()
	if err != nil {
		if storage.IsNotFound(err) {
			writeError(w, http.StatusNotFound, errors.New("not found"))
			return
		}
		writeInternalError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", info.ContentType)
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	_, _ = io.Copy(w, obj)
}
