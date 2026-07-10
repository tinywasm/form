# PLAN — Kind unification (phase B): `input` decorates `model.Kind`, the widget registry dies

> This plan is dispatched via the CodeJob workflow. See skill: agents-workflow.
> Phase B of `tinywasm/docs/KIND_UNIFICATION_MASTER_PLAN.md` (Kind unification wave). Requires
> the published phase-A `tinywasm/model`. Runs parallel to orm/sqlt/postgres/mcp.

## Context (zero-context summary)

Phase A changed `tinywasm/model`: `Field.Widget` and the `model.Widget`
interface are **gone**; `Field.Type` is now the interface

```go
type Kind interface {
    Storage() FieldType
    Name() string
    Validate(value string) error
}
```

and `model` ships base kinds (`model.Text()`, `model.Int()`, …) with
fail-closed validation. Definitions declare form fields by choosing an
`input.*` kind directly:

```go
{Name: "email", Type: input.Email(), NotNull: true}
```

This repo's role changes from "provides optional widgets looked up beside
the type" to "**decorates** kinds with rendering". A field is a form field
iff its `Kind` also implements this package's `input.Input` contract —
capability by interface, no registry, no flags.

**Ecosystem rules:** WASM-compiled package — no stdlib (`tinywasm/fmt`), no
`any`/`map` in public APIs, typed constants, `gotest` (native + wasm),
value embedding only (never `*dom.Element`).

## Stage 1 — `input.Input` embeds `model.Kind`, owns `Clone`

In `input/interface.go`:

- `Input` currently embeds `model.Widget` (deleted upstream). It now embeds
  `model.Kind` and declares `Clone(parentID, name string) Input` itself —
  clone-for-render-positioning is a form concern, which is exactly why phase
  A evicted it from the schema package.
- In `input/base.go`: rename `Base.Type() string` → `Name() string`
  (satisfies `Kind.Name`); add `Storage() model.FieldType`; `Validate`
  stays — `Base` already embeds `model.Permitted` and delegates to it,
  which is exactly the intended pattern (kinds implement their baselines
  with the `Permitted` engine). Update `Clone` signatures to return `Input`.
- Note: `Base.Validate`'s `Required` check duplicates `Field.NotNull` when
  reached through `Field.Validate` (which checks NotNull first). Keep
  `Required` for the HTML `required` attribute rendering, but the
  validation path must not error twice or diverge — if simplifying, the
  presence check belongs to `Field.Validate`, not the kind.

## Stage 2 — every `input.*` constructor becomes a storage-annotated kind

For each constructor (`text.go`, `email.go`, `password.go`, `phone.go`,
`number.go`, `date.go`, `hour.go`, `textarea.go`, `search.go`, `rut.go`,
`ip.go`, `address.go`, `gender.go`, `datalist.go`, `select.go`, `radio.go`,
`checkbox.go`, `filepath.go`, `tilde.go`):

- Add the ormc directive comment immediately above the constructor:
  `//ormc:storage text` (the syntax ormc's phase-B plan consumes).
- Storage assignment: **`text` for all string-valued inputs** (the vast
  majority). `number.go` → `int`. `checkbox.go` → `bool` if it binds a
  single boolean (verify how existing consumers pair it; if today's usage
  pairs it with `FieldText`, keep `text` and note it). If any input's
  natural storage contradicts how existing consumers used it, **STOP and
  report** — that decision goes back to the master plan, not guessed here.
- Ensure each constructor's returned value implements the full new `Input`
  (Kind + rendering) — compile-time asserted:
  `var _ Input = Email()` per file or grouped in one test.

## Stage 3 — form binding by capability, registry deleted

- `form.New` currently selects fields via `Field.Widget != nil`. It now
  binds every field whose `Type` asserts to `input.Input`:

  ```go
  if in, ok := field.Type.(input.Input); ok { /* clone, position, render */ }
  ```

  Fields with base kinds (`model.Text()` etc.) are simply not form inputs —
  same semantics as the old `Widget: nil`, but validation no longer
  disappears with them.
- Delete `registry.go` (`RegisterInput` and its map): name-based lookup is
  dead machinery — binding is by interface capability. Grep the repo for
  `RegisterInput` and remove every trace.
- Validation call sites: keep calling `Field.Validate` (phase A made it
  delegate to the kind unconditionally) — do not call `Kind.Validate`
  directly from form, the Field wrapper owns NotNull/Permitted ordering.

## Stage 4 — tests

- `form.New` over a fixture Definition mixing `input.*` and `model.*` kinds:
  exactly the `input.*` fields become inputs, in schema order.
- SSR render (`SetSSR(true).String()`) emits `name='<field>'` for each bound
  field (contract consumed by downstream endpoint handlers).
- Compile-time interface assertions for every constructor.
- A base-kind field still validates on submit (fail-closed regression:
  the old `Widget: nil` = no-validation hole must be provably gone).
- `grep -rn "RegisterInput"` → empty; `gotest ./...` green (native + wasm).

## Stage 5 — documentation

- `README.md`: authoring rule — choose the kind once in the Definition;
  `input.*` = form + validation, `model.*` = validation only.
- Update any doc mentioning `RegisterInput` or `Field.Widget`.

## Harness checklist (mandatory)

- Pin the phase-A `tinywasm/model` version.
- The `//ormc:storage` directive strings must match ormc's syntax exactly —
  they are part of the wave's shared contract.
- No `any`/`map` public API; no stdlib; typed constants.
- If the phase-A `Kind` contract is insufficient for rendering needs,
  **STOP and report** to the master plan — never re-add a parallel widget
  slot locally.
- Breaking change: next minor version. No deprecated shims.

## Acceptance criteria

1. Every `input.*` constructor implements `model.Kind` + `input.Input` and
   carries its `//ormc:storage` directive.
2. `form.New` binds by interface capability; `registry.go` is gone.
3. Fixture SSR HTML exposes the bound field names; base-kind fields render
   nothing but still validate.
4. `gotest ./...` green (native + wasm).

## Stages

| Stage | File(s) | Action |
|---|---|---|
| 1 | `input/interface.go`, `input/base.go` | embed `model.Kind`, own `Clone`, `Name()`/`Storage()` |
| 2 | `input/*.go` (all constructors) | storage directives + full-contract compliance |
| 3 | `form.go`/binding path, `registry.go` (delete) | capability-based binding |
| 4 | `*_test.go` | binding, SSR contract, fail-closed regression |
| 5 | `README.md`, docs | authoring rule |
