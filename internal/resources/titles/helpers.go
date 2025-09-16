// Copyright 2025 Jamf Software LLC.

package titles

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"

	"github.com/nfnt/resize"
)

const (
	OverlayImageBase64 = "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"
	BaseImageSize      = 512
	OverlaySize        = 128
	TargetImageSize    = 512
)

type Processor struct{}

// NewProcessor creates a new instance of the image processor.
func NewProcessor() *Processor {
	return &Processor{}
}

// ProcessUninstallIcon processes a base64 encoded image by resizing it and adding an overlay.
// It takes a base64 encoded string of the original image and returns a base64 encoded string
// of the processed image. The processing includes resizing to the standard size and adding
// an uninstall overlay to the bottom right corner.
// Returns an error if any step of the image processing fails.
func (p *Processor) ProcessUninstallIcon(baseImageB64 string) (*string, error) {
	baseImageBytes, err := base64.StdEncoding.DecodeString(baseImageB64)
	if err != nil {
		return nil, fmt.Errorf("error decoding base image: %w", err)
	}

	baseImg, _, err := image.Decode(bytes.NewReader(baseImageBytes))
	if err != nil {
		return nil, fmt.Errorf("error decoding base image bytes: %w", err)
	}

	baseImg = resize.Resize(BaseImageSize, BaseImageSize, baseImg, resize.Lanczos3)

	bounds := baseImg.Bounds()
	rgba := image.NewRGBA(bounds)

	draw.Draw(rgba, bounds, baseImg, image.Point{}, draw.Src)

	overlayBytes, err := base64.StdEncoding.DecodeString(OverlayImageBase64)
	if err != nil {
		return nil, fmt.Errorf("error decoding overlay image: %w", err)
	}

	overlayImg, _, err := image.Decode(bytes.NewReader(overlayBytes))
	if err != nil {
		return nil, fmt.Errorf("error decoding overlay image bytes: %w", err)
	}

	overlayImg = resize.Resize(OverlaySize, OverlaySize, overlayImg, resize.Lanczos3)
	overlayBounds := overlayImg.Bounds()

	x := bounds.Max.X - overlayBounds.Dx()
	y := bounds.Max.Y - overlayBounds.Dy()

	offset := image.Point{X: x, Y: y}

	draw.Draw(rgba, image.Rectangle{
		Min: offset,
		Max: offset.Add(overlayBounds.Size()),
	}, overlayImg, image.Point{}, draw.Over)

	var buf bytes.Buffer
	if err := png.Encode(&buf, rgba); err != nil {
		return nil, fmt.Errorf("error encoding processed image: %w", err)
	}

	result := base64.StdEncoding.EncodeToString(buf.Bytes())
	return &result, nil
}
