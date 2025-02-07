package mssql

import (
	"testing"
)

func TestMssqlPasswordValidation(t *testing.T) {
	creator := Creator{}

	t.Run("default password is ok", func(t *testing.T) {
		defaultOpts := creator.GetDefaultOpts()
		if err := creator.IsPasswordValid(defaultOpts.Password); err != nil {
			t.Error(err)
		}
	})

	t.Run("empty password", func(t *testing.T) {
		err := creator.IsPasswordValid("")
		if err != ErrPasswordEmpty {
			t.Errorf("Want ErrPasswordEmpty, got: %s", err)
		}
	})

	t.Run("password too short", func(t *testing.T) {
		err := creator.IsPasswordValid("Password1")
		if err != ErrPasswordTooShort {
			t.Errorf("Want ErrPasswordTooShort, got: %s", err)
		}
	})

	complexityCases := [...]struct {
		input  string
		output error
	}{
		{"PASSWORD1!", nil},
		{"password1!", nil},
		{"password12", ErrPasswordTooSimple},
		{"PASSWORD12", ErrPasswordTooSimple},
		{"PASSWORD!!", ErrPasswordTooSimple},
		{"0123456789", ErrPasswordTooSimple},
	}
	for _, tt := range complexityCases {
		t.Run("Password Complexity validation: "+tt.input, func(t *testing.T) {
			err := creator.IsPasswordValid(tt.input)
			if err != tt.output {
				t.Errorf("Want '%s', got: '%s'", tt.output, err)
			}
		})
	}

	complexityOrderCases := [...]string{
		"pAssword12",
		"12Password",
		"Pass12word",
	}
	for _, tt := range complexityOrderCases {
		t.Run("Complexity doesn't depend on order: "+tt, func(t *testing.T) {
			err := creator.IsPasswordValid(tt)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
