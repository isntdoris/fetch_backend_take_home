package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/isntdoris/fetch_backend_take_home/api"
)

var mockPointsStore sync.Map

func TestProcessReceipt(t *testing.T) {
	app := fiber.New()
	app.Post("/receipts/process", api.ProcessReceipt)

	// define test cases
	tests := []struct {
		description  string
		receiptInput map[string]interface{}
		expectedCode int
	}{
		{
			description: "Valid receipt with items",
			receiptInput: map[string]interface{}{
				"retailer":     "Test Store",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "15:00",
				"items": []map[string]string{
					{"shortDescription": "Item 1", "price": "10.00"},
				},
				"total": "10.00",
			},
			expectedCode: fiber.StatusOK,
		},
		{
			description: "Receipt missing items",
			receiptInput: map[string]interface{}{
				"retailer":     "Test Store",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "15:00",
				"total":        "10.00",
			},
			expectedCode: fiber.StatusBadRequest,
		},
		{
			description: "Invalid retailer format",
			receiptInput: map[string]interface{}{
				"retailer":     "Test@Store",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "15:00",
				"items": []map[string]interface{}{
					{"shortDescription": "Item 1", "price": "10.00"},
				},
				"total": "10.00",
			},
			expectedCode: fiber.StatusBadRequest,
		},
		{
			description: "Invalid purchase date format",
			receiptInput: map[string]interface{}{
				"retailer":     "Test Store",
				"purchaseDate": "01-01-2022", // Invalid date format
				"purchaseTime": "15:00",
				"items": []map[string]interface{}{
					{"shortDescription": "Item 1", "price": "10.00"},
				},
				"total": "10.00",
			},
			expectedCode: fiber.StatusBadRequest,
		},
		{
			description: "Invalid total format",
			receiptInput: map[string]interface{}{
				"retailer":     "Test Store",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "15:00",
				"items": []map[string]interface{}{
					{"shortDescription": "Item 1", "price": "10.00"},
				},
				"total": "10,00", // Invalid total format
			},
			expectedCode: fiber.StatusBadRequest,
		},
		{
			description: "Empty items array",
			receiptInput: map[string]interface{}{
				"retailer":     "Test Store",
				"purchaseDate": "2022-01-01",
				"purchaseTime": "15:00",
				"items":        []map[string]interface{}{}, // Empty items array
				"total":        "10.00",
			},
			expectedCode: fiber.StatusBadRequest,
		},
	}

	// iterate over test cases
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tc.receiptInput)
			req := httptest.NewRequest("POST", "/receipts/process", bytes.NewReader(payloadBytes))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatalf("%s: Test request failed: %v", tc.description, err)
			}

			if resp.StatusCode != tc.expectedCode {
				t.Errorf("%s: Expected status code %d, got %d", tc.description, tc.expectedCode, resp.StatusCode)
			}
		})
	}
}

func TestGetPoints(t *testing.T) {
	app := fiber.New()
	app.Get("/receipts/:id/points", api.GetPoints)

	// Define test cases
	tests := []struct {
		description  string
		setupFunc    func() string // Function to setup test case, returns receipt ID
		expectedCode int
	}{
		{
			description: "Valid receipt ID with points",
			setupFunc: func() string {
				// Setup: add points for a test receipt ID
				testReceiptID := uuid.New().String()
				mockPointsStore.Store(testReceiptID, 100)
				return testReceiptID
			},
			expectedCode: 200,
		},
		{
			description: "Invalid UUID format",
			setupFunc: func() string {
				// No setup needed, return invalid UUID format
				return "invalid-uuid-format"
			},
			expectedCode: fiber.StatusBadRequest,
		},
		{
			description: "Valid UUID format but not existing in the store",
			setupFunc: func() string {
				// Return a valid but unknown UUID
				return uuid.New().String()
			},
			expectedCode: fiber.StatusNotFound,
		},
	}

	// Iterate over test cases
	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			receiptID := tc.setupFunc() // Execute setup function
			req := httptest.NewRequest("GET", fmt.Sprintf("/receipts/%s/points", receiptID), nil)
			resp, err := app.Test(req, -1)
			if err != nil {
				t.Fatalf("%s: Test request failed: %v", tc.description, err)
			}

			if resp.StatusCode != tc.expectedCode {
				t.Errorf("%s: Expected status code %d, got %d", tc.description, tc.expectedCode, resp.StatusCode)
			}

		})
	}
}
