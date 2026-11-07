# Duration

A duration entity that stores a length of time as calendar and clock components

## Programmatic Information

Follows the [RFC 3339](https://www.rfc-editor.org/info/rfc3339/) specification for the `duration` format (Appendix A).

- Duration
    type: string
    desc: An RFC 3339 `duration` string, e.g. "P1Y2M3DT4H5M6S", "PT30M", or "P2W".
    validation:
        - must be a valid RFC 3339 `duration`

Components:
    - Years
        type: integer
        validation:
            - must be >= 0
    - Months
        type: integer
        validation:
            - must be >= 0
    - Weeks
        type: integer
        validation:
            - must be >= 0
    - Days
        type: integer
        validation:
            - must be >= 0
    - Hours
        type: integer
        validation:
            - must be >= 0
    - Minutes
        type: integer
        validation:
            - must be >= 0
    - Seconds
        type: integer
        validation:
            - must be >= 0

Operations:
    - Time:
        inputs: (none)
        outputs:
            - the native language duration
            - error
        desc: >
            Returns the Duration as the language's fixed-length duration type.
            Weeks convert as 7 days of 24 hours; Hours, Minutes, and Seconds convert exactly.
            Returns an error if the Duration is invalid, if any of Years, Months, or Days is
            non-zero (they have no fixed length), or if the result does not fit the native type.

## Validation

- `Weeks` cannot be combined with any other component.

## Notes

- RFC 3339 Appendix A defines the `duration` production: https://www.rfc-editor.org/rfc/rfc3339#appendix-A
- The grammar is followed strictly:
    - Values are whole, non-negative numbers; there are no fractions or signs.
    - At least one component is required ("P" and "PT" are invalid).
    - `T` separates the date components (Y, M, D) from the time components (H, M, S) and must be followed by at least one time component.
    - Within the date part, and within the time part, components appear in order and cannot skip a unit: "P1Y1D" and "PT1H1S" are invalid, written instead as "P1Y0M1D" and "PT1H0M1S". Skipping across `T` is allowed ("P1YT1S").
    - Weeks stand alone ("P2W"); "P1Y2W" and "P2WT1H" are invalid.
- Components are stored as written and not converted to one another: years, months, and days have no fixed length (months are 28–31 days, and a day can be 23 or 25 hours across a DST change), so "P1M" and "P30D" are different durations.
- Serialize with the fewest components the grammar allows, filling skipped units with zeros (e.g. Years 1 and Days 1 serialize as "P1Y0M1D"). A zero duration serializes as "PT0S".
- A zero duration is valid.
- Days are not converted by `Time` even though they are often treated as 24 hours: across a DST change a day is 23 or 25 hours. Use "PT24H" for an exact 24 hours.
- For go implementation the native duration is `time.Duration`, whose maximum is about 292 years.
