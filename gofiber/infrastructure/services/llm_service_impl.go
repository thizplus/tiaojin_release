package services

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"gofiber-template/domain/dto"
	"gofiber-template/domain/services"
	"gofiber-template/pkg/config"
)

type LLMServiceImpl struct {
	config *config.Config
	client *http.Client
	cache  *ExtractionCache
}

type CacheEntry struct {
	Result    *dto.PassportExtractionResult
	Timestamp time.Time
}

type ExtractionCache struct {
	entries map[string]CacheEntry
	mutex   sync.RWMutex
	ttl     time.Duration
}

type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
	Error      *GeminiError      `json:"error,omitempty"`
}

type GeminiCandidate struct {
	Content GeminiContent `json:"content"`
}

type GeminiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type PassportData struct {
	PassportNumber   *string `json:"passport_number,omitempty"`
	TitleTH          *string `json:"title_th,omitempty"`
	NameTH           *string `json:"name_th,omitempty"`
	SurnameTH        *string `json:"surname_th,omitempty"`
	TitleEN          *string `json:"title_en,omitempty"`
	NameEN           *string `json:"name_en,omitempty"`
	SurnameEN        *string `json:"surname_en,omitempty"`
	DateOfBirth      *string `json:"date_of_birth,omitempty"`
	DateOfExpiry     *string `json:"date_of_expiry,omitempty"`
	DateOfIssue      *string `json:"date_of_issue,omitempty"`
	Nationality      *string `json:"nationality,omitempty"`
	CountryCode      *string `json:"country_code,omitempty"`
	Gender           *string `json:"gender,omitempty"`
	IDNumber         *string `json:"id_number,omitempty"`
	PlaceOfBirth     *string `json:"place_of_birth,omitempty"`
	IssuingAuthority *string `json:"issuing_authority,omitempty"`
	Confidence       float32 `json:"confidence"`
}

func NewLLMService(cfg *config.Config) (services.LLMService, error) {
	if !cfg.LLM.Enabled {
		return nil, fmt.Errorf("LLM service is disabled")
	}

	if cfg.LLM.APIKey == "" {
		return nil, fmt.Errorf("LLM API key is required")
	}

	cache := &ExtractionCache{
		entries: make(map[string]CacheEntry),
		ttl:     1 * time.Hour, // Cache results for 1 hour
	}

	// Start cache cleanup goroutine
	go cache.cleanup()

	return &LLMServiceImpl{
		config: cfg,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		cache: cache,
	}, nil
}

