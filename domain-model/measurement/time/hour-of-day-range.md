# Hour Of Day Range

An hour of day entity stores a range of Hour Of Day entities

## Programmatic Information

Stores a range of [Hour Of Day](hour-of-day.md) entities.

Properties: 
    - Start
        type: [Hour Of Day](hour-of-day.md)
        desc: The start Hour Of Day
        validation:
            - call the validation on the Hour Of Day
    - End
        type: [Hour Of Day](hour-of-day.md)
        desc: The end Hour Of Day
        validation:
            - call the validation on the Hour Of Day
Operations:
    - Between:
        inputs: 
            - [Hour Of Day](hour-of-day.md)
        outputs:
            - boolean
        desc: Returns true if Hour Of Day is between the start and end properties, inclusive of both, false otherwise.
    - Before:
        inputs: 
            - [Hour Of Day](hour-of-day.md)
        outputs:
            - boolean
        desc: Returns true if Hour Of Day is strictly before the start property, false otherwise.
    - After:
        inputs: 
            - [Hour Of Day](hour-of-day.md)
        outputs:
            - boolean
        desc: Returns true if Hour Of Day is strictly after the end property, false otherwise.
    - Overlaps:
        inputs:
            - [Hour Of Day Range](hour-of-day-range.md)
        outputs:
            - boolean
        desc: >
            Returns true if the input shares at least one Hour Of Day with this range, false otherwise.
            Bounds are inclusive, so ranges that only touch at an endpoint overlap.

## Validation

- All required fields (see above) must be present.
- `Start` must be <= `End`. A range where `Start` equals `End` is a single hour.

## Notes

- For any Hour Of Day, exactly one of `Before`, `Between`, and `After` is true.
- Both ends are inclusive: Start 9 and End 17 covers the hours 9 through 17, i.e. 9:00 to 17:59.
- A range does not wrap past midnight: Start 22 and End 2 is invalid. Use a [DateTime Range](datetime-range.md) for a period that crosses midnight.
- An Hour Of Day is always valid when unset (midnight is the zero value), so a missing `Start` or `End` cannot be detected by validation; it reads as 0.
- The data was based on the legacy `HourOfDayRange` in `z_idiomatic/measurement/time/hour-of-day.go`, whose `Min` and `Max` fields are renamed `Start` and `End` to match the other ranges.
