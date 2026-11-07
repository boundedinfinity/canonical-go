# Quarter Range

A quarter entity stores a range of Quarter entities

## Programmatic Information

Stores a range of [Quarter](quarter.md) entities.

Properties: 
    - Start
        type: [Quarter](quarter.md)
        desc: The start Quarter
        validation:
            - call the validation on the Quarter
    - End
        type: [Quarter](quarter.md)
        desc: The end Quarter
        validation:
            - call the validation on the Quarter
Operations:
    - Between:
        inputs: 
            - [Quarter](quarter.md)
        outputs:
            - boolean
        desc: Returns true if Quarter is between the start and end properties, inclusive of both, false otherwise.
    - Before:
        inputs: 
            - [Quarter](quarter.md)
        outputs:
            - boolean
        desc: Returns true if Quarter is strictly before the start property, false otherwise.
    - After:
        inputs: 
            - [Quarter](quarter.md)
        outputs:
            - boolean
        desc: Returns true if Quarter is strictly after the end property, false otherwise.
    - Overlaps:
        inputs:
            - [Quarter Range](quarter-range.md)
        outputs:
            - boolean
        desc: >
            Returns true if the input shares at least one Quarter with this range, false otherwise.
            Bounds are inclusive, so ranges that only touch at an endpoint overlap.

## Validation

- All required fields (see above) must be present.
- `Start` must be <= `End`. A range where `Start` equals `End` is a single quarter.

## Notes

- For any Quarter, exactly one of `Before`, `Between`, and `After` is true.
- A range does not wrap past Q4: "Q4 to Q1" is invalid. A period that crosses into the next year is a [Date Range](date-range.md).
