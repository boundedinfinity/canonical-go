# Validation

## Keep validation in the model's `Validate` method

All of a model's validation lives in its `Validate` method. That includes every
field rule from `## Programmatic Information`, the per-element checks on lists,
and the cross-field rules from the spec's `# Validation` section. A reader can
then see every rule for a model in one place, in the same order as the spec.

Do not split validation into separate helper functions (e.g. a function that
validates each element of a list, or one that validates a group of fields),
unless:

- the model's spec says to in its `## Notes`, or
- there is a good reason, which is written up and agreed before the split is
  made (for example, logic that is complex enough to need its own tests).

This is not a split, and is expected:

- **Declarative field rules.** Package-level validators that state a field's
  rule, such as `validation.String[string]("name").Min(1).Max(30)`. They declare
  the rule once; `Validate` still decides when to apply it.
- **Nested models.** Calling `Validate` on a nested model, as the spec's "call
  the validation on ..." rules require. That validation belongs to the nested
  model.

## DTOs

A [DTO](dto.md) has its own `Validate` method that follows the same rule: all of
its validation is written out in full, even where it repeats the model's.

The one exception is a DTO whose fields are identical to its model's (a model
that holds no other **Top Level Model**). It may delegate to the model's
`Validate`, because the rules are the same by definition, so there is nothing
to keep in sync.
