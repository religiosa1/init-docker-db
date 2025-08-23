package mssql

import (
	"testing"
)

func TestValidatePassword(t *testing.T) {
	creator := Creator{}

	t.Run("default password is ok", func(t *testing.T) {
		defaultOpts := creator.GetDefaultOpts()
		if err := creator.ValidatePassword(defaultOpts.Password); err != nil {
			t.Error(err)
		}
	})

	t.Run("empty password", func(t *testing.T) {
		err := creator.ValidatePassword("")
		if err != ErrPasswordEmpty {
			t.Errorf("Want ErrPasswordEmpty, got: %s", err)
		}
	})

	t.Run("password too short", func(t *testing.T) {
		err := creator.ValidatePassword("Password1")
		if err != ErrPasswordTooShort {
			t.Errorf("Want ErrPasswordTooShort, got: %s", err)
		}
	})

	t.Run("Password Complexity validation", func(t *testing.T) {
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
			t.Run(tt.input, func(t *testing.T) {
				err := creator.ValidatePassword(tt.input)
				if err != tt.output {
					t.Errorf("Want '%s', got: '%s'", tt.output, err)
				}
			})
		}
	})

	t.Run("Complexity validation doesn't depend on char classes order", func(t *testing.T) {
		complexityOrderCases := [...]string{
			"pAssword12",
			"12Password",
			"Pass12word",
		}
		for _, tt := range complexityOrderCases {
			t.Run(tt, func(t *testing.T) {
				err := creator.ValidatePassword(tt)
				if err != nil {
					t.Error(err)
				}
			})
		}
	})
}
