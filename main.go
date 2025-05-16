package main

import (
	"FetchAssessment/models"
	"FetchAssessment/points"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

var receipts = make(map[string]models.Receipt) // Storage for receipts

// ProcessReceipt handles POST requests to /receipts/process
func ProcessReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt models.Receipt

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&receipt); err != nil {
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if receipt.Retailer == "" ||
		receipt.PurchaseDate == "" ||
		receipt.PurchaseTime == "" ||
		len(receipt.Items) == 0 ||
		receipt.Total == "" {
		http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
		return
	}

	for _, item := range receipt.Items {
		if item.ShortDescription == "" || item.Price == "" {
			http.Error(w, "The receipt is invalid.", http.StatusBadRequest)
			return
		}
	}

	// Generate an ID for the receipt
	id := uuid.New().String()
	receipts[id] = receipt

	// Respond with the ID
	response := map[string]string{"id": id}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetPoints handles GET requests to /receipts/{id}/points
func GetPoints(w http.ResponseWriter, r *http.Request) {
	// Extract the receipt ID from the URL path
	path := r.URL.Path
	id := strings.TrimPrefix(path, "/receipts/")
	id = strings.TrimSuffix(id, "/points")

	// Retrieve the receipt by ID
	receipt, found := receipts[id]
	if !found {
		http.Error(w, "No receipt found for that ID.", http.StatusNotFound)
		return
	}

	// Calculate the points
	points := points.CalculatePoints(receipt)

	// Respond with the points
	response := map[string]int{"points": points}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Set up the routes
	http.HandleFunc("/receipts/process", ProcessReceipt)
	http.HandleFunc("/receipts/", GetPoints)

	// Start the server
	port := ":8080"
	fmt.Printf("Application running on port %s...\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
