# JSON serialization

## Go implementation

**JSON tags**: kebab-case of the field name (e.g. `StreetAddress1` → `street-address1`, `AccountNumber` → `account-number`), `omitempty` on every optional field. This applies to every language implementation.

## Typescript implementation

**JSON tags**: kebab-case of the field name (e.g. `StreetAddress1` → `street-address1`, `AccountNumber` → `account-number`)on every optional field. This applies to every language implementation.

For typescript classes types manually encode the serialization of the class internals into an internal interface, and converte the camel-cased class files to the kebab-cased field names of the interface.