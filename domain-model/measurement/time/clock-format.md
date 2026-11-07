# Clock Format

How times of day are displayed: on a 12-hour clock (2 PM) or a 24-hour clock (14:00)

## Programmatic Information

- Clock Format
    type: string
    desc: The name of the format, e.g. "24-hour".
    validation:
        - must be one of the names below, ignoring case

Values:

| Name    | Example  |
| ------- | -------- |
| 12-hour | 2:30 PM  |
| 24-hour | 14:30    |

## Notes

- Serialize as the name (e.g. "24-hour"). Parsing also accepts any letter case ("24-HOUR"), but nothing else ("24", "24h", "military").
- There is no zero format: an unset Clock Format fails validation.
- A Clock Format only says how to display a time. Times are always stored on a 24-hour clock.
- The language's string representation (e.g. Go's `String()`) returns the Name, e.g. "24-hour".
- The legacy implementation in `z_idiomatic/measurement/time/clock-format.go` used the names "measurement.time.clock-format.12-hour" and "measurement.time.clock-format.24-hour".
- For go implementation store the format in an integer type.
