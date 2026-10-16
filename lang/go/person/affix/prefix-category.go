package affix

// Source of truth: domain-model/person/affix/prefix-category.md

// PrefixCategory is a category of Prefix, such as "Honorific". It has the fields,
// validation, and JSON encoding of Category, which it embeds.
type PrefixCategory struct {
	Category
}

// NewPrefixCategory returns a validated PrefixCategory with a new ID.
// An empty description is treated as not present.
func NewPrefixCategory(name, description string) (PrefixCategory, error) {
	category, err := NewCategory(name, description)

	if err != nil {
		return PrefixCategory{}, err
	}

	return PrefixCategory{Category: category}, nil
}

// MustNewPrefixCategory is like NewPrefixCategory but panics if the
// PrefixCategory is invalid.
func MustNewPrefixCategory(name, description string) PrefixCategory {
	category, err := NewPrefixCategory(name, description)

	if err != nil {
		panic(err)
	}

	return category
}
