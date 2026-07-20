package model

var _ Model = &EnumerationModel{}

type EnumerationModel struct {
	Name            string   `json:"name"`
	Values          []string `json:"values"`
	CaseInsensitive bool     `json:"case-insensitive"`
}

func (_ EnumerationModel) Kind() Kind {
	return Kinds.Enumeration
}

func (this EnumerationModel) Validate() error {
	return nil
}
