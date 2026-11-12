package dto

import "time"

// VisionOCRResponse represents the raw OCR response from Google Vision API
type VisionOCRResponse struct {
	Text       string             `json:"text"`
	Confidence float32            `json:"confidence"`
	Words      []VisionWord       `json:"words,omitempty"`
	Blocks     []VisionTextBlock  `json:"blocks,omitempty"`
}

// VisionWord represents a single word detected by Vision API
type VisionWord struct {
	Text       string    `json:"text"`
	Confidence float32   `json:"confidence"`
	BoundingBox *BoundingBox `json:"bounding_box,omitempty"`
}

// VisionTextBlock represents a block of text detected by Vision API
type VisionTextBlock struct {
	Text       string    `json:"text"`
	Confidence float32   `json:"confidence"`
	BoundingBox *BoundingBox `json:"bounding_box,omitempty"`
}

// BoundingBox represents the coordinates of detected text
type BoundingBox struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// PassportExtractionResult represents extracted passport data
type PassportExtractionResult struct {
	// OCR Information
	RawText       string  `json:"raw_text"`
	OverallConfidence float32 `json:"overall_confidence"`
	ProcessedSuccessfully bool `json:"processed_successfully"`
	ErrorMessage  *string `json:"error_message,omitempty"`

	// Extracted Passport Data
	PassportNumber   *string    `json:"passport_number,omitempty"`
	PassportType     *string    `json:"passport_type,omitempty"`
	CountryCode      *string    `json:"country_code,omitempty"`

	// Personal Information (Thai)
	TitleTH      *string    `json:"title_th,omitempty"`
	NameTH       *string    `json:"name_th,omitempty"`
	SurnameTH    *string    `json:"surname_th,omitempty"`

	// Personal Information (English)
	TitleEN      *string    `json:"title_en,omitempty"`
	NameEN       *string    `json:"name_en,omitempty"`
	SurnameEN    *string    `json:"surname_en,omitempty"`

	// Additional Information
	Gender         *string    `json:"gender,omitempty"`
	Nationality    *string    `json:"nationality,omitempty"`
	DateOfBirth    *time.Time `json:"date_of_birth,omitempty"`
	PlaceOfBirth   *string    `json:"place_of_birth,omitempty"`
	DateOfIssue    *time.Time `json:"date_of_issue,omitempty"`
	DateOfExpiry   *time.Time `json:"date_of_expiry,omitempty"`
	IssuingAuthority *string  `json:"issuing_authority,omitempty"`
	IDNumber       *string    `json:"id_number,omitempty"`

	// Confidence Scores for individual fields
	FieldConfidences map[string]float32 `json:"field_confidences,omitempty"`
}

// PassportProcessingRequest represents request for passport processing
type PassportProcessingRequest struct {
	ImageData    []byte  `json:"image_data" validate:"required"`
	BookingID    string  `json:"booking_id" validate:"required"`
	FileName     string  `json:"file_name" validate:"required"`
	ProcessOCR   bool    `json:"process_ocr"`
}

// PassportProcessingResponse represents the full passport processing response
type PassportProcessingResponse struct {
	PassportID       string                    `json:"passport_id"`
	ImageURL         string                    `json:"image_url"`
	ProcessingStatus string                    `json:"processing_status"`
	ExtractionResult *PassportExtractionResult `json:"extraction_result,omitempty"`
	CreatedAt        string                    `json:"created_at"`
	UpdatedAt        string                    `json:"updated_at"`
}