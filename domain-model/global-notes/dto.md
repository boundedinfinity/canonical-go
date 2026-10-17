# Data Transfer Objects (DTO)

Every **Top Level Model** has a matching Data Transfer Object (DTO). A DTO has
the same fields as its model, except that any field holding another **Top Level
Model** holds that model's [Id](../id/id.md) instead of embedding it.

A DTO lets a model be stored or sent without copying the models it refers to.
Each referenced model is stored once in its own table (see
[SQL generation](sqlite.md)) and referred to everywhere else by `id`.

The DTO is derived from the model's spec. Individual specs do not describe it.

## Rules

- **References become IDs.** A field whose type is a **Top Level Model** becomes
  an [Id](../id/id.md). The cardinality is kept: a single model becomes a single
  ID, and a `List` of models becomes a `List` of IDs. If the referenced model
  has a [Natural Key](natural-key.md), the field holds that key's value instead
  of an Id (e.g. a Country reference holds `"US"`).
- **Other fields are unchanged.** Scalars, lists of scalars, and models that are
  *not* **Top Level Models** (e.g. [Affix](../person/affix/affix.md) or
  [Category](../person/affix/category.md), which are base types) are copied
  exactly as the model declares them.
- **Optionality is unchanged.** An optional reference is an optional ID (or an
  optional list of IDs).
- **Field names are unchanged.** A reference field keeps the model's field name
  and [JSON](json.md) key: `prefix` holds a list of IDs in the DTO, rather than
  a list of objects.
- **Validation matches the model.** Every field rule from the model's spec
  applies to the DTO. Each reference ID must be a valid [Id](../id/id.md), and a
  referenced list must contain at least 1 element if present. Cross-field rules
  from the spec's `# Validation` section also apply.
- **Models with no references still get a DTO.** Their DTO has the same fields
  as the model, so callers can always use the DTO form of any **Top Level
  Model**.

## Example

A [Name](../person/name/name.md) holds lists of [Prefix](../person/affix/prefix.md)
and [Suffix](../person/affix/suffix.md) models:

```json
{
  "id": "b12e7628-35b0-4ae5-b2d9-0f34f4828229",
  "prefix": [
    {
      "id": "0b6f1c9e-2d4a-4f3b-9c8e-7a5d3e1f2b4c",
      "name": "Dr.",
      "category": { "id": "9a7e4c2b-1d3f-4b5a-8e6c-2f0d9b8a7c6e", "name": "Honorific" }
    }
  ],
  "first": ["Ada"]
}
```

The Name DTO holds their IDs instead:

```json
{
  "id": "b12e7628-35b0-4ae5-b2d9-0f34f4828229",
  "prefix": ["0b6f1c9e-2d4a-4f3b-9c8e-7a5d3e1f2b4c"],
  "first": ["Ada"]
}
```

The Prefix DTO likewise holds `"category": "9a7e4c2b-…"` instead of the
category object.

## Conversion

- **Model → DTO** replaces each referenced model with its ID. It cannot fail.
- **DTO → Model** looks up each referenced ID with a resolver supplied by the
  caller (for example, a database query), then validates the rebuilt model. It
  fails if any ID cannot be resolved or the rebuilt model is invalid, and
  reports every failure together.

## Go implementation

- The DTO for model `Xxx` is `XxxDto`, declared in the model's package in a
  `<model>-dto.go` file next to the model (e.g. `name.NameDto` in
  `lang/go/person/name/name-dto.go`, `affix.PrefixDto` in
  `lang/go/person/affix/prefix-dto.go`). Its test is `<model>-dto_test.go`.
- A reference field has type `id.ID`, `[]id.ID`, `optioner.Option[id.ID]`, or
  `optioner.Option[[]id.ID]`.
- `XxxDto` has `Validate()` and validating `MarshalJSON` / `UnmarshalJSON`, the
  same as the model. It reports failures with the model's sentinels
  (`ErrInvalid`, `ErrPrefix`, etc.).
- `Xxx.ToDto() XxxDto` converts a model to its DTO.
- `XxxDto.ToModel(...) (Xxx, error)` converts back. It takes one
  `resolver.Resolver[K, T]` (from `lang/go/resolver`) per referenced model type,
  in field order, where `K` is the reference's key type: `id.ID`, or the
  [Natural Key](natural-key.md)'s type (e.g. `string` for a Country). A DTO with
  no references has `ToModel()` with no arguments.
- Resolution failures wrap `resolver.ErrResolve` and the resolver's own error.
