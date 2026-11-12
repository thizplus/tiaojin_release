package services

import (
	"context"
	"gofiber-template/domain/dto"
)

type LLMService interface {
	ExtractPassportFields(ctx context.Context, ocrText string) (*dto.PassportExtractionResult, error)
	ParseStructuredData(ctx context.Context, text string, prompt string) (string, error)
}