func (l *LLMServiceImpl) ExtractPassportFields(ctx context.Context, ocrText string) (*dto.PassportExtractionResult, error) {
	log.Printf("[INFO] Starting LLM-based passport field extraction")

	// Check cache first
	cacheKey := l.generateCacheKey(ocrText)
	if cachedResult := l.cache.Get(cacheKey); cachedResult != nil {
		log.Printf("[INFO] Returning cached result for OCR text")
		return cachedResult, nil
	}

	prompt := l.buildPassportExtractionPrompt(ocrText)

	response, err := l.ParseStructuredData(ctx, ocrText, prompt)
	if err != nil {
		log.Printf("[ERROR] LLM API call failed: %v", err)
		return nil, fmt.Errorf("LLM API call failed: %w", err)
	}

	// Parse JSON response
	var passportData PassportData
	if err := json.Unmarshal([]byte(response), &passportData); err != nil {
		log.Printf("[ERROR] Failed to parse LLM response: %v", err)
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}

	// Convert to PassportExtractionResult
	result := &dto.PassportExtractionResult{
		RawText:               ocrText,
		ProcessedSuccessfully: true,
		OverallConfidence:     passportData.Confidence,
		FieldConfidences:      make(map[string]float32),
		PassportNumber:        passportData.PassportNumber,
		TitleTH:               passportData.TitleTH,
		NameTH:                passportData.NameTH,
		SurnameTH:             passportData.SurnameTH,
		TitleEN:               passportData.TitleEN,
		NameEN:                passportData.NameEN,
		SurnameEN:             passportData.SurnameEN,
		Nationality:           passportData.Nationality,
		CountryCode:           passportData.CountryCode,
		Gender:                passportData.Gender,
		PlaceOfBirth:          passportData.PlaceOfBirth,
		IssuingAuthority:      passportData.IssuingAuthority,
		IDNumber:              passportData.IDNumber,
	}

	// Parse dates with multiple format support
	if passportData.DateOfBirth != nil && *passportData.DateOfBirth != "" {
		if birthDate, err := time.Parse("2006-01-02", *passportData.DateOfBirth); err == nil {
			result.DateOfBirth = &birthDate
			result.FieldConfidences["date_of_birth"] = 0.9
		}
	}

	if passportData.DateOfExpiry != nil && *passportData.DateOfExpiry != "" {
		if expiryDate, err := time.Parse("2006-01-02", *passportData.DateOfExpiry); err == nil {
			result.DateOfExpiry = &expiryDate
			result.FieldConfidences["date_of_expiry"] = 0.9
		}
	}

	if passportData.DateOfIssue != nil && *passportData.DateOfIssue != "" {
		if issueDate, err := time.Parse("2006-01-02", *passportData.DateOfIssue); err == nil {
			result.DateOfIssue = &issueDate
			result.FieldConfidences["date_of_issue"] = 0.9
		}
	}

	// Set confidence scores for all extracted fields
	if result.PassportNumber != nil {
		result.FieldConfidences["passport_number"] = 0.95
	}
	if result.TitleTH != nil {
		result.FieldConfidences["title_th"] = 0.85
	}
	if result.NameTH != nil {
		result.FieldConfidences["name_th"] = 0.9
	}
	if result.SurnameTH != nil {
		result.FieldConfidences["surname_th"] = 0.9
	}
	if result.TitleEN != nil {
		result.FieldConfidences["title_en"] = 0.85
	}
	if result.NameEN != nil {
		result.FieldConfidences["name_en"] = 0.9
	}
	if result.SurnameEN != nil {
		result.FieldConfidences["surname_en"] = 0.9
	}
	if result.Nationality != nil {
		result.FieldConfidences["nationality"] = 0.95
	}
	if result.CountryCode != nil {
		result.FieldConfidences["country_code"] = 0.95
	}
	if result.Gender != nil {
		result.FieldConfidences["gender"] = 0.9
	}
	if result.PlaceOfBirth != nil {
		result.FieldConfidences["place_of_birth"] = 0.8
	}
	if result.IssuingAuthority != nil {
		result.FieldConfidences["issuing_authority"] = 0.85
	}
	if result.IDNumber != nil {
		result.FieldConfidences["id_number"] = 0.9
	}

	log.Printf("[SUCCESS] LLM extraction completed, fields: %d, confidence: %.2f",
		len(result.FieldConfidences), result.OverallConfidence)

	// Cache the result
	l.cache.Set(cacheKey, result)

	return result, nil
}

func (l *LLMServiceImpl) ParseStructuredData(ctx context.Context, text string, prompt string) (string, error) {
	switch l.config.LLM.Provider {
	case "gemini":
		return l.callGeminiAPI(ctx, prompt)
	default:
		return "", fmt.Errorf("unsupported LLM provider: %s", l.config.LLM.Provider)
	}
}

func (l *LLMServiceImpl) callGeminiAPI(ctx context.Context, prompt string) (string, error) {
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		l.config.LLM.Model, l.config.LLM.APIKey)

	request := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: prompt},
				},
			},
		},
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := l.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Gemini API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Gemini API error %d: %s", resp.StatusCode, string(body))
	}

	var response GeminiResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if response.Error != nil {
		return "", fmt.Errorf("Gemini API error: %s", response.Error.Message)
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from Gemini API")
	}

	responseText := response.Candidates[0].Content.Parts[0].Text

	// Extract JSON from response (remove any markdown formatting)
	jsonStart := strings.Index(responseText, "{")
	jsonEnd := strings.LastIndex(responseText, "}")

	if jsonStart == -1 || jsonEnd == -1 {
		return "", fmt.Errorf("no valid JSON found in response: %s", responseText)
	}

	return responseText[jsonStart : jsonEnd+1], nil
}

