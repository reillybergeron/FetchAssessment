package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"FetchAssessment/models"
)

func validReceipt() models.Receipt {
	return models.Receipt{
		Retailer:     "Reilly",
		PurchaseDate: "2000-10-05",
		PurchaseTime: "15:00",
		Total:        "6.25",
		Items: []models.Item{
			{ShortDescription: "Pen", Price: "5.00"},
			{ShortDescription: "Paper", Price: "1.25"},
		},
	}
}

func TestProcessReceipt_Valid(t *testing.T) {
	receipt := validReceipt()

	body, _ := json.Marshal(receipt)
	req := httptest.NewRequest(http.MethodPost, "/receipts/process", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	ProcessReceipt(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal("response is not valid JSON")
	}
	if _, ok := resp["id"]; !ok {
		t.Error("response does not contain an 'id' field")
	}
}

func TestProcessReceipt_Invalid(t *testing.T) {
	invalidJSON := `{"Retailer": "Reilly"}`
	req := httptest.NewRequest(http.MethodPost, "/receipts/process", strings.NewReader(invalidJSON))
	rec := httptest.NewRecorder()

	ProcessReceipt(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", rec.Code)
	}
}

func TestGetPoints_ValidID(t *testing.T) {
	receipt := validReceipt()
	body, _ := json.Marshal(receipt)
	req := httptest.NewRequest(http.MethodPost, "/receipts/process", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	ProcessReceipt(rec, req)

	var resp map[string]string
	json.Unmarshal(rec.Body.Bytes(), &resp)
	id := resp["id"]

	pointsReq := httptest.NewRequest(http.MethodGet, "/receipts/"+id+"/points", nil)
	pointsRec := httptest.NewRecorder()
	GetPoints(pointsRec, pointsReq)

	if pointsRec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", pointsRec.Code)
	}

	var pointsResp map[string]int
	if err := json.Unmarshal(pointsRec.Body.Bytes(), &pointsResp); err != nil {
		t.Fatal("response is not valid JSON")
	}
	if _, ok := pointsResp["points"]; !ok {
		t.Error("response does not contain a 'points' field")
	}
}

func TestGetPoints_InvalidID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/receipts/invalid-id/points", nil)
	rec := httptest.NewRecorder()

	GetPoints(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}
}
