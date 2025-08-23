package mssql

import (
	"strings"
	"unicode"
)

const specialCharClass = "!@#$%^&*()_-+={}[]\\|/<>~,.;:'\""

func isLatinLower(c rune) bool {
	return 97 <= c && c <= 122
}

func isLatinUpper(c rune) bool {
	return 65 <= c && c <= 90
}

// https://learn.microsoft.com/en-us/sql/relational-databases/security/password-policy?view=sql-server-ver16#password-complexity
func isPasswordComplexEnough(password string) bool {
	var hasLower, hasUpper, hasDigit, hasSpecial bool
	var numberOfCharClassesMatched int

	for _, c := range password {
		if !hasLower && isLatinLower(c) {
			hasLower = true
			numberOfCharClassesMatched++
		} else if !hasUpper && isLatinUpper(c) {
			hasUpper = true
			numberOfCharClassesMatched++
		} else if !hasDigit && unicode.IsDigit(c) {
			hasDigit = true
			numberOfCharClassesMatched++
		} else if !hasSpecial && strings.ContainsRune(specialCharClass, c) {
			hasSpecial = true
			numberOfCharClassesMatched++
		}
		if numberOfCharClassesMatched >= 3 {
			return true
		}
	}
	return false
}
