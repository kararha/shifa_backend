// pkg/utils/helpers.go

package utils

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

// GenerateRandomToken generates a random token of the specified length
func GenerateRandomToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// FormatTime formats a time.Time into a standard string format
func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ParseTime parses a string into a time.Time
func ParseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

// pkg/utils/validators.go

// package utils

// import (
// 	"regexp"
// 	"unicode"
// )

// var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

// // ValidateEmail checks if the provided email is valid
// func ValidateEmail(email string) bool {
// 	return emailRegex.MatchString(email)
// }

// // ValidatePassword checks if the provided password meets the required criteria
// func ValidatePassword(password string) bool {
// 	var (
// 		hasMinLen  = false
// 		hasUpper   = false
// 		hasLower   = false
// 		hasNumber  = false
// 		hasSpecial = false
// 	)
// 	if len(password) >= 8 {
// 		hasMinLen = true
// 	}
// 	for _, char := range password {
// 		switch {
// 		case unicode.IsUpper(char):
// 			hasUpper = true
// 		case unicode.IsLower(char):
// 			hasLower = true
// 		case unicode.IsNumber(char):
// 			hasNumber = true
// 		case unicode.IsPunct(char) || unicode.IsSymbol(char):
// 			hasSpecial = true
// 		}
// 	}
// 	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
// }