func (l *LLMServiceImpl) buildPassportExtractionPrompt(ocrText string) string {
	return fmt.Sprintf(`Extract ALL passport information from the following OCR text and return ONLY a valid JSON object with ALL requested fields. Do not include any explanations or markdown formatting.

OCR Text:
%s

Extract ALL these fields and return as JSON (use null for missing fields):

BASIC INFORMATION:
- passport_number: The passport number (usually starts with letter + numbers, e.g., "G1234567", "AA1234567")

THAI NAMES (if available):
- title_th: Thai title (e.g., "นาย", "นาง", "นางสาว", "เด็กชาย", "เด็กหญิง")
- name_th: Thai first name only (separate from title)
- surname_th: Thai surname/family name

ENGLISH NAMES:
- title_en: English title (e.g., "MR", "MRS", "MISS", "MASTER")
- name_en: English first name only (separate from title)
- surname_en: English surname/family name

PERSONAL DETAILS:
- gender: "M" for male, "F" for female
- date_of_birth: Birth date in YYYY-MM-DD format
- place_of_birth: Place of birth (city/province)
- id_number: National ID number (13 digits for Thai)

PASSPORT DETAILS:
- nationality: Full nationality name (e.g., "THAI", "AMERICAN")
- country_code: 3-letter country code (e.g., "THA", "USA")
- date_of_issue: Issue date in YYYY-MM-DD format
- date_of_expiry: Expiry date in YYYY-MM-DD format
- issuing_authority: Issuing authority (e.g., "MINISTRY OF FOREIGN AFFAIRS")

PROCESSING:
- confidence: Overall confidence score as decimal between 0.0 and 1.0

IMPORTANT EXTRACTION RULES:
1. For Thai names: Separate title from name (e.g., "นายสมชาย" -> title_th: "นาย", name_th: "สมชาย")
2. For English names: Separate title from name (e.g., "MR SOMCHAI" -> title_en: "MR", name_en: "SOMCHAI")
3. Look for dates in various formats: DD/MM/YYYY, DD-MM-YYYY, DD MMM YYYY, etc.
4. Convert all dates to YYYY-MM-DD format
5. Use null for missing fields, never empty strings
6. Gender: Look for "M", "F", "MALE", "FEMALE", "ชาย", "หญิง"
7. Passport numbers often have specific formats per country
8. Confidence should reflect extraction certainty (0.7+ for good extraction)

Example response format:
{
  "passport_number": "G1234567",
  "title_th": "นาย",
  "name_th": "สมชาย",
  "surname_th": "ใจดี",
  "title_en": "MR",
  "name_en": "SOMCHAI",
  "surname_en": "JAIDEE",
  "gender": "M",
  "date_of_birth": "1990-05-15",
  "place_of_birth": "Bangkok",
  "id_number": "1234567890123",
  "nationality": "THAI",
  "country_code": "THA",
  "date_of_issue": "2020-05-16",
  "date_of_expiry": "2030-05-15",
  "issuing_authority": "Ministry of Foreign Affairs",
  "confidence": 0.92
}`, ocrText)
}

// Cache methods
func (l *LLMServiceImpl) generateCacheKey(ocrText string) string {
	hash := md5.Sum([]byte(strings.TrimSpace(ocrText)))
	return hex.EncodeToString(hash[:])
}

func (c *ExtractionCache) Get(key string) *dto.PassportExtractionResult {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	entry, exists := c.entries[key]
	if !exists {
		return nil
	}

	// Check if entry is still valid
	if time.Since(entry.Timestamp) > c.ttl {
		delete(c.entries, key)
		return nil
	}

	return entry.Result
}

func (c *ExtractionCache) Set(key string, result *dto.PassportExtractionResult) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries[key] = CacheEntry{
		Result:    result,
		Timestamp: time.Now(),
	}
}

func (c *ExtractionCache) cleanup() {
	ticker := time.NewTicker(30 * time.Minute) // Cleanup every 30 minutes
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()

		for key, entry := range c.entries {
			if now.Sub(entry.Timestamp) > c.ttl {
				delete(c.entries, key)
			}
		}
		c.mutex.Unlock()

		log.Printf("[INFO] Cache cleanup completed, entries: %d", len(c.entries))
	}
}