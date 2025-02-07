package mssql

import (
	"errors"
	"fmt"
	"strings"
)

func escapeId(name string) (string, error) {
	if name == "" {
		return "", errors.New("mssql identifier cannot be empty")
	}
	if strings.ContainsRune(name, '[') || strings.ContainsRune(name, ']') {
		return "", errors.New("mssql identifier cannot contain '[' or ']' characters")
	}
	if !isStringPrintable(name) {
		return "", errors.New("mssql identifiers cannot contain non-printable characters")
	}
	return fmt.Sprintf("[%s]", name), nil
}

func escapeUser(name string) (string, error) {
	if len(name) >= 128 {
		return "", errors.New("user name cannot be longer than 128 charaters")
	}
	return escapeId(name)
}

func isRunePrintable(r rune) bool {
	return r >= 32 && r < 127
}

func isStringPrintable(str string) bool {
	for _, c := range str {
		if !isRunePrintable(c) {
			return false
		}
	}
	return true
}

func escapeStr(str string) string {
	if str == "" {
		return "''"
	}

	var sb strings.Builder
	tokens := tokenizeString(str)
	for i, token := range tokens {
		if i != 0 {
			sb.WriteString(" + ")
		}
		if token.printable {
			sb.WriteString("'")
			sb.WriteString(strings.ReplaceAll(token.value, "'", "''"))
			sb.WriteString("'")
		} else {
			sb.WriteString(fmt.Sprintf("CHAR(%d)", rune(token.value[0])))
		}
	}
	return sb.String()
}

type Token struct {
	value     string
	printable bool
}

func tokenizeString(str string) []Token {
	var result []Token
	lastNonPrintableIndex := -1
	runes := []rune(str)

	for i, r := range runes {
		if isRunePrintable(r) {
			continue
		}
		if i != 0 && lastNonPrintableIndex != i-1 {
			result = append(result, Token{string(runes[lastNonPrintableIndex+1 : i]), true})
		}
		result = append(result, Token{string(r), false})
		lastNonPrintableIndex = i
	}

	if lastNonPrintableIndex != len(runes)-1 {
		result = append(result, Token{string(runes[lastNonPrintableIndex+1:]), true})
	}

	return result
}
