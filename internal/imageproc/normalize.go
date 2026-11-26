package imageproc

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"io"

	"github.com/disintegration/imaging"
	"github.com/gen2brain/vpx/webp"
)

type NormalizedImage struct {
	Buffer *bytes.Buffer
	Width  int
	Height int
}

func Normalize(
	ctx context.Context,
	source io.Reader,
	cfg Config,
) (*NormalizedImage, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	sourceImage, err := imaging.Decode(source, imaging.AutoOrientation(true))
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	if bounds := sourceImage.Bounds(); bounds.Dx() > cfg.MaxWidth || bounds.Dy() > cfg.MaxHeight {
		sourceImage = imaging.Fit(sourceImage, cfg.MaxWidth, cfg.MaxHeight, imaging.Lanczos)
	}

	encodedImage, err := encodeWithinSizeLimit(ctx, sourceImage, cfg)
	if err != nil {
		return nil, fmt.Errorf("webp encode: %w", err)
	}

	bounds := sourceImage.Bounds()
	return &NormalizedImage{Buffer: encodedImage, Width: bounds.Dx(), Height: bounds.Dy()}, nil
}

func encodeWithinSizeLimit(ctx context.Context, sourceImage image.Image, cfg Config) (*bytes.Buffer, error) {
	for quality := cfg.InitialQuality; quality >= cfg.MinQuality; quality -= cfg.QualityStep {
		var encodedImage bytes.Buffer
		if err := webp.Encode(&encodedImage, sourceImage, webp.EncodeOptions{Quality: quality, Method: cfg.EncodeMethod}); err != nil {
			return nil, err
		}
		if encodedImage.Len() <= cfg.MaxSizeBytes {
			return &encodedImage, nil
		}
	}
	return nil, fmt.Errorf("failed to fit into %d bytes", cfg.MaxSizeBytes)
}
