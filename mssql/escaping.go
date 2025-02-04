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
	if !isPrintable(name) {
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

func isPrintable(str string) bool {
	for _, charCode := range str {
		if charCode < 32 || charCode >= 127 {
			return false
		}
	}
	return true
}

func escapeStr(str string) string {
	if str == "" {
		return "''"
	}

	panic("TODO")
}
