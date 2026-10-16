package affix

// Source of truth: domain-model/person/affix/suffix-category.md

// SuffixCategory is a category of Suffix, such as "Generational". It has the fields,
// validation, and JSON encoding of Category, which it embeds.
type SuffixCategory struct {
	Category
}

// NewSuffixCategory returns a validated SuffixCategory with a new ID.
// An empty description is treated as not present.
func NewSuffixCategory(name, description string) (SuffixCategory, error) {
	category, err := NewCategory(name, description)

	if err != nil {
		return SuffixCategory{}, err
	}

	return SuffixCategory{Category: category}, nil
}

// MustNewSuffixCategory is like NewSuffixCategory but panics if the
// SuffixCategory is invalid.
func MustNewSuffixCategory(name, description string) SuffixCategory {
	category, err := NewSuffixCategory(name, description)

	if err != nil {
		panic(err)
	}

	return category
}
