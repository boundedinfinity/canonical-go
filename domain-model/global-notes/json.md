# JSON serialization

## Go implementation

**JSON tags**: kebab-case of the field name (e.g. `StreetAddress1` → `street-address1`, `AccountNumber` → `account-number`), `omitzero` on every optional field. This applies to every language implementation.

Use `omitzero`, not `omitempty`. Optional fields are `optioner.Option[T]`, which is a struct, and `omitempty` never omits structs, so an unset option would be written as `null`. `omitzero` omits any zero value, and an unset option is the zero value.

## Typescript implementation

**JSON tags**: kebab-case of the field name (e.g. `StreetAddress1` → `street-address1`, `AccountNumber` → `account-number`)on every optional field. This applies to every language implementation.

For typescript classes types manually encode the serialization of the class internals into an internal interface, and converte the camel-cased class files to the kebab-cased field names of the interface.