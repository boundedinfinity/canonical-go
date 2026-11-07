# Date Range

A date entity stores a range of Date entities

## Programmatic Information

Stores a range of [Date](date.md) entities.

Properties: 
    - Start
        type: [Date](date.md)
        desc: The start Date
        validation:
            - call the validation on the Date
    - End
        type: [Date](date.md)
        desc: The end Date
        validation:
            - call the validation on the Date
Operations:
    - Between:
        inputs: 
            - [Date](date.md)
        outputs:
            - boolean
        desc: Returns true if Date is between the start and end properties, inclusive of both, false otherwise.
    - Before:
        inputs: 
            - [Date](date.md)
        outputs:
            - boolean
        desc: Returns true if Date is strictly before the start property, false otherwise.
    - After:
        inputs: 
            - [Date](date.md)
        outputs:
            - boolean
        desc: Returns true if Date is strictly after the end property, false otherwise.
    - Overlaps:
        inputs:
            - [Date Range](date-range.md)
        outputs:
            - boolean
        desc: >
            Returns true if the input shares at least one Date with this range, false otherwise.
            Bounds are inclusive, so ranges that only touch at an endpoint overlap.

## Validation

- All required fields (see above) must be present.
- `Start` must be <= `End`. A range where `Start` equals `End` is a single day.

## Notes

- For any Date, exactly one of `Before`, `Between`, and `After` is true.
