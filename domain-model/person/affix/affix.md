
# Affix

The "base" struct for prefixes and suffixes used with personal names.

## Programmatic Information

- Id
    type: [Id](../../id/id.md)
    desc: The unique identifier of the affix.
    validation:
        - call the validation on the [ID](../../id/id.md)
- Name
    type: string
    desc: The name of the affix.
    validation:
        - length >= 1 and <= 30
- Description
    type: string
    desc: The description of the affix.
    optional: true
    validation:
        - length >= 2 and <= 500
- Category
    type: [Prefix Category](prefix-category.md) for a [Prefix](prefix.md), or [Suffix Category](suffix-category.md) for a [Suffix](suffix.md)
    desc: The category of the affix.
    validation:
        - call the validation on the category, which is either a [Prefix Category](prefix-category.md) or a [Suffix Category](suffix-category.md) (see [Category](category.md))

## Validation

- All required fields (see above) must be present.
- `ID` must be a valid [Id](../../id/id.md).
- `Name` length must be >= 1 and <= 30.
- `Description`, if present, length must be >= 2 and <= 500.
- `Category` must be present and must be a valid [Category](category.md).
- A [Prefix](prefix.md) must use a [Prefix Category](prefix-category.md), and a [Suffix](suffix.md) must use a [Suffix Category](suffix-category.md).
