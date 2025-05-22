package points

import (
	"math"
	"regexp"
	"strconv"
	"time"
)

// Used by pointsForRetailerName, returns true if character is alphanumeric.
func isAlphanumeric(r rune) bool {
	return regexp.MustCompile(`[A-Za-z0-9]`).MatchString(string(r))
}

// Used by pointsForWholeDollar to check if total is a whole number, converts the total to an int and compares it to the original float value, returns true if they are equal.
func isWholeDollar(total string) bool {
	if f, err := strconv.ParseFloat(total, 64); err == nil {
		return f == float64(int(f))
	}
	return false
}

// Used by pointsForMultipleOfQuarter , returns true if the total is a multiple of the given factor.
func isMultipleOf(total string, factor float64) bool {
	if f, err := strconv.ParseFloat(total, 64); err == nil {
		return math.Mod(f, factor) == 0
	}
	return false
}

// Used by pointsForOddPurchaseDay, returns true if date is odd.
func isDayOdd(date string) bool {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false
	}
	return t.Day()%2 == 1
}

// Used by pointsForPurchaseTime, returns true if given time string is between the start and end time.
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
