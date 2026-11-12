package services

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	vision "cloud.google.com/go/vision/apiv1"
	visionpb "google.golang.org/genproto/googleapis/cloud/vision/v1"
	"google.golang.org/api/option"
	"gofiber-template/pkg/config"
	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
)

type VisionServiceImpl struct {
	client     *vision.ImageAnnotatorClient
	projectID  string
	llmService services.LLMService
}

func NewVisionService(cfg *config.Config) (services.VisionService, error) {
	ctx := context.Background()

	// Initialize Vision client with credentials
	var opts []option.ClientOption
	if cfg.GoogleCloud.CredentialsPath != "" {
		opts = append(opts, option.WithCredentialsFile(cfg.GoogleCloud.CredentialsPath))
	}

	client, err := vision.NewImageAnnotatorClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create vision client: %w", err)
	}

	// Initialize LLM service if enabled
	var llmService services.LLMService
	if cfg.LLM.Enabled {
		llmService, err = NewLLMService(cfg)
		if err != nil {
			log.Printf("Warning: Failed to initialize LLM service: %v", err)
			llmService = nil
		} else {
			log.Println("✓ LLM service initialized for enhanced passport parsing")
		}
	}

	return &VisionServiceImpl{
		client:     client,
		projectID:  cfg.GoogleCloud.ProjectID,
		llmService: llmService,
	}, nil
}

func (v *VisionServiceImpl) ExtractTextFromImage(ctx context.Context, imageData []byte) (*dto.VisionOCRResponse, error) {
	log.Printf("[INFO] Starting OCR text extraction, image size: %d bytes", len(imageData))

	// Create image object
	image := &visionpb.Image{
		Content: imageData,
	}

	// Perform text detection
	annotations, err := v.client.DetectTexts(ctx, image, nil, 10)

	if err != nil {
		log.Printf("[ERROR] Failed to detect text: %v", err)
		return nil, fmt.Errorf("failed to detect text: %w", err)
	}

	if len(annotations) == 0 {
		log.Printf("[WARNING] No text detected in image")
		return &dto.VisionOCRResponse{
			Text:       "",
			Confidence: 0.0,
		}, nil
	}

	// Extract the full text (first annotation contains all text)
	fullText := annotations[0]

	response := &dto.VisionOCRResponse{
		Text:       fullText.Description,
		Confidence: fullText.Confidence,
	}

	// Extract individual words if available
	if len(annotations) > 1 {
		for i := 1; i < len(annotations); i++ {
			word := annotations[i]
			response.Words = append(response.Words, dto.VisionWord{
				Text:       word.Description,
				Confidence: word.Confidence,
			})
		}
	}

	log.Printf("[SUCCESS] OCR completed, detected text length: %d, confidence: %.2f",
		len(response.Text), response.Confidence)

	return response, nil
}

func (v *VisionServiceImpl) ExtractPassportData(ctx context.Context, ocrText string) (*dto.PassportExtractionResult, error) {
	log.Printf("[INFO] Starting passport data extraction from OCR text")

	// Try LLM-based extraction first if available
	if v.llmService != nil {
		log.Printf("[INFO] Using LLM-based passport extraction")
		result, err := v.llmService.ExtractPassportFields(ctx, ocrText)
		if err != nil {
			log.Printf("[WARNING] LLM extraction failed, falling back to regex: %v", err)
		} else {
			log.Printf("[SUCCESS] LLM extraction completed with confidence: %.2f", result.OverallConfidence)
			return result, nil
		}
	}

	// Fallback to regex-based extraction
	log.Printf("[INFO] Using regex-based passport extraction (fallback)")
	result := &dto.PassportExtractionResult{
		RawText:               ocrText,
		ProcessedSuccessfully: true,
		FieldConfidences:      make(map[string]float32),
	}

	// Clean and normalize text
	text := strings.ToUpper(strings.TrimSpace(ocrText))
	lines := strings.Split(text, "\n")

	// Extract passport number (Thai passport pattern: G1234567)
	passportRegex := regexp.MustCompile(`[A-Z]\d{7}`)
	if matches := passportRegex.FindString(text); matches != "" {
		result.PassportNumber = &matches
		result.FieldConfidences["passport_number"] = 0.9
		log.Printf("[SUCCESS] Extracted passport number: %s", matches)
	}

	// Extract nationality (look for THA, THAILAND, etc.)
	if strings.Contains(text, "THA") || strings.Contains(text, "THAILAND") {
		nationality := "THA"
		result.Nationality = &nationality
		result.CountryCode = &nationality
		result.FieldConfidences["nationality"] = 0.95
	}

	// Extract dates (DD MMM YYYY or DD/MM/YYYY format)
	dateRegex := regexp.MustCompile(`(\d{1,2})\s+(JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|SEP|OCT|NOV|DEC)\s+(\d{4})`)
	dateMatches := dateRegex.FindAllStringSubmatch(text, -1)

	for _, match := range dateMatches {
		if len(match) == 4 {
			day, _ := strconv.Atoi(match[1])
			month := v.parseMonth(match[2])
			year, _ := strconv.Atoi(match[3])

			date := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

			// Determine if this is birth date or expiry date based on context
			if year < 2010 { // Likely birth date
				result.DateOfBirth = &date
				result.FieldConfidences["date_of_birth"] = 0.85
			} else if year > time.Now().Year() { // Future date, likely expiry
				result.DateOfExpiry = &date
				result.FieldConfidences["date_of_expiry"] = 0.85
			}
		}
	}

	// Extract gender
	if strings.Contains(text, "M/MALE") || strings.Contains(text, "MALE") {
		gender := "M"
		result.Gender = &gender
		result.FieldConfidences["gender"] = 0.9
	} else if strings.Contains(text, "F/FEMALE") || strings.Contains(text, "FEMALE") {
		gender := "F"
		result.Gender = &gender
		result.FieldConfidences["gender"] = 0.9
	}

	// Extract names (Thai passport typically has Thai and English names)
	result = v.extractNames(result, lines)

	// Calculate overall confidence
	totalConfidence := float32(0)
	fieldCount := 0
	for _, confidence := range result.FieldConfidences {
		totalConfidence += confidence
		fieldCount++
	}

	if fieldCount > 0 {
		result.OverallConfidence = totalConfidence / float32(fieldCount)
	}

	log.Printf("[SUCCESS] Passport data extraction completed, fields extracted: %d, overall confidence: %.2f",
		fieldCount, result.OverallConfidence)

	return result, nil
}

