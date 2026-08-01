package sic

import "errors"

var ErrInvalidCode = errors.New("invalid SEC SIC code")

type Code struct {
	Code          string `json:"code"`
	Office        string `json:"office"`
	IndustryTitle string `json:"industry-title"`
}

func (code Code) Validate() error {
	if _, ok := validCodes[code]; !ok {
		return ErrInvalidCode
	}

	return nil
}
