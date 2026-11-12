package slipverification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"gofiber-template/domain/ports"
)

// SlipOKAdapter — implement SlipVerificationPort ด้วย SlipOK API
type SlipOKAdapter struct {
	branchID string
	apiKey   string
	client   *http.Client
}

// NewSlipOKAdapter สร้าง adapter จาก env: SLIPOK_BRANCH_ID, SLIPOK_API_KEY
func NewSlipOKAdapter() *SlipOKAdapter {
	return &SlipOKAdapter{
		branchID: os.Getenv("SLIPOK_BRANCH_ID"),
		apiKey:   os.Getenv("SLIPOK_API_KEY"),
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// IsConfigured เช็คว่า config ครบหรือยัง
func (a *SlipOKAdapter) IsConfigured() bool {
	return a.branchID != "" && a.apiKey != ""
}

// slipOKResponse — response จาก SlipOK API
type slipOKResponse struct {
	Success bool            `json:"success"`
	Code    int             `json:"code,omitempty"`
	Message string          `json:"message,omitempty"`
	Data    *slipOKData     `json:"data,omitempty"`
}

type slipOKData struct {
	Success        bool           `json:"success"`
	Message        string         `json:"message"`
	Language       string         `json:"language"`
	TransRef       string         `json:"transRef"`
	SendingBank    string         `json:"sendingBank"`
	ReceivingBank  string         `json:"receivingBank"`
	TransDate      string         `json:"transDate"`
	TransTime      string         `json:"transTime"`
	TransTimestamp string         `json:"transTimestamp,omitempty"`
	Sender         slipOKParty    `json:"sender"`
	Receiver       slipOKParty    `json:"receiver"`
	Amount         float64        `json:"amount"`
	CountryCode    string         `json:"countryCode"`
	TransFeeAmount json.Number    `json:"transFeeAmount"`
	Ref1           string         `json:"ref1"`
	Ref2           string         `json:"ref2"`
	Ref3           string         `json:"ref3"`
}

type slipOKParty struct {
	DisplayName string         `json:"displayName"`
	Name        string         `json:"name"`
	Proxy       slipOKProxy    `json:"proxy"`
	Account     slipOKAccount  `json:"account"`
}

type slipOKProxy struct {
	Type  *string `json:"type"`
	Value *string `json:"value"`
}

type slipOKAccount struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// VerifySlip ส่งรูปสลิปไป SlipOK API ตรวจสอบ
func (a *SlipOKAdapter) VerifySlip(ctx context.Context, imageData []byte, amount *float64) (*ports.SlipVerificationResult, error) {
	if !a.IsConfigured() {
		return nil, fmt.Errorf("SlipOK not configured: missing SLIPOK_BRANCH_ID or SLIPOK_API_KEY")
	}

	url := fmt.Sprintf("https://api.slipok.com/api/line/apikey/%s", a.branchID)

	// สร้าง multipart form — ส่งรูปเป็น files
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("files", "slip.jpg")
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := part.Write(imageData); err != nil {
		return nil, fmt.Errorf("failed to write image data: %w", err)
	}

	// ส่ง log: true เพื่อเช็คสลิปซ้ำ + เช็คบัญชี
	if err := writer.WriteField("log", "true"); err != nil {
		return nil, fmt.Errorf("failed to write log field: %w", err)
	}

	// ส่ง amount ถ้ามี
	if amount != nil && *amount > 0 {
		if err := writer.WriteField("amount", fmt.Sprintf("%.2f", *amount)); err != nil {
			return nil, fmt.Errorf("failed to write amount field: %w", err)
		}
	}

	writer.Close()

	// สร้าง request
	req, err := http.NewRequestWithContext(ctx, "POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("x-authorization", a.apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// ส่ง request
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("SlipOK request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read SlipOK response: %w", err)
	}

	log.Printf("[SlipOK] Response status=%d body=%s", resp.StatusCode, string(respBody))

	// Parse response
	var slipResp slipOKResponse
	if err := json.Unmarshal(respBody, &slipResp); err != nil {
		return nil, fmt.Errorf("failed to parse SlipOK response: %w", err)
	}

	result := &ports.SlipVerificationResult{
		RawJSON: string(respBody),
	}

	// Error case
	if resp.StatusCode != 200 || !slipResp.Success {
		result.Success = false
		result.ErrorCode = slipResp.Code
		result.Message = slipResp.Message

		// กรณี 1013 (ยอดไม่ตรง) ยังมี data กลับมา
		if slipResp.Data != nil {
			a.mapDataToResult(slipResp.Data, result)
		}
		return result, nil
	}

	// Success case
	if slipResp.Data == nil {
		return nil, fmt.Errorf("SlipOK returned success but no data")
	}

	result.Success = true
	result.Message = slipResp.Data.Message
	a.mapDataToResult(slipResp.Data, result)

	return result, nil
}

// mapDataToResult แปลง SlipOK data → SlipVerificationResult
func (a *SlipOKAdapter) mapDataToResult(data *slipOKData, result *ports.SlipVerificationResult) {
	result.TransRef = data.TransRef
	result.TransDate = data.TransDate
	result.TransTime = data.TransTime
	result.Amount = data.Amount
	result.SendingBank = data.SendingBank
	result.ReceivingBank = data.ReceivingBank
	result.CountryCode = data.CountryCode
	result.Ref1 = data.Ref1
	result.Ref2 = data.Ref2
	result.Ref3 = data.Ref3

	if fee, err := data.TransFeeAmount.Float64(); err == nil {
		result.TransFeeAmount = fee
	}

	// Parse timestamp
	if data.TransTimestamp != "" {
		if t, err := time.Parse(time.RFC3339Nano, data.TransTimestamp); err == nil {
			result.TransTimestamp = t
		}
	} else if data.TransDate != "" && data.TransTime != "" {
		// Fallback: parse จาก transDate + transTime
		if t, err := time.Parse("20060102 15:04:05", data.TransDate+" "+data.TransTime); err == nil {
			result.TransTimestamp = t
		}
	}

	// Sender
	result.SenderDisplayName = data.Sender.DisplayName
	result.SenderName = data.Sender.Name
	result.SenderAccount = data.Sender.Account.Value

	// Receiver
	result.ReceiverDisplayName = data.Receiver.DisplayName
	result.ReceiverName = data.Receiver.Name
	result.ReceiverAccount = data.Receiver.Account.Value
}

// CheckQuota เช็คโควต้าคงเหลือ
func (a *SlipOKAdapter) CheckQuota(ctx context.Context) (*ports.SlipQuotaResult, error) {
	if !a.IsConfigured() {
		return nil, fmt.Errorf("SlipOK not configured")
	}

	url := fmt.Sprintf("https://api.slipok.com/api/line/apikey/%s/quota", a.branchID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("x-authorization", a.apiKey)

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("SlipOK quota request failed: %w", err)
	}
	defer resp.Body.Close()

	var quotaResp struct {
		Success bool `json:"success"`
		Data    struct {
			Quota        int    `json:"quota"`
			OverQuota    int    `json:"overQuota"`
			SpecialQuota int    `json:"specialQuota"`
			EndDate      string `json:"endDate"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&quotaResp); err != nil {
		return nil, fmt.Errorf("failed to parse quota response: %w", err)
	}

	return &ports.SlipQuotaResult{
		Quota:        quotaResp.Data.Quota,
		OverQuota:    quotaResp.Data.OverQuota,
		SpecialQuota: quotaResp.Data.SpecialQuota,
		EndDate:      quotaResp.Data.EndDate,
	}, nil
}
