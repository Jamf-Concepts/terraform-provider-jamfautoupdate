// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package titles

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"testing"
)

// createTestPNG generates a minimal valid PNG image of the given dimensions and returns it as a base64 string.
func createTestPNG(t *testing.T, width, height int) string {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("failed to encode test PNG: %v", err)
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func TestProcessUninstallIcon_ValidPNG(t *testing.T) {
	input := createTestPNG(t, 64, 64)
	result, err := processUninstallIcon(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected non-nil result")
	}

	decoded, err := base64.StdEncoding.DecodeString(*result)
	if err != nil {
		t.Fatalf("result is not valid base64: %v", err)
	}

	img, err := png.Decode(bytes.NewReader(decoded))
	if err != nil {
		t.Fatalf("result is not valid PNG: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() != BaseImageSize || bounds.Dy() != BaseImageSize {
		t.Errorf("expected %dx%d image, got %dx%d", BaseImageSize, BaseImageSize, bounds.Dx(), bounds.Dy())
	}
}

func TestProcessUninstallIcon_InvalidBase64(t *testing.T) {
	_, err := processUninstallIcon("not-valid-base64!!!")
	if err == nil {
		t.Fatal("expected error for invalid base64")
	}
}

func TestProcessUninstallIcon_InvalidImageData(t *testing.T) {
	input := base64.StdEncoding.EncodeToString([]byte("not an image"))
	_, err := processUninstallIcon(input)
	if err == nil {
		t.Fatal("expected error for invalid image data")
	}
}

func TestOverlayImageBase64_ValidPNG(t *testing.T) {
	decoded, err := base64.StdEncoding.DecodeString(OverlayImageBase64)
	if err != nil {
		t.Fatalf("OverlayImageBase64 is not valid base64: %v", err)
	}

	img, err := png.Decode(bytes.NewReader(decoded))
	if err != nil {
		t.Fatalf("OverlayImageBase64 is not valid PNG: %v", err)
	}

	bounds := img.Bounds()
	if bounds.Dx() == 0 || bounds.Dy() == 0 {
		t.Error("overlay image has zero dimensions")
	}
}
