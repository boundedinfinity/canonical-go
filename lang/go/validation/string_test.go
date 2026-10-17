package validation_test

import (
	"regexp"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/validation"
	"github.com/stretchr/testify/assert"
)

func Test_String_Match(t *testing.T) {
	v := validation.String[string]("code").Match(regexp.MustCompile(`^[A-Z]{2}$`))

	tcs := []struct {
		input string
		ok    bool
	}{
		{input: "US", ok: true},
		{input: "us"},
		{input: "USA"},
		{input: ""},
	}

	for _, tc := range tcs {
		t.Run(tc.input, func(tt *testing.T) {
			err := v.Validate(tc.input)

			if tc.ok {
				assert.NoError(tt, err)
				return
			}

			assert.ErrorIs(tt, err, validation.StringValidationError)
		})
	}
}

func Test_String_URL(t *testing.T) {
	v := validation.String[string]("reference").URL("http", "https")

	tcs := []struct {
		input string
		ok    bool
	}{
		{input: "https://en.wikipedia.org/wiki/Alabama", ok: true},
		{input: "http://example.com", ok: true},
		{input: "ftp://example.com"},
		{input: "/wiki/Alabama"},
		{input: "https://"},
		{input: "not a url"},
		{input: "https://exa mple.com"},
		{input: ""},
	}

	for _, tc := range tcs {
		t.Run(tc.input, func(tt *testing.T) {
			err := v.Validate(tc.input)

			if tc.ok {
				assert.NoError(tt, err)
				return
			}

			assert.ErrorIs(tt, err, validation.StringValidationError)
		})
	}
}
