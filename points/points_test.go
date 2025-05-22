package points

import (
	"testing"

	"FetchAssessment/models"
)

func makeItem(desc, price string) models.Item {
	return models.Item{
		ShortDescription: desc,
		Price:            price,
	}
}

func TestPointsForRetailerName(t *testing.T) {
	tests := []struct {
		retailer string
		expected int
	}{
		{"Reilly123", 9},   // all alphanumeric
		{"Rei lly!", 6},    // excludes space and !
		{"  ", 0},          // only spaces
		{"$#@!%^", 0},      // no alphanumeric
		{"Rei-lly_123", 9}, // only letters/numbers counted
	}

	for _, tt := range tests {
		got := pointsForRetailerName(tt.retailer)
		if got != tt.expected {
			t.Errorf("pointsForRetailerName(%q) = %d when it should be %d", tt.retailer, got, tt.expected)
		}
	}
}

func TestPointsForWholeDollar(t *testing.T) {
	tests := []struct {
		total    string
		expected int
	}{
		{"10.00", 50},
		{"9.99", 0},
		{"10", 50},
		{"reilly", 0},
		{"0", 50},
		{" ", 0},
	}

	for _, tt := range tests {
		got := pointsForWholeDollar(tt.total)
		if got != tt.expected {
			t.Errorf("pointsForWholeDollar(%q) = %d when it should be %d", tt.total, got, tt.expected)
		}
	}
}

func TestPointsForMultipleOfQuarter(t *testing.T) {
	tests := []struct {
		total    string
		expected int
	}{
		{"0.25", 25},
		{"0.50", 25},
		{"0.75", 25},
		{"1.00", 25},
		{"1.10", 0},
		{"reilly", 0},
		{"0", 25},
		{" ", 0},
	}

	for _, tt := range tests {
		got := pointsForMultipleOfQuarter(tt.total)
		if got != tt.expected {
			t.Errorf("pointsForMultipleOfQuarter(%q) = %d when it should be %d", tt.total, got, tt.expected)
		}
	}
}

func TestPointsForItemCount(t *testing.T) {
	tests := []struct {
		items    []models.Item
		expected int
	}{
		{[]models.Item{}, 0},
		{[]models.Item{makeItem("item1", "1")}, 0},
		{[]models.Item{makeItem("item1", "1"), makeItem("item2", "1")}, 5},
		{[]models.Item{makeItem("item1", "1"), makeItem("item2", "1"), makeItem("item3", "1")}, 5},
		{[]models.Item{makeItem("item1", "1"), makeItem("item2", "1"), makeItem("item3", "1"), makeItem("item4", "1")}, 10},
	}

	for _, tt := range tests {
		got := pointsForItemCount(tt.items)
		if got != tt.expected {
			t.Errorf("pointsForItemCount(%d items) = %d when it should be %d", len(tt.items), got, tt.expected)
		}
	}
}

func TestPointsForItemDescriptions(t *testing.T) {
	tests := []struct {
		items    []models.Item
		expected int
	}{
		{[]models.Item{}, 0},
		{[]models.Item{makeItem("r", "10.00")}, 0},
		{[]models.Item{makeItem("rei", "10.00")}, 2},
		{[]models.Item{makeItem("reilly", "20.00")}, 4},
		{[]models.Item{makeItem("reillyy", "10.00")}, 0},
		{
			[]models.Item{
				makeItem("rei", "10.00"),      // 2 points
				makeItem("reilly", "5.50"),    // 2 points (due to rounding up)
				makeItem("reilly", "4.50"),    // 1 point (due to rounding down)
				makeItem("re", "3.00"),        // 0 points
				makeItem("   rei   ", "6.00"), // 2 points (due to rounding up)
				makeItem("reilly", "reilly"),  // 0 points
			},
			7,
		},
	}

	for _, tt := range tests {
		got := pointsForItemDescriptions(tt.items)
		if got != tt.expected {
			t.Errorf("pointsForItemDescriptions(%d items) = %d when it should be %d", len(tt.items), got, tt.expected)
		}
	}
}

func TestPointsForOddPurchaseDay(t *testing.T) {
	tests := []struct {
		date     string
		expected int
	}{
		{"2025-05-21", 6},
		{"2025-05-20", 0},
		{"invalid-date", 0},
		{"5-21-2025", 0}, // Incorrect Format
	}

	for _, tt := range tests {
		got := pointsForOddPurchaseDay(tt.date)
		if got != tt.expected {
			t.Errorf("pointsForOddPurchaseDay(%q) = %d when it should be %d", tt.date, got, tt.expected)
		}
	}
}

func TestPointsForPurchaseTime(t *testing.T) {
	tests := []struct {
		timeStr  string
		expected int
	}{
		{"14:00", 0},
		{"14:01", 10},
		{"15:59", 10},
		{"16:00", 0},
		{"invalid", 0},
	}

	for _, tt := range tests {
		got := pointsForPurchaseTime(tt.timeStr)
		if got != tt.expected {
			t.Errorf("pointsForPurchaseTime(%q) = %dwhen it should be %d", tt.timeStr, got, tt.expected)
		}
	}
}

func TestCalculatePoints(t *testing.T) {
	receipt := models.Receipt{
		Retailer:     "reilly-reilly_123",
		PurchaseDate: "2000-10-05",
		PurchaseTime: "14:30",
		Total:        "20.00",
		Items: []models.Item{
			makeItem("rei", "10.00"),
			makeItem("item", "5.00"),
			makeItem("lly", "6.00"),
		},
	}

	// Retailer: "reilly-reilly_123" length 15 alphanumeric -> 15 points
	// Whole dollar total: 20.00 -> 50 points
	// Multiple of 0.25: 20.00 -> 25 points
	// Items count: 3 items -> (3/2)*5=5 points
	// Items descriptions: "abc"(3 chars, 10.00)-> ceil(10*0.2)=2, "item"(4 chars)->0, "xyz"(3 chars, 6.00)-> ceil(6*0.2)=2 -> total 4
	// Odd day: day 21 odd -> 6 points
	// Purchase time 14:30 between 14:00 and 16:00 -> 10 points
	// Total = 12+50+25+5+4+6+10 = 115

	expected := 115
	got := CalculatePoints(receipt)

	if got != expected {
		t.Errorf("CalculatePoints() = %d when it should be %d", got, expected)
	}
}
