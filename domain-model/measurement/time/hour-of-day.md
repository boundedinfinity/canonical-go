# Hour Of Day

An hour of the day on a 24-hour clock, with no date, minutes, or time offset

## Programmatic Information

Follows [ISO 8601](https://www.iso.org/iso-8601-date-and-time-format.html) 24-hour clock hours.

- Hour Of Day
    type: integer
    desc: The hour, from 0 (midnight to 1 AM) through 23 (11 PM to midnight), e.g. 14 for 2 PM.
    validation:
        - must be >= 0 and <= 23

Operations:
    - Number:
        inputs: (none)
        outputs:
            - integer
        desc: Returns the hour, 0 through 23.

## Notes

- Serialize as a JSON number (e.g. `14`). Its text form, used for example as a map key, is the decimal hour without padding (e.g. `"7"`); parsing text also accepts one leading zero (`"07"`), but nothing else (`"7 AM"`, `"7:00"`, `" 7"`, `"24"`).
- Midnight (0) is a valid hour, so, like [Time](time.md), there is no "must not be zero" rule. The zero value is midnight.
- An Hour Of Day can be built from a native date or date-time; the hour is taken in that value's own time offset.
- The language's string representation (e.g. Go's `String()`) is the decimal hour, e.g. "14".
- The data was based on the legacy implementation in `z_idiomatic/measurement/time/hour-of-day.go`.
