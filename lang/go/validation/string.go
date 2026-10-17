package validation

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

var (
	StringValidationError = errors.New("string validation")
)

func String[T ~string](name string) *stringValidator[T] {
	return &stringValidator[T]{
		name: name,
	}
}

type stringValidator[T ~string] struct {
	name        string
	validations []func(T) error
}

func (this stringValidator[T]) Validate(a any) error {
	t, ok := a.(T)

	if !ok {
		var t T
		panic(fmt.Errorf("%T is not of type %T", a, t))
	}

	var errs []error

	for _, fn := range this.validations {
		if err := fn(t); err != nil {
			if this.name != "" {
				err = fmt.Errorf("%s %w", this.name, err)
			}
			errs = append(errs, fmt.Errorf("%w : %w", StringValidationError, err))
		}
	}

	return errors.Join(errs...)
}

func (this *stringValidator[T]) Name(v string) *stringValidator[T] {
	this.name = v
	return this
}

func (this *stringValidator[T]) Add(fn func(T) error) *stringValidator[T] {
	this.validations = append(this.validations, fn)
	return this
}

func (this *stringValidator[T]) Min(t int) *stringValidator[T] {
	validation := func(v T) error {
		l := utf8.RuneCountInString(string(v))

		if l < t {
			return fmt.Errorf("length %d less than minimium %d", l, t)
		}

		return nil
	}

	return this.Add(validation)
}

func (this *stringValidator[T]) Max(t int) *stringValidator[T] {
	validation := func(v T) error {
		l := utf8.RuneCountInString(string(v))

		if l > t {
			return fmt.Errorf("length %d greater than maximum %d", l, t)
		}

		return nil
	}

	return this.Add(validation)
}

func (this *stringValidator[T]) OneOf(ts ...string) *stringValidator[T] {
	validation := func(v T) error {
		s := string(v)

		if !slices.Contains(ts, s) {
			return fmt.Errorf("%s is not one of %s", s, strings.Join(ts, ", "))
		}

		return nil
	}

	return this.Add(validation)
}

func (this *stringValidator[T]) NoneOf(ts ...string) *stringValidator[T] {
	validation := func(v T) error {
		s := string(v)

		if slices.Contains(ts, s) {
			return fmt.Errorf("%s is one of %s", s, strings.Join(ts, ", "))
		}

		return nil
	}

	return this.Add(validation)
}

// Match requires the string to match re. Anchor the pattern with ^ and $ to
// require a whole-string match.
func (this *stringValidator[T]) Match(re *regexp.Regexp) *stringValidator[T] {
	validation := func(v T) error {
		s := string(v)

		if !re.MatchString(s) {
			return fmt.Errorf("%q does not match %s", s, re)
		}

		return nil
	}

	return this.Add(validation)
}

// URL requires an absolute URL with a host and one of the given schemes.
func (this *stringValidator[T]) URL(schemes ...string) *stringValidator[T] {
	validation := func(v T) error {
		s := string(v)
		u, err := url.Parse(s)

		if err != nil {
			return fmt.Errorf("%q is not a URL: %w", s, err)
		}

		if !u.IsAbs() || u.Host == "" {
			return fmt.Errorf("%q is not an absolute URL", s)
		}

		if !slices.Contains(schemes, u.Scheme) {
			return fmt.Errorf("%q scheme is not one of %s", s, strings.Join(schemes, ", "))
		}

		return nil
	}

	return this.Add(validation)
}
