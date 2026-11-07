# Time Range

A time entity stores a range of Time entities

## Programmatic Information

Stores a range of [Time](time.md) entities.

Properties: 
    - Start
        type: [Time](time.md)
        desc: The start Time
        validation:
            - call the validation on the Time
    - End
        type: [Time](time.md)
        desc: The end Time
        validation:
            - call the validation on the Time
Operations:
    - Between:
        inputs: 
            - [Time](time.md)
        outputs:
            - boolean
        desc: Returns true if Time is between the start and end properties, inclusive of both, false otherwise.
    - Before:
        inputs: 
            - [Time](time.md)
        outputs:
            - boolean
        desc: Returns true if Time is strictly before the start property, false otherwise.
    - After:
        inputs: 
            - [Time](time.md)
        outputs:
            - boolean
        desc: Returns true if Time is strictly after the end property, false otherwise.
    - Overlaps:
        inputs:
            - [Time Range](time-range.md)
        outputs:
            - boolean
        desc: >
            Returns true if the input shares at least one instant with this range, false otherwise.
            Bounds are inclusive, so ranges that only touch at an endpoint overlap.

## Validation

- All required fields (see above) must be present.
- `Start` must be <= `End`. A range where `Start` equals `End` is a single instant.
- A range does not wrap past midnight. A period that crosses midnight is a [DateTime Range](datetime-range.md).

## Notes

- For any Time, exactly one of `Before`, `Between`, and `After` is true.
- Comparisons place every Time on the same day and apply its offset, then compare by instant. So "18:00:00-05:00" is before "20:00:00-05:00" (even though in UTC they are 23:00 and 01:00), and Times in different offsets compare by their real instant on that day.
- This differs from Time equality, which wraps at midnight, only for Times in different offsets near midnight: "23:00:00-05:00" equals "04:00:00Z" as a Time, but is after it in a range.
- A Time is always valid (midnight is the zero value), so a missing `Start` or `End` cannot be detected by validation; it reads as "00:00:00Z".
