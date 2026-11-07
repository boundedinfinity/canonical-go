# Date

A date entity that stores the date

## Programmatic Information

Follows the [RFC 3339](https://www.rfc-editor.org/info/rfc3339/) specification for the `full-date` format.

- Date
    type: string
    desc: An RFC 3339 `full-date` string, e.g. "2026-11-07".
    validation:
        - must be a valid RFC 3339 `full-date`
        - must not be the zero value

## Notes

- RFC 3339 section 5.6 defines the `full-date` production: https://www.rfc-editor.org/rfc/rfc3339#section-5.6
- A Date has no time of day or time offset.
- For go implementation use the standard `time` package (`time.DateOnly`), storing the date at midnight UTC so equal dates compare equal with `==`.
