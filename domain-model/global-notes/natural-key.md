# Natural Key

Most **Top Level Models** are identified by a generated [Id](../id/id.md). Some
models are already identified by a well-known, unique code, such as the ISO 3166
alpha-2 code of a country (`US`). Generating an Id for these would only add a
second identifier that has to be kept in sync with the first.

A **Top Level Model** marked with a **Natural Key** annotation names the field
that identifies it:

```markdown
**Top Level Model**

**Natural Key**: `Alpha2`
```

## Rules

- **No Id field.** The model does not have an `Id` field. The natural key field
  takes its place.
- **Unique and stable.** The natural key must be unique across every instance of
  the model and must never change. If the real-world code changes, that is a new
  instance (the spec may record the old code, e.g. `PreviousCode`).
- **Required.** The natural key field is always required and has validation that
  fully defines its format.
- **SQL.** The unique `id` column described in [SQL generation](sqlite.md) holds
  the natural key's value instead of a generated Id. The `rowid` primary key is
  unchanged.
- **DTO.** In a [DTO](dto.md), a field that refers to a model with a natural key
  holds that model's natural key value instead of an Id. For example, a
  reference to a Country holds `"US"`.
- **Lookup.** Looking a model up "by id" means looking it up by its natural key.
