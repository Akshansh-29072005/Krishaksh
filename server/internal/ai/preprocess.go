package ai

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"time"
)

func FetchAndPreprocess(ctx context.Context, imageURL string) ([]byte, string, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, imageURL, nil)
	client := &http.Client{Timeout: 12 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		return nil, "", fmt.Errorf("image fetch status %d", res.StatusCode)
	}
	raw, err := io.ReadAll(io.LimitReader(res.Body, 12*1024*1024))
	if err != nil {
		return nil, "", err
	}
	img, _, err := image.Decode(bytes.NewReader(raw))
	if err != nil {
		return raw, "image/jpeg", nil
	}
	maxW, maxH := 1024, 1024
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if w > maxW || h > maxH {
		img = resizeNearest(img, maxW, maxH)
	}
	buf := bytes.NewBuffer(nil)
	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 75}); err != nil {
		return raw, "image/jpeg", nil
	}
	return buf.Bytes(), "image/jpeg", nil
}

func resizeNearest(src image.Image, maxW, maxH int) image.Image {
	b := src.Bounds()
	sw, sh := b.Dx(), b.Dy()
	ratioW := float64(maxW) / float64(sw)
	ratioH := float64(maxH) / float64(sh)
	ratio := ratioW
	if ratioH < ratio {
		ratio = ratioH
	}
	if ratio > 1 {
		ratio = 1
	}
	dw, dh := int(float64(sw)*ratio), int(float64(sh)*ratio)
	dst := image.NewRGBA(image.Rect(0, 0, dw, dh))
	for y := 0; y < dh; y++ {
		sy := int(float64(y) / ratio)
		for x := 0; x < dw; x++ {
			sx := int(float64(x) / ratio)
			dst.Set(x, y, src.At(b.Min.X+sx, b.Min.Y+sy))
		}
	}
	return dst
}
