package points

import (
	"math"
	"strconv"
	"strings"

	"FetchAssessment/models"
)

func CalculatePoints(receipt models.Receipt) int {
	points := 0

	points += pointsForRetailerName(receipt.Retailer)
	points += pointsForWholeDollar(receipt.Total)
	points += pointsForMultipleOfQuarter(receipt.Total)
	points += pointsForItemCount(receipt.Items)
	points += pointsForItemDescriptions(receipt.Items)
	points += pointsForOddPurchaseDay(receipt.PurchaseDate)
	points += pointsForPurchaseTime(receipt.PurchaseTime)

	return points
}

// One point for every alphanumeric character in the retailer name.
func pointsForRetailerName(retailer string) int {
	points := 0
	for _, char := range retailer {
		if isAlphanumeric(char) {
			points++
		}
	}
	return points
}

// 50 points if the total is a round dollar amount with no cents.
func pointsForWholeDollar(total string) int {
	if isWholeDollar(total) {
		return 50
	}
	return 0
}

// 25 points if the total is a multiple of 0.25.
func pointsForMultipleOfQuarter(total string) int {
	if isMultipleOf(total, 0.25) {
		return 25
	}
	return 0
}

// 5 points for every two items on the receipt.
func pointsForItemCount(items []models.Item) int {
	return (len(items) / 2) * 5
}

// If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer. The result is the number of points earned.
func pointsForItemDescriptions(items []models.Item) int {
	points := 0
	for _, item := range items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			if price, err := strconv.ParseFloat(item.Price, 64); err == nil {
				points += int(math.Ceil(price * 0.2))
			}
		}
	}
	return points
}

// 6 points if the day in the purchase date is odd.
func pointsForOddPurchaseDay(date string) int {
	if isDayOdd(date) {
		return 6
	}
	return 0
}

// 10 points if the time of purchase is between 2:00pm and 4:00pm.
func pointsForPurchaseTime(timeStr string) int {
	if isBetweenTime(timeStr, "14:00", "16:00") {
		return 10
	}
	return 0
}
