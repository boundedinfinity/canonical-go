# Time

A time entity that stores the time of day

## Programmatic Information

Follows the [RFC 3339](https://www.rfc-editor.org/info/rfc3339/) specification for the `full-time` format.

- Time
    type: string
    desc: An RFC 3339 `full-time` string, e.g. "14:30:00Z" or "09:30:00.123-05:00".
    validation:
        - must be a valid RFC 3339 `full-time`, including the time offset (`Z` or `±hh:mm`)
        - fractional seconds are optional

## Notes

- RFC 3339 section 5.6 defines the `full-time` production: https://www.rfc-editor.org/rfc/rfc3339#section-5.6
- A Time has no date.
- Midnight (`00:00:00`) is a valid time, so unlike [DateTime](datetime.md) and [Date](date.md) there is no "must not be zero" rule. The zero value is `00:00:00Z`.
- The time offset is preserved as given; values are not normalized to UTC.
- Two Times are equal when they are the same time of day in UTC, wrapping at midnight (e.g. "23:00:00-05:00" equals "04:00:00Z").
- Serialize with fractional seconds only when present (trailing zeros dropped).
- For go implementation use the standard `time` package, storing the time on the zero date (0001-01-01) so the Go zero value is `00:00:00Z`.
