package mssql

import (
	"strings"
	"testing"
)

func Test_escapeId(t *testing.T) {
	bracketsValues := [...]struct {
		input  string
		output string
	}{
		{input: "foo", output: "[foo]"},
		{input: "_foo", output: "[_foo]"},
		{input: "foo.", output: "[foo.]"},
	}

	for _, tt := range bracketsValues {
		t.Run("always wrap identifier in square brackets: "+tt.input+" => "+tt.output, func(t *testing.T) {
			got, err := escapeId(tt.input)
			if err != nil {
				t.Error(err)
			}
			if got != tt.output {
				t.Errorf("Unexpected value, want %s, got %s", tt.output, got)
			}
		})
	}

	t.Run("treats characters from range [_$#@.\\-] as valid", func(t *testing.T) {
		got, err := escapeId("_$#.@-")
		if err != nil {
			t.Error(err)
		}
		if want := "[_$#.@-]"; want != got {
			t.Errorf("Unexpected value, want %s, got %s", want, got)
		}
	})

	t.Run("throws on empty strings", func(t *testing.T) {
		got, err := escapeId("")
		if got != "" || err == nil {
			t.Error("expected escapeId to throw, but it didn't")
		}
	})

	t.Run("throw on non-printable characters in identifiers", func(t *testing.T) {
		_, err := escapeId("f\noo")
		if err == nil {
			t.Error("expected escapeId to throw, but it didn't")
		}
		_, err = escapeId("f\boo")
		if err == nil {
			t.Error("expected escapeId to throw, but it didn't")
		}
	})

}

func Test_escapeUser(t *testing.T) {
	t.Run("throws on empty strings", func(t *testing.T) {
		val, err := escapeUser("")
		if val != "" || err == nil {
			t.Error("expected escapeUser to throw, but it didn't")
		}

		val, err = escapeUser("a")
		if val == "" || err != nil {
			t.Error(err)
		}
	})

	t.Run("throws on non-printable characters", func(t *testing.T) {
		val, err := escapeUser("as\nda")
		if val != "" || err == nil {
			t.Error("expected escapeUser to throw, but it didn't")
		}

		val, err = escapeUser("asda")
		if val == "" || err != nil {
			t.Error(err)
		}
	})

	t.Run("throws if user name gte 128 characters long", func(t *testing.T) {
		val, err := escapeUser(strings.Repeat("a", 128))
		if val != "" || err == nil {
			t.Error("expected escapeUser to throw, but it didn't")
		}
		val, err = escapeUser(strings.Repeat("a", 127))
		if val == "" || err != nil {
			t.Error(err)
		}
	})

	t.Run("escapes user name in the same way as identifier", func(t *testing.T) {
		values := [...]string{
			"!asd",
			"_$#.@-",
			"as.d",
		}
		for _, value := range values {
			t.Run("value: "+value, func(t *testing.T) {
				user, err := escapeUser(value)
				if err != nil {
					t.Error(err)
				}
				id, err := escapeId(value)
				if err != nil {
					t.Error(err)
				}
				if user != id {
					t.Errorf("Values do not match, want '%s', got '%s'", id, user)
				}
			})
		}
	})
}

func Test_escapeStr(t *testing.T) {
	t.Run("always wraps a string in single quots", func(t *testing.T) {
		got := escapeStr("foo")
		if want := "'foo'"; want != got {
			t.Errorf("Values do not match, want '%s', got '%s'", want, got)
		}
	})

	singleQuotesValues := [...]struct {
		input  string
		output string
	}{
		{input: "o'foo", output: "'o''foo'"},
		{input: "'ofoo", output: "'''ofoo'"},
	}
	for _, tt := range singleQuotesValues {
		t.Run("turns single quots in strin into a pair of singlequots", func(t *testing.T) {
			got := escapeStr(tt.input)
			if want := tt.output; want != got {
				t.Errorf("Values do not match, want '%s', got '%s'", want, got)
			}
		})
	}

	t.Run("turns non-ascii chars into a CHAR() concatenation", func(t *testing.T) {
		got := escapeStr("foo\r\n")
		if want := "'foo' + CHAR(13) + CHAR(10)"; want != got {
			t.Errorf("Values do not match, want '%s', got '%s'", want, got)
		}
	})

	t.Run("still escapes concatenated pairs", func(t *testing.T) {
		got := escapeStr("foo\nb'ar")
		if want := "'foo' + CHAR(10) + 'b''ar'"; want != got {
			t.Errorf("Values do not match, want '%s', got '%s'", want, got)
		}
	})

	t.Run("doesn't append extra chars at start", func(t *testing.T) {
		got := escapeStr("\nfoo")
		if want := "CHAR(10) + 'foo'"; want != got {
			t.Errorf("Values do not match, want '%s', got '%s'", want, got)
		}
	})

	t.Run("encodes empty strings", func(t *testing.T) {
		got := escapeStr("")
		if want := "''"; want != got {
			t.Errorf("Values do not match, want '%s', got '%s'", want, got)
		}
	})
}
