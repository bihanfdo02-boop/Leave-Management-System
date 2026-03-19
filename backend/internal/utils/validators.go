package utils

import (
	"regexp"
	"time"
)

// IsValidEmail checks if email format is valid
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsValidPassword checks if password meets requirements
func IsValidPassword(password string) bool {
	if len(password) < MinPasswordLength {
		return false
	}
	if len(password) > MaxPasswordLength {
		return false
	}
	return true
}

// IsValidDateRange checks if date range is valid
func IsValidDateRange(startDate, endDate time.Time) bool {
	return startDate.Before(endDate) || startDate.Equal(endDate)
}

// IsValidDateString checks if date string is in correct format
func IsValidDateString(dateStr string) bool {
	_, err := time.Parse(DateFormat, dateStr)
	return err == nil
}

// ParseDateString parses a date string
func ParseDateString(dateStr string) (time.Time, error) {
	return time.Parse(DateFormat, dateStr)
}

// CalculateDaysBetween calculates business days between two dates
func CalculateDaysBetween(startDate, endDate time.Time) float32 {
	if startDate.After(endDate) {
		return 0
	}

	duration := endDate.Sub(startDate)
	days := float32(duration.Hours() / 24)
	return days + 1 // Include both start and end dates
}