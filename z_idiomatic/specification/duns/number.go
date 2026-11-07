package duns

import "github.com/boundedinfinity/go-commoner/errorer"

var (
	ErrNumber   = errorer.New("DUNS number err")
	errNumberfn = errorer.Func(ErrNumber)
)

type Number string

func (this Number) Validate() error {
	if len(this) != 9 {
		return errNumberfn("%s : length is not 9", this)
	}

	for i := 0; i < len(this); i++ {
		if this[i] < '0' || this[i] > '9' {
			return errNumberfn("%s : contains non-digit characters", this)
		}
	}

	return nil
}
