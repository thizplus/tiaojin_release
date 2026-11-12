package services

import (
	"context"
	"gofiber-template/domain/dto"
)

// VisionService handles Google Vision API text extraction
type VisionService interface {
	// ExtractTextFromImage performs OCR on an image and returns extracted text
	ExtractTextFromImage(ctx context.Context, imageData []byte) (*dto.VisionOCRResponse, error)

	// ExtractPassportData parses OCR text and extracts passport-specific information
	ExtractPassportData(ctx context.Context, ocrText string) (*dto.PassportExtractionResult, error)

	// ProcessPassportImage combines OCR and passport data extraction
	ProcessPassportImage(ctx context.Context, imageData []byte) (*dto.PassportExtractionResult, error)
}