# ID

This data type represents a unique identifier, formatted as a version 4 UUID.

## Programmatic Information

- ID
    type: string
    desc: A UUIDv4 string, e.g. "f47ac10b-58cc-4372-a567-0e02b2c3d479".
    validation:
        - must be a valid UUID

## Notes

- UUID version 4 (random) is defined by RFC 9562, which obsoletes RFC 4122: https://www.rfc-editor.org/rfc/rfc9562
- For go implementation use the built-in uuid package