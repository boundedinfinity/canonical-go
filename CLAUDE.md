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
- **Struct models:** `lang/go/label` is the reference implementation for structs that hold several fields.
  - Spec fields map to exported Go fields with lowercase JSON tags. A field whose spec type links to another spec uses that package's type, so `[ID](...)` becomes `id.ID`.
  - Fields marked `optional: true` use `optioner.Option[string]` with the `omitzero` JSON tag.
  - Length bounds become exported constants (`NameMin`, `NameMax`). Lengths are counted in runes.
  - `Validate()` checks every field and returns all failures together via `errorer.Join`. Each failure wraps `ErrInvalid` plus a per-field sentinel such as `ErrName`. It calls `Validate()` on any nested types.
  - `New(...)` creates a new `id.New()` and validates the result. `MustNew` panics instead of returning an error.
  - `UnmarshalJSON` decodes through an alias type and then calls `Validate()`.
- **Dependencies:** use the standard library plus go-commoner's `errorer` (errors, via `errorer.Func(ErrInvalid, ...)`) and `optioner` (optional fields). Tests also use `testify`. Don't import `idiomatic/` packages. Go 1.27 includes a standard `uuid` package, which `lang/go/id` uses.
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
