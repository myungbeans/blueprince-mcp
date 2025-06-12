package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"golang.org/x/image/draw"
)

// CompressedImageResult holds the compressed image data and metadata
type CompressedImageResult struct {
	Data             string  `json:"data"`            // base64 encoded compressed image
	Format           string  `json:"format"`          // jpeg, png, etc.
	OriginalSize     int64   `json:"original_size"`   // original file size in bytes
	CompressedSize   int     `json:"compressed_size"` // compressed size in bytes
	Width            int     `json:"width"`
	Height           int     `json:"height"`
	CompressionRatio float64 `json:"compression_ratio"`
}

// CompressImage compresses an image file with multiple strategies
func CompressImage(filepath string, maxWidth, maxHeight int, quality int) (*CompressedImageResult, error) {
	// Read original file
	originalData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	originalSize := int64(len(originalData))

	// Decode the image
	img, format, err := image.Decode(bytes.NewReader(originalData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Get original dimensions
	bounds := img.Bounds()
	originalWidth := bounds.Dx()
	originalHeight := bounds.Dy()

	// Resize if needed
	resizedImg := img
	if originalWidth > maxWidth || originalHeight > maxHeight {
		resizedImg = resizeImage(img, maxWidth, maxHeight)
	}

	// Compress based on format
	var compressedData []byte
	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		compressedData, err = compressJPEG(resizedImg, quality)
	case "png":
		compressedData, err = compressPNG(resizedImg)
	default:
		// Convert to JPEG for unsupported formats
		compressedData, err = compressJPEG(resizedImg, quality)
		format = "jpeg"
	}

	if err != nil {
		return nil, fmt.Errorf("failed to compress image: %w", err)
	}

	// Calculate compression ratio
	compressionRatio := float64(originalSize) / float64(len(compressedData))

	// Encode to base64
	base64Data := base64.StdEncoding.EncodeToString(compressedData)

	return &CompressedImageResult{
		Data:             base64Data,
		Format:           format,
		OriginalSize:     originalSize,
		CompressedSize:   len(compressedData),
		Width:            resizedImg.Bounds().Dx(),
		Height:           resizedImg.Bounds().Dy(),
		CompressionRatio: compressionRatio,
	}, nil
}

// resizeImage resizes an image while maintaining aspect ratio
func resizeImage(src image.Image, maxWidth, maxHeight int) image.Image {
	bounds := src.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Calculate new dimensions maintaining aspect ratio
	ratio := float64(width) / float64(height)
	newWidth := maxWidth
	newHeight := int(float64(newWidth) / ratio)

	if newHeight > maxHeight {
		newHeight = maxHeight
		newWidth = int(float64(newHeight) * ratio)
	}

	// Create new image
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Use bilinear scaling for better quality
	draw.BiLinear.Scale(dst, dst.Bounds(), src, src.Bounds(), draw.Over, nil)

	return dst
}

// compressJPEG compresses an image as JPEG with specified quality
func compressJPEG(img image.Image, quality int) ([]byte, error) {
	var buf bytes.Buffer

	options := &jpeg.Options{
		Quality: quality, // 1-100, lower = more compression
	}

	err := jpeg.Encode(&buf, img, options)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// compressPNG compresses an image as PNG (lossless)
func compressPNG(img image.Image) ([]byte, error) {
	var buf bytes.Buffer

	encoder := &png.Encoder{
		CompressionLevel: png.BestCompression, // Maximum compression
	}

	err := encoder.Encode(&buf, img)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