func (v *VisionServiceImpl) ProcessPassportImage(ctx context.Context, imageData []byte) (*dto.PassportExtractionResult, error) {
	log.Printf("[INFO] Starting complete passport image processing")

	// Step 1: Extract text using OCR
	ocrResponse, err := v.ExtractTextFromImage(ctx, imageData)
	if err != nil {
		log.Printf("[ERROR] OCR failed: %v", err)
		return &dto.PassportExtractionResult{
			ProcessedSuccessfully: false,
			ErrorMessage:          stringPtr(fmt.Sprintf("OCR failed: %v", err)),
		}, err
	}

	// Step 2: Extract passport data from OCR text
	result, err := v.ExtractPassportData(ctx, ocrResponse.Text)
	if err != nil {
		log.Printf("[ERROR] Passport data extraction failed: %v", err)
		return &dto.PassportExtractionResult{
			RawText:               ocrResponse.Text,
			ProcessedSuccessfully: false,
			ErrorMessage:          stringPtr(fmt.Sprintf("Data extraction failed: %v", err)),
		}, err
	}

	// Add OCR confidence to the result
	result.OverallConfidence = (result.OverallConfidence + ocrResponse.Confidence) / 2

	log.Printf("[SUCCESS] Complete passport processing finished successfully")
	return result, nil
}

// Helper functions

func (v *VisionServiceImpl) parseMonth(monthStr string) time.Month {
	months := map[string]time.Month{
		"JAN": time.January, "FEB": time.February, "MAR": time.March,
		"APR": time.April, "MAY": time.May, "JUN": time.June,
		"JUL": time.July, "AUG": time.August, "SEP": time.September,
		"OCT": time.October, "NOV": time.November, "DEC": time.December,
	}
	if month, ok := months[monthStr]; ok {
		return month
	}
	return time.January
}

func (v *VisionServiceImpl) extractNames(result *dto.PassportExtractionResult, lines []string) *dto.PassportExtractionResult {
	// This is a simplified name extraction - could be enhanced with more sophisticated parsing
	for i, line := range lines {
		line = strings.TrimSpace(line)

		// Look for patterns that indicate names
		if strings.Contains(line, "NAME") || strings.Contains(line, "ชื่อ") {
			// Try to extract names from following lines
			if i+1 < len(lines) {
				nextLine := strings.TrimSpace(lines[i+1])
				if nextLine != "" && !strings.Contains(nextLine, "PASSPORT") {
					// Simple heuristic: if it contains Thai characters, it's Thai name
					if v.containsThai(nextLine) {
						parts := strings.Fields(nextLine)
						if len(parts) >= 2 {
							result.NameTH = &parts[0]
							result.SurnameTH = &parts[len(parts)-1]
							result.FieldConfidences["name_th"] = 0.7
						}
					} else {
						// Likely English name
						parts := strings.Fields(nextLine)
						if len(parts) >= 2 {
							result.NameEN = &parts[0]
							result.SurnameEN = &parts[len(parts)-1]
							result.FieldConfidences["name_en"] = 0.7
						}
					}
				}
			}
		}
	}
	return result
}

func (v *VisionServiceImpl) containsThai(text string) bool {
	for _, char := range text {
		if char >= 0x0E00 && char <= 0x0E7F {
			return true
		}
	}
	return false
}

func stringPtr(s string) *string {
	return &s
}