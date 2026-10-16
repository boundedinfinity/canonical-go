package affix

// Source of truth: domain-model/person/affix/suffix-category.md,
// domain-model/global-notes/dto.md

// SuffixCategoryDto is the data transfer object for SuffixCategory.
// SuffixCategory holds no other Top Level Models, so SuffixCategoryDto has
// the same fields, validation, and JSON encoding, all from the embedded
// Category.
type SuffixCategoryDto struct {
	Category
}

// ToModel returns the validated SuffixCategory for this SuffixCategoryDto.
func (this SuffixCategoryDto) ToModel() (SuffixCategory, error) {
	category := SuffixCategory(this)

	if err := category.Validate(); err != nil {
		return SuffixCategory{}, err
	}

	return category, nil
}

// ToDto returns the SuffixCategoryDto for this SuffixCategory.
func (this SuffixCategory) ToDto() SuffixCategoryDto {
	return SuffixCategoryDto(this)
}
