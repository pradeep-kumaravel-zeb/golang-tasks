package utils

import (
	"regexp"
	"strings"
)

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func ValidateName(name string) bool {
	return len(strings.TrimSpace(name)) > 0
}

func ValidateText(text string) bool {
	return len(strings.TrimSpace(text)) > 0
}

func CalculatePaymentStatus(feePaid, courseFee float64) string {
	if feePaid == 0 {
		return "unpaid"
	} else if feePaid < courseFee {
		return "partial"
	}
	return "paid"
}
