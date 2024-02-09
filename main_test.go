package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"net/http/httptest"
	"testing"
)

func TestProcessReceipt(t *testing.T) {
	app := fiber.New()
	app.Post("/receipts/process", processReceipt)

	receiptPayload := map[string]interface{}{
		"retailer":     "Test Store",
		"purchaseDate": "2022-01-01",
		"purchaseTime": "15:00",
		"items": []map[string]string{
			{"shortDescription": "Item 1", "price": "10.00"},
		},
		"total": "10.00",
	}
	payloadBytes, _ := json.Marshal(receiptPayload)
	req := httptest.NewRequest("POST", "/receipts/process", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		t.Errorf("TestProcessReceipt failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %v", resp.StatusCode)
	}

	// further validation
}

func TestGetPoints(t *testing.T) {
	app := fiber.New()
	app.Get("/receipts/:id/points", getPoints)

	// add points for a test receipt ID
	testReceiptID := uuid.New().String()
	pointsStore[testReceiptID] = 100

	req := httptest.NewRequest("GET", fmt.Sprintf("/receipts/%s/points", testReceiptID), nil)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Errorf("TestGetPoints failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %v", resp.StatusCode)
	}

	// further validation
}
