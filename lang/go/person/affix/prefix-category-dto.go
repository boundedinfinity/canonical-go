package affix

// Source of truth: domain-model/person/affix/prefix-category.md,
// domain-model/global-notes/dto.md

// PrefixCategoryDto is the data transfer object for PrefixCategory.
// PrefixCategory holds no other Top Level Models, so PrefixCategoryDto has
// the same fields, validation, and JSON encoding, all from the embedded
// Category.
type PrefixCategoryDto struct {
	Category
}

// ToModel returns the validated PrefixCategory for this PrefixCategoryDto.
func (this PrefixCategoryDto) ToModel() (PrefixCategory, error) {
	category := PrefixCategory(this)

	if err := category.Validate(); err != nil {
		return PrefixCategory{}, err
	}

	return category, nil
}

// ToDto returns the PrefixCategoryDto for this PrefixCategory.
func (this PrefixCategory) ToDto() PrefixCategoryDto {
	return PrefixCategoryDto(this)
}
