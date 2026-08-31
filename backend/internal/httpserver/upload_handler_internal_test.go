package httpserver

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func encodeJPEG(t *testing.T, width, height int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x % 256), G: uint8(y % 256), B: 128, A: 255})
		}
	}
	var buf bytes.Buffer
	require.NoError(t, jpeg.Encode(&buf, img, &jpeg.Options{Quality: 95}))
	return buf.Bytes()
}

func encodePNG(t *testing.T, width, height int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	var buf bytes.Buffer
	require.NoError(t, png.Encode(&buf, img))
	return buf.Bytes()
}

func TestDownscaleIfNeeded_ShrinksOversizedJPEG(t *testing.T) {
	original := encodeJPEG(t, 3130, 2075)

	out, size, err := downscaleIfNeeded(bytes.NewReader(original), "image/jpeg")

	require.NoError(t, err)
	data, err := io.ReadAll(out)
	require.NoError(t, err)
	assert.Equal(t, int64(len(data)), size)

	cfg, err := jpeg.DecodeConfig(bytes.NewReader(data))
	require.NoError(t, err)
	assert.LessOrEqual(t, cfg.Width, maxImageDimension)
	assert.LessOrEqual(t, cfg.Height, maxImageDimension)
	assert.Less(t, len(data), len(original), "a downscaled photo should end up smaller than the original")
}

func TestDownscaleIfNeeded_ShrinksOversizedPNG(t *testing.T) {
	original := encodePNG(t, 2400, 1200)

	out, _, err := downscaleIfNeeded(bytes.NewReader(original), "image/png")

	require.NoError(t, err)
	data, err := io.ReadAll(out)
	require.NoError(t, err)

	cfg, err := png.DecodeConfig(bytes.NewReader(data))
	require.NoError(t, err)
	assert.Equal(t, maxImageDimension, cfg.Width, "width was the longer side, so it should land exactly on the cap")
	assert.Less(t, cfg.Height, 1200)
}

func TestDownscaleIfNeeded_LeavesSmallImageByteForByteUnchanged(t *testing.T) {
	original := encodeJPEG(t, 400, 300)

	out, size, err := downscaleIfNeeded(bytes.NewReader(original), "image/jpeg")

	require.NoError(t, err)
	data, err := io.ReadAll(out)
	require.NoError(t, err)
	assert.Equal(t, original, data, "an already-small image shouldn't be re-encoded at all")
	assert.Equal(t, int64(len(original)), size)
}

func TestDownscaleIfNeeded_PassesThroughGifAndWebpUnchanged(t *testing.T) {
	fakeGifBytes := []byte("not actually decodable but that's fine, gif/webp never get decoded")

	out, size, err := downscaleIfNeeded(bytes.NewReader(fakeGifBytes), "image/gif")

	require.NoError(t, err)
	data, err := io.ReadAll(out)
	require.NoError(t, err)
	assert.Equal(t, fakeGifBytes, data)
	assert.Equal(t, int64(len(fakeGifBytes)), size)
}

func TestDownscaleIfNeeded_FallsBackToOriginalOnUndecodableBytes(t *testing.T) {
	garbage := []byte("this claims to be a jpeg but isn't one")

	out, size, err := downscaleIfNeeded(bytes.NewReader(garbage), "image/jpeg")

	require.NoError(t, err, "a failed optimization attempt should never fail the upload itself")
	data, err := io.ReadAll(out)
	require.NoError(t, err)
	assert.Equal(t, garbage, data)
	assert.Equal(t, int64(len(garbage)), size)
}

func TestResizeToWidth_ShrinksToRequestedWidth(t *testing.T) {
	original := encodeJPEG(t, 2000, 1000)

	data, ok := resizeToWidth(original, "image/jpeg", 400)

	require.True(t, ok)
	cfg, err := jpeg.DecodeConfig(bytes.NewReader(data))
	require.NoError(t, err)
	assert.Equal(t, 400, cfg.Width)
	assert.Equal(t, 200, cfg.Height, "height should follow the same aspect ratio as the original")
	assert.Less(t, len(data), len(original))
}

func TestResizeToWidth_ClampsToThumbnailMaxWidth(t *testing.T) {
	original := encodeJPEG(t, 2000, 1000)

	data, ok := resizeToWidth(original, "image/jpeg", thumbnailMaxWidth*10)

	require.True(t, ok, "an oversized request should be clamped, not rejected")
	cfg, err := jpeg.DecodeConfig(bytes.NewReader(data))
	require.NoError(t, err)
	assert.Equal(t, thumbnailMaxWidth, cfg.Width)
}

func TestResizeToWidth_LeavesAlreadySmallerImageAlone(t *testing.T) {
	original := encodeJPEG(t, 300, 200)

	_, ok := resizeToWidth(original, "image/jpeg", 400)

	assert.False(t, ok, "an image already narrower than the requested width has nothing to shrink")
}

func TestParsePositiveInt(t *testing.T) {
	cases := map[string]struct {
		want int
		ok   bool
	}{
		"400":  {400, true},
		"0":    {0, false},
		"-40":  {0, false},
		"":     {0, false},
		"abc":  {0, false},
		"40.5": {0, false},
	}
	for input, want := range cases {
		n, ok := parsePositiveInt(input)
		assert.Equal(t, want.ok, ok, "input %q", input)
		assert.Equal(t, want.want, n, "input %q", input)
	}
}
