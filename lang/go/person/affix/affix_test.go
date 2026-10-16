package affix_test

import (
	"strings"
	"testing"

	"github.com/boundedinfinity/canonical-go/lang/go/id"
	"github.com/boundedinfinity/canonical-go/lang/go/person/affix"
	"github.com/boundedinfinity/go-commoner/functional/optioner"
	"github.com/stretchr/testify/assert"
)

func Test_Affix_Validate(t *testing.T) {
	validId := id.MustParse("b12e7628-35b0-4ae5-b2d9-0f34f4828229")

	tcs := []struct {
		name  string
		input affix.Affix
		errs  []error
	}{
		{name: "name only", input: affix.Affix{Id: validId, Name: "Dr."}},
		{name: "all fields", input: affix.Affix{Id: validId, Name: "Dr.", Description: optioner.Some("Doctor")}},
		{name: "name max", input: affix.Affix{Id: validId, Name: strings.Repeat("a", 30)}},
		{name: "name max counts runes", input: affix.Affix{Id: validId, Name: strings.Repeat("é", 30)}},
		{name: "description min", input: affix.Affix{Id: validId, Name: "a", Description: optioner.Some("ab")}},
		{name: "description max", input: affix.Affix{Id: validId, Name: "a", Description: optioner.Some(strings.Repeat("a", 500))}},
		{name: "description max counts runes", input: affix.Affix{Id: validId, Name: "a", Description: optioner.Some(strings.Repeat("é", 500))}},
		{name: "zero id", input: affix.Affix{Name: "a"}, errs: []error{affix.ErrId, id.ErrInvalid}},
		{name: "name empty", input: affix.Affix{Id: validId}, errs: []error{affix.ErrName}},
		{name: "name too long", input: affix.Affix{Id: validId, Name: strings.Repeat("a", 31)}, errs: []error{affix.ErrName}},
		{name: "description empty", input: affix.Affix{Id: validId, Name: "a", Description: optioner.Some("")}, errs: []error{affix.ErrDescription}},
		{name: "description too short", input: affix.Affix{Id: validId, Name: "a", Description: optioner.Some("a")}, errs: []error{affix.ErrDescription}},
		{name: "description too long", input: affix.Affix{Id: validId, Name: "a", Description: optioner.Some(strings.Repeat("a", 501))}, errs: []error{affix.ErrDescription}},
		{
			name:  "reports every failure",
			input: affix.Affix{Description: optioner.Some("a")},
			errs:  []error{affix.ErrId, affix.ErrName, affix.ErrDescription},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(tt *testing.T) {
			err := tc.input.Validate()

			if len(tc.errs) > 0 {
				assert.ErrorIs(tt, err, affix.ErrInvalid)
				for _, e := range tc.errs {
					assert.ErrorIs(tt, err, e)
				}
				return
			}

			assert.NoError(tt, err)
		})
	}
}
