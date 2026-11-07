# DateTime

A datetime entity that stores the date and time

## Programmatic Information

Follows the [RFC 3339](https://www.rfc-editor.org/info/rfc3339/) specification for the `date-time` format.

- DateTime
    type: string
    desc: An RFC 3339 `date-time` string, e.g. "2026-11-07T14:30:00Z" or "2026-11-07T09:30:00.123-05:00".
    validation:
        - must be a valid RFC 3339 `date-time`, including the time offset (`Z` or `±hh:mm`)
        - fractional seconds are optional
        - must not be the zero value

## Notes

- RFC 3339 section 5.6 defines the `date-time` production: https://www.rfc-editor.org/rfc/rfc3339#section-5.6
- The time offset is preserved as given; values are not normalized to UTC.
- Serialize with fractional seconds only when present (trailing zeros dropped).
- For go implementation use the standard `time` package (`time.RFC3339Nano`).
