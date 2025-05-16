package points

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"FetchAssessment/models"
)

func CalculatePoints(receipt models.Receipt) int {
	points := 0

	// One point for every alphanumeric character in the retailer name.
	for _, char := range receipt.Retailer {
		if isAlphanumeric(char) {
			points++
		}
	}

	// 50 points if the total is a round dollar amount with no cents.
	if isWholeDollar(receipt.Total) {
		points += 50
	}

	// 25 points if the total is a multiple of 0.25.
	if isMultipleOf(receipt.Total, 0.25) {
		points += 25
	}

	// 5 points for every two items on the receipt.
	points += (len(receipt.Items) / 2) * 5

	// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
	for _, item := range receipt.Items {
		desc := strings.TrimSpace(item.ShortDescription)
		if len(desc)%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err == nil {
				itemPoints := int(math.Ceil(price * 0.2))
				points += itemPoints
			}
		}
	}

	// 6 points if the day in the purchase date is odd.
	if isDayOdd(receipt.PurchaseDate) {
		points += 6
	}

	// 10 points if the time of purchase is between 2:00pm and 4:00pm.
	if isBetweenTime(receipt.PurchaseTime, "14:00", "16:00") {
		points += 10
	}

	return points
}

func isAlphanumeric(r rune) bool {
	return regexp.MustCompile(`[A-Za-z0-9]`).MatchString(string(r))
}

func isWholeDollar(total string) bool {
	if f, err := strconv.ParseFloat(total, 64); err == nil {
		return f == float64(int(f))
	}
	return false
}

func isMultipleOf(total string, factor float64) bool {
	if f, err := strconv.ParseFloat(total, 64); err == nil {
		return math.Mod(f, factor) == 0
	}
	return false
}

func isDayOdd(date string) bool {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false
	}
	return t.Day()%2 == 1
}

func isBetweenTime(tStr, startStr, endStr string) bool {
	layout := "15:04"
	t, err1 := time.Parse(layout, tStr)
	start, err2 := time.Parse(layout, startStr)
	end, err3 := time.Parse(layout, endStr)

	if err1 != nil || err2 != nil || err3 != nil {
		return false
	}
	return t.After(start) && t.Before(end)
}
