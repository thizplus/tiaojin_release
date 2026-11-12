package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type TestCase struct {
	Name        string
	Description string
	FilePath    string
	Expected    ExpectedResults
}

type ExpectedResults struct {
	SumRoom               float64
	SumRate               float64
	SumAdditionalServices float64
	SumNetBeforeDiscount  float64
	TotalCommission       float64
	SumNetAfterDiscount   float64
	SumDeposit            float64
	TotalPassengers       int
}

type BookingResponse struct {
	ID                 string  `json:"id"`
	InvoiceID          string  `json:"invoice_id"`
	CustomerName       string  `json:"customer_name"`
	SumRoom            float64 `json:"sum_room"`
	SumRate            float64 `json:"sum_rate"`
	SumNet             float64 `json:"sum_net"`
	SumCommission      float64 `json:"sum_commission"`
	Discount           float64 `json:"discount"`
	SumDeposit         float64 `json:"sum_deposit"`
	Status             string  `json:"status"`
	RateTypes          []interface{} `json:"rate_types"`
	AdditionalServices []interface{} `json:"additional_services"`
}

func main() {
	fmt.Println("🚀 Comprehensive Booking Test Runner")
	fmt.Println("=====================================")

	baseURL := "http://localhost:8080/api/v1" // Adjust to your API URL
	userID := "880e8400-e29b-41d4-a716-446655440001" // Adjust to valid user ID

	testCases := []TestCase{
		{
			Name:        "Test Case 1: Basic Family Booking",
			Description: "2 Adults + 1 Child with bed, Standard room, No additional services",
			FilePath:    "test_case_1_basic_family.json",
			Expected: ExpectedResults{
				SumRoom:               0.00,
				SumRate:               24500.00,
				SumAdditionalServices: 0.00,
				SumNetBeforeDiscount:  24500.00,
				TotalCommission:       2450.00,
				SumNetAfterDiscount:   22050.00,
				SumDeposit:            6615.00,
				TotalPassengers:       3,
			},
		},
		{
			Name:        "Test Case 2: Exclude Logic Testing",
			Description: "2 Adults + 1 Child + 1 Infant (exclude=true)",
			FilePath:    "test_case_2_exclude_logic.json",
			Expected: ExpectedResults{
				SumRoom:               0.00,
				SumRate:               24500.00,
				SumAdditionalServices: 0.00,
				SumNetBeforeDiscount:  24500.00,
				TotalCommission:       2450.00,
				SumNetAfterDiscount:   22050.00,
				SumDeposit:            6615.00,
				TotalPassengers:       4,
			},
		},
		{
			Name:        "Test Case 3: Additional Services - Per Unit",
			Description: "Family + Travel Insurance + Spa Treatment",
			FilePath:    "test_case_3_per_unit_services.json",
			Expected: ExpectedResults{
				SumRoom:               0.00,
				SumRate:               24500.00,
				SumAdditionalServices: 1150.00,
				SumNetBeforeDiscount:  25650.00,
				TotalCommission:       2450.00,
				SumNetAfterDiscount:   23200.00,
				SumDeposit:            6960.00,
				TotalPassengers:       3,
			},
		},
		{
			Name:        "Test Case 4: Additional Services - Per Booking",
			Description: "Family + Airport Transfer + Photography",
			FilePath:    "test_case_4_per_booking_services.json",
			Expected: ExpectedResults{
				SumRoom:               0.00,
				SumRate:               24500.00,
				SumAdditionalServices: 2000.00,
				SumNetBeforeDiscount:  26500.00,
				TotalCommission:       2450.00,
				SumNetAfterDiscount:   24050.00,
				SumDeposit:            7215.00,
				TotalPassengers:       3,
			},
		},
		{
			Name:        "Test Case 5: Mixed Additional Services",
			Description: "Combination of per_unit and per_booking services",
			FilePath:    "test_case_5_mixed_services.json",
			Expected: ExpectedResults{
				SumRoom:               0.00,
				SumRate:               32000.00,
				SumAdditionalServices: 3000.00,
				SumNetBeforeDiscount:  35000.00,
				TotalCommission:       3200.00,
				SumNetAfterDiscount:   31800.00,
				SumDeposit:            9540.00,
				TotalPassengers:       4,
			},
		},
		{
			Name:        "Test Case 6: Room Upgrades",
			Description: "Different room types with prices, Single supplement",
			FilePath:    "test_case_6_room_upgrades.json",
			Expected: ExpectedResults{
				SumRoom:               2000.00,
				SumRate:               25500.00,
				SumAdditionalServices: 800.00,
				SumNetBeforeDiscount:  28300.00,
				TotalCommission:       2550.00,
				SumNetAfterDiscount:   25750.00,
				SumDeposit:            7725.00,
				TotalPassengers:       3,
			},
		},
		{
			Name:        "Test Case 7: Maximum Complexity",
			Description: "Multiple passenger types, room types, additional services",
			FilePath:    "test_case_7_maximum_complexity.json",
			Expected: ExpectedResults{
				SumRoom:               2000.00,
				SumRate:               55500.00,
				SumAdditionalServices: 7150.00,
				SumNetBeforeDiscount:  64650.00,
				TotalCommission:       5550.00,
				SumNetAfterDiscount:   59100.00,
				SumDeposit:            17730.00,
				TotalPassengers:       8,
			},
		},
	}

	fmt.Printf("Testing against API: %s\n", baseURL)
	fmt.Printf("Using User ID: %s\n\n", userID)

	passed := 0
	failed := 0

	for _, testCase := range testCases {
		fmt.Printf("📋 %s\n", testCase.Name)
		fmt.Printf("   %s\n", testCase.Description)
		fmt.Printf("   File: %s\n", testCase.FilePath)

		// Read test data
		requestData, err := os.ReadFile(testCase.FilePath)
		if err != nil {
			fmt.Printf("   ❌ Failed to read test file: %v\n\n", err)
			failed++
			continue
		}

		// Make API request
		url := fmt.Sprintf("%s/users/%s/bookings", baseURL, userID)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestData))
		if err != nil {
			fmt.Printf("   ❌ API request failed: %v\n\n", err)
			failed++
			continue
		}
		defer resp.Body.Close()

		// Read response
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("   ❌ Failed to read response: %v\n\n", err)
			failed++
			continue
		}

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			fmt.Printf("   ❌ API returned status %d: %s\n\n", resp.StatusCode, string(responseBody))
			failed++
			continue
		}

		// Parse response
		var bookingResponse BookingResponse
		if err := json.Unmarshal(responseBody, &bookingResponse); err != nil {
			fmt.Printf("   ❌ Failed to parse response: %v\n\n", err)
			failed++
			continue
		}

		// Validate results
		tolerance := 0.01 // Allow small floating point differences
		success := true

		if !floatEqual(bookingResponse.SumRoom, testCase.Expected.SumRoom, tolerance) {
			fmt.Printf("   ❌ SumRoom mismatch: expected %.2f, got %.2f\n", testCase.Expected.SumRoom, bookingResponse.SumRoom)
			success = false
		}

		if !floatEqual(bookingResponse.SumRate, testCase.Expected.SumRate, tolerance) {
			fmt.Printf("   ❌ SumRate mismatch: expected %.2f, got %.2f\n", testCase.Expected.SumRate, bookingResponse.SumRate)
			success = false
		}

		expectedAdditionalServices := testCase.Expected.SumAdditionalServices
		// Calculate actual additional services from response
		actualAdditionalServices := bookingResponse.SumNet - bookingResponse.SumRoom - bookingResponse.SumRate + bookingResponse.Discount
		if !floatEqual(actualAdditionalServices, expectedAdditionalServices, tolerance) {
			fmt.Printf("   ❌ AdditionalServices mismatch: expected %.2f, calculated %.2f\n", expectedAdditionalServices, actualAdditionalServices)
			success = false
		}

		if !floatEqual(bookingResponse.Discount, testCase.Expected.TotalCommission, tolerance) {
			fmt.Printf("   ❌ Commission/Discount mismatch: expected %.2f, got %.2f\n", testCase.Expected.TotalCommission, bookingResponse.Discount)
			success = false
		}

		expectedNet := testCase.Expected.SumNetBeforeDiscount - testCase.Expected.TotalCommission
		if !floatEqual(bookingResponse.SumNet, expectedNet, tolerance) {
			fmt.Printf("   ❌ SumNet mismatch: expected %.2f, got %.2f\n", expectedNet, bookingResponse.SumNet)
			success = false
		}

		if !floatEqual(bookingResponse.SumDeposit, testCase.Expected.SumDeposit, tolerance) {
			fmt.Printf("   ❌ SumDeposit mismatch: expected %.2f, got %.2f\n", testCase.Expected.SumDeposit, bookingResponse.SumDeposit)
			success = false
		}

		if success {
			fmt.Printf("   ✅ Test passed!\n")
			fmt.Printf("      InvoiceID: %s\n", bookingResponse.InvoiceID)
			fmt.Printf("      SumRoom: %.2f, SumRate: %.2f, SumNet: %.2f\n",
				bookingResponse.SumRoom, bookingResponse.SumRate, bookingResponse.SumNet)
			fmt.Printf("      Commission: %.2f, Deposit: %.2f\n",
				bookingResponse.Discount, bookingResponse.SumDeposit)
			passed++
		} else {
			failed++
		}

		fmt.Println()
	}

	fmt.Println("🏁 Test Summary")
	fmt.Println("===============")
	fmt.Printf("Total tests: %d\n", len(testCases))
	fmt.Printf("Passed: %d ✅\n", passed)
	fmt.Printf("Failed: %d ❌\n", failed)

	if failed == 0 {
		fmt.Println("\n🎉 All tests passed! Booking system is working correctly.")
	} else {
		fmt.Printf("\n⚠️  %d test(s) failed. Please check the implementation.\n", failed)
	}
}

func floatEqual(a, b, tolerance float64) bool {
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff <= tolerance
}

// Helper function to convert string to float64
func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return f
}