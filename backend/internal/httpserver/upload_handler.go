package httpserver

import (
	"bytes"
	"errors"
	"image"
	_ "image/jpeg" // registers the jpeg decoder for image.DecodeConfig
	_ "image/png"  // registers the png decoder for image.DecodeConfig
	"io"
	"net/http"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"floway-backend/internal/storage"
)

const (
	maxUploadSize = 10 << 20 // 10 MB

	// Admin-uploaded photos are typically straight off a phone camera —
	// several thousand pixels wide, multiple megabytes — far beyond
	// anything a single full-size image view on the site needs. Downscaling
	// once here, at upload time, is what actually fixes "photos take a
	// moment to load": no amount of client-side caching or preloading helps
	// until the file itself is a reasonable size.
	maxImageDimension = 2000 // px, longest side
	jpegQuality       = 85

	// thumbnailMaxWidth caps the `?w=` query param on GET /uploads/{key}
	// (see serve/resizeToWidth below) — that param exists for in-page
	// thumbnails (e.g. the homepage gallery strip, cards up to ~66vh tall),
	// not as a general resize API. Anything wanting more detail should
	// request the unscaled URL, already capped at maxImageDimension from
	// upload time.
	thumbnailMaxWidth = 900
)

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

	file, _, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	defer func() { _ = file.Close() }()

	// The multipart Content-Type header is client-supplied and trivially
	// spoofed (a script can name any file "photo.png" with any header it
	// likes) — sniff the actual bytes instead of trusting it, then seek
	// back to the start so storage.Upload reads the whole file.
	sniffBuf := make([]byte, 512)
	n, err := file.Read(sniffBuf)
	if err != nil && err != io.EOF {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	contentType := http.DetectContentType(sniffBuf[:n])
	ext, ok := allowedImageExtensions[contentType]
	if !ok {
		writeError(w, http.StatusBadRequest, errors.New("unsupported image type, allowed: jpeg, png, webp, gif"))
		return
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		writeInternalError(w, r, err)
		return
	}

	body, size, err := downscaleIfNeeded(file, contentType)
	if err != nil {
		writeInternalError(w, r, err)
		return
	}

	key := uuid.NewString() + ext
	if err := h.storage.Upload(r.Context(), key, body, size, contentType); err != nil {
		writeInternalError(w, r, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"url": "/uploads/" + key})
}

// downscaleIfNeeded shrinks jpeg/png uploads whose longest side exceeds
// maxImageDimension. gif and webp pass through untouched — resizing a gif
// here would need frame-by-frame handling to not break animation, and this
// avoids taking on webp encoding at all. Any failure to decode/resize falls
// back to uploading the original bytes — this is a size optimization, not
// something that should ever fail the upload itself.
func downscaleIfNeeded(r io.Reader, contentType string) (io.Reader, int64, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, 0, err
	}
	if contentType == "image/jpeg" || contentType == "image/png" {
		if resized, ok := resizeLongestSide(data, contentType, maxImageDimension); ok {
			return bytes.NewReader(resized), int64(len(resized)), nil
		}
	}
	return bytes.NewReader(data), int64(len(data)), nil
}

// resizeLongestSide decodes data and, if its longest side exceeds maxDim,
// resizes down to maxDim (preserving aspect ratio) and re-encodes at
// jpegQuality (jpeg) or standard compression (png). Returns ok=false —
// meaning "use the original bytes" — when the image is already within
// maxDim (checked via image.DecodeConfig, which reads only the header, so
// this costs nothing for the common case) or when decoding/encoding fails.
func resizeLongestSide(data []byte, contentType string, maxDim int) ([]byte, bool) {
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || (cfg.Width <= maxDim && cfg.Height <= maxDim) {
		return nil, false
	}

	img, err := imaging.Decode(bytes.NewReader(data), imaging.AutoOrientation(true))
	if err != nil {
		return nil, false
	}

	var resized image.Image
	if cfg.Width >= cfg.Height {
		resized = imaging.Resize(img, maxDim, 0, imaging.Lanczos)
	} else {
		resized = imaging.Resize(img, 0, maxDim, imaging.Lanczos)
	}
	return encodeImage(resized, contentType)
}

// resizeToWidth is resizeLongestSide's counterpart for the `?w=` thumbnail
// param on serve() below: a fixed target width (height following the
// aspect ratio) rather than a cap on whichever side is longest, since
// thumbnail slots are a fixed width regardless of a photo's orientation.
func resizeToWidth(data []byte, contentType string, width int) ([]byte, bool) {
	if width > thumbnailMaxWidth {
		width = thumbnailMaxWidth
	}
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || cfg.Width <= width {
		return nil, false
	}

	img, err := imaging.Decode(bytes.NewReader(data), imaging.AutoOrientation(true))
	if err != nil {
		return nil, false
	}
	return encodeImage(imaging.Resize(img, width, 0, imaging.Lanczos), contentType)
}

func encodeImage(img image.Image, contentType string) ([]byte, bool) {
	var buf bytes.Buffer
	var err error
	if contentType == "image/jpeg" {
		err = imaging.Encode(&buf, img, imaging.JPEG, imaging.JPEGQuality(jpegQuality))
	} else {
		err = imaging.Encode(&buf, img, imaging.PNG)
	}
	if err != nil {
		return nil, false
	}
	return buf.Bytes(), true
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
	// The Content-Type above is now trustworthy (sniffed from real bytes at
	// upload time, not client-declared) — nosniff tells the browser to
	// honor it as-is instead of re-guessing from the body, closing the
	// usual "served image gets reinterpreted as HTML/script" MIME-confusion
	// class of bug.
	w.Header().Set("X-Content-Type-Options", "nosniff")

	// @nuxt/image's IPX sends a HEAD request for exactly these three headers
	// before ever fetching/transforming an image, to decide its own
	// Cache-Control — this route only had GET registered, so that HEAD
	// 405'd, IPX silently treated the failure as "no cache metadata", and
	// every transformed image it served came back with no Cache-Control at
	// all (confirmed live: every reload re-downloaded every image). No body
	// for HEAD, same as any well-behaved static-file-like endpoint.
	if r.Method == http.MethodHead {
		return
	}

	// `?w=` lets a caller ask for a small thumbnail-sized rendition instead
	// of the full (already upload-time-capped) image — the homepage gallery
	// strip uses it so its several-photos-in-a-row thumbnails don't each
	// pull down a ~2000px full-size file just to display at ~200px. Cached
	// forever per the header above, same as the unscaled URL, so this only
	// costs a resize once per browser.
	if width, ok := parsePositiveInt(r.URL.Query().Get("w")); ok &&
		(info.ContentType == "image/jpeg" || info.ContentType == "image/png") {
		data, err := io.ReadAll(obj)
		if err != nil {
			writeInternalError(w, r, err)
			return
		}
		if resized, ok := resizeToWidth(data, info.ContentType, width); ok {
			_, _ = w.Write(resized)
			return
		}
		_, _ = w.Write(data)
		return
	}
	_, _ = io.Copy(w, obj)
}

func parsePositiveInt(s string) (int, bool) {
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return 0, false
	}
	return n, true
}
