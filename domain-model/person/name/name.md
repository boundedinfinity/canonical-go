# Name

**Top Level Model**

The name of a person (can be real, fictional or otherwise)

## Programmatic Information

- ID
    type: [Id](id.md)
    desc: The unique identifier of the label.
    validation:
        - call the validation on the ID
- First
    description: The first names
    type: List strings
    optional: true
    validation:
        - if present, at least 1 element
        - each string length >= 1 and <= 30 
- Middle
    description: The middle names
    type: List strings
    optional: true
    validation:
        - if present, at least 1 element
        - each string length >= 1 and <= 30 
- Last
    description: The last names
    type: List strings
    optional: true
    validation:
        - if present, at least 1 element
        - each string length >= 1 and <= 30 

# Validation

- At least one of `First`, `Middle`, or `Last` must be present — a `Personal Name` with none of them set is meaningless, even though each is individually optional.
- `First`, `Middle`, and `Last`, if present, must each contain at least one element, and each name string must have length >= 1 and <= 30.
