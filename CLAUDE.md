# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

`github.com/boundedinfinity/canonical-go` is a library of "canonical" domain models for real-world concepts such as people, phone numbers, currency, bookkeeping, vehicles, measurements, and ISO/ANSI specifications. There is no application or main binary.

The project is driven from the `domain-model/` directory. Its markdown specs are the source of truth, and language implementations are built from them under `lang/<language>/`.

## domain-model → lang/go

Every spec at `domain-model/<path>/<name>.md` maps to a Go file at `lang/go/<path>/<name>.go`. The test goes next to it at `lang/go/<path>/<name>_test.go`. For example, `domain-model/id/id.md` maps to `lang/go/id/id.go`.

- **Package:** the package name is the last directory segment. Hyphenated directories need a valid Go package name, so `mailing-address` becomes `package mailingaddress`. The import path is `github.com/boundedinfinity/canonical-go/lang/go/<path>`, and these packages live in the root `go.mod` module.
- **Template:** `lang/go/id` is the reference implementation, so copy its shape for new types.
  - A package doc comment names the source spec (`// Source of truth: domain-model/...`).
  - Each type gets a constructor, `Parse` / `MustParse`, `String`, and `IsZero`, plus explicit `MarshalJSON` / `UnmarshalJSON` (encode as the spec's declared type, validate on decode, treat JSON `null` as a no-op) and `MarshalText` / `UnmarshalText` for map keys and other text encoders.
  - Validation rules from the spec's `validation:` list are enforced in `Parse` / `UnmarshalText`. Failures return an error that wraps a package-level `ErrInvalid` sentinel.
  - Tests are table-driven, live in `<pkg>_test`, and cover the valid cases, every validation rule, and a JSON round-trip.
- **Global notes:** `domain-model/domain-model.md` and `domain-model/global-notes/` apply to every model. `json.md` covers JSON tags (kebab-case field names, omitted when empty), `sqlite.md` covers SQL generation (rowid primary key plus a unique `id` column for every **Top Level Model**), `dto.md` covers the DTO every **Top Level Model** gets, `validation.md` covers where validation code lives, and `natural-key.md` covers **Top Level Models** identified by an existing code (e.g. an ISO country code) instead of a generated `Id`. Read them before implementing a model.
- **Struct models:** `lang/go/label` (scalar fields) and `lang/go/person/name` (list fields) are the reference implementations for structs that hold several fields.
  - File order: `var` block (validators, errors), type, methods (`Validate`, JSON), then constructors (`New`, `MustNew`) at the bottom.
  - Spec fields map to exported Go fields with kebab-case JSON tags. A field whose spec type links to another spec uses that package's type, so `[Id](...)` becomes `id.ID`.
  - Optional fields use `optioner.Option[T]` with the `omitzero` JSON tag, never `omitempty`: `omitempty` does not omit structs, so a `None` option would be written as `null`. Optional lists (`type: List strings`) are `optioner.Option[[]string]`, built with `optioner.OfSlice`, which turns an empty slice into `None`. A JSON `[]` decodes as present-but-empty and fails the "at least 1 element" rule.
  - Field rules are declared once as package-level validators built with `lang/go/validation`, e.g. `validation.String[string]("name").Min(1).Max(20)` or `validation.Slice[string]("first").Min(1).Each(elemValidator)`. String lengths count runes.
  - `Validate()` checks every field (optional ones only when present) and returns all failures together via `errorer.Join`. Each failure wraps `ErrInvalid` plus a per-field sentinel such as `ErrName`. `ErrId` aliases `id.ErrInvalid`. It calls `Validate()` on any nested types, and cross-field rules from the spec's `# Validation` section get their own sentinel (e.g. `name.ErrRequired`). Keep every check inline in `Validate()`, including per-element loops over lists. Don't split validation into helper functions (no `validateEach` / `validateParts`) unless the spec's `## Notes` says to, or you've proposed a good reason and the user agreed. Package-level validators and calls to nested `Validate()` are fine.
  - `New(...)` creates a new `id.New()` and validates the result. `MustNew` panics instead of returning an error.
  - `MarshalJSON` / `UnmarshalJSON` go through an alias type. `UnmarshalJSON` calls `Validate()` and leaves the receiver untouched on error or JSON `null`.
- **DTOs:** every **Top Level Model** `Xxx` also gets an `XxxDto` in `<model>-dto.go` (test: `<model>-dto_test.go`), following `domain-model/global-notes/dto.md`. `lang/go/person/name/name-dto.go` (lists of references) and `lang/go/person/affix/prefix-dto.go` (a single reference) are the references; `lang/go/label/label-dto.go` shows a model with no references.
  - A field holding another **Top Level Model** holds its `id.ID` instead, keeping the cardinality and optionality (`id.ID`, `[]id.ID`, `optioner.Option[[]id.ID]`, ...) and the same JSON key. Everything else, including embedded non-top-level base types like `affix.Affix`, is copied as-is.
  - The DTO has `Validate()` (same rules and sentinels as the model, with each reference checked as an ID) and validating JSON. Write its checks out in full, even where they repeat the model's (see `name-dto.go`). The exception is a DTO with no references, which can convert to the model and reuse its `Validate` (see `label-dto.go`).
  - `Xxx.ToDto()` converts to the DTO. `XxxDto.ToModel(...)` takes one `resolver.Resolver[K, T]` (`lang/go/resolver`) per referenced type in field order (`K` is `id.ID`, or the natural key's type), resolves with `resolver.One` / `All` / `Option` / `OptionAll` (failures wrap `resolver.ErrResolve`), reports every failure, then validates the rebuilt model. Tests can build resolvers with `resolver.Func`.
- **Natural keys and data sets:** `lang/go/specification/iso-3166` is the reference for a **Natural Key** model (no `Id`; the key field is validated like any other) and for a spec with a data file. The package keeps a copy of the spec's JSON, refreshed by `go generate ./lang/go/specification/...` (a `//go:generate cp` line), and embeds it with `go:embed`. A test fails if the copy drifts from `domain-model/`. The data loads into package-level lists (`CountryCodes`, `Subdivisions`) that provide the spec's `## Lookup` operations and implement `resolver.Resolver` by natural key. Format rules use `validation.String(...).Match(regexp)` and `.URL("http", "https")`.
- **Dependencies:** use the standard library plus go-commoner's `errorer` (errors, via `errorer.Func(ErrInvalid, ...)`) and `optioner` (optional fields), and the in-repo `lang/go/validation` package for field rules. Tests also use `testify`. Don't import `idiomatic/` packages. Go 1.27 includes a standard `uuid` package, which `lang/go/id` uses.
- **Spec format:** a spec's `## Programmatic Information` section lists each field with its `type`, `desc`, and `validation`. `## Notes` holds references and implementation hints, such as which library to use.
- **Change order:** when a spec changes, update the `.md` first and then the generated language code. Don't add behavior to `lang/go` that the spec doesn't describe. If the spec is ambiguous, ask or update the spec.

## Commands

- `go test ./lang/...`: test all of the new spec-driven code. This should always pass.
- `go test ./lang/go/id -run Test_Parse`: run a single test
- `go vet ./lang/...`, `gofmt -l lang`
- `just test`: `go test ./...` across the whole repo, including legacy code that is known to fail
- `just commit` / `just push` / `just publish m=<tag>`: git add-all, commit, push, and tag helpers

## Legacy code (`idiomatic/`)

`idiomatic/<domain>/` holds the older hand-written models. Use it as reference material when building `lang/go` packages, but don't extend it.

- It depends on `go-commoner`, through `replace github.com/boundedinfinity/go-commoner => ../go-commoner` in `go.mod`, so that sibling repo must be checked out.
- It doesn't fully build or pass tests: `specification/iso/iso3166/internal` has no `main`, and tests fail in `measurement/length`, `phone`, `check_algo`, and `sql`. Don't assume a red `go test ./...` was caused by your change.
- Its conventions differ from `lang/go`. It uses `ider.Id` instead of `id.ID`, and static data embedded in `*.data.go` / `*.data.json` files that are loaded in `init()`.
- The `go-commoner` dependency it shares with `lang/go` requires the sibling checkout at `../go-commoner`.

## Style

- Methods use `this` as the receiver name. Staticcheck (`staticcheck.conf`) disables ST1003 and ST1006 to allow it.
- `.editorconfig` sets 4-space indentation and LF line endings. `gofmt` still uses tabs for `.go` files.
