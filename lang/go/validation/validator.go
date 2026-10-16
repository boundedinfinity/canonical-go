package validation

type Validator interface {
	Validate(a any) error
}

type ValidationManager struct {
}
