# tinywasm/form — Plan: Signal-Bound Fields, No Lifecycle Hooks

> **Master:** tinywasm/docs/PLAN.md · **Engine:** tinywasm/dom/docs/PLAN.md
> **Module:** `github.com/tinywasm/form`
> **Type:** Breaking-aligned migration.

---

## Prerequisites

```bash
# Canonical test runner (WASM tests run against a real DOM). Required: external agents have no global gotest.
go install github.com/tinywasm/devflow/cmd/gotest@latest
```

## Development Rules

- **Documentation First:** update `docs/ARCHITECTURE.md` (events/validation) before code.
- **WASM only:** reactive code in `//go:build wasm` (`mount.go`); keep `mount_stub.go`'s method set in
  sync (it shrinks — see below).
- **DOM access only via `tinywasm/dom`.** No `syscall/js`.
- **Tests:** `gotest` (never `go test`); stdlib only; dual WASM/stdlib. Publish with `gopush 'msg'`.
- **Minimal public API:** export only what a component *user* types; unexport anything only this package uses (helpers, field models, single-use constants). State lives in unexported fields exposed via signals.

## Signals API recap (from the dom engine — self-contained)

```go
s := dom.NewString("")            // observable string cell; Set patches bound nodes surgically
in.Bind(s)                        // two-way <input> <-> signal (cursor/IME safe; node never replaced)
errSpan.BindText(errSig)          // a single text node bound to a SignalString
el.BindClassFunc("visible", func() bool { return errSig.Get() != "" }) // computed class, auto-tracked (no deps)
```

State the UI shows lives in a typed `Signal` (**no generics**). No `Update()`, no `OnMount`.

---

## Context

`Form.OnMount()` (mount.go:11-89) attaches **delegated** `input`/`change`/`submit`
listeners to the form root by ID and updates errors imperatively via `dom.Get(errID).SetText(...)`
(mount.go:30-39). That whole imperative, lifecycle-dependent machine is replaced by
per-field signal bindings: surgical by construction, IME-safe, nothing to re-bind.

---

## Change — Bind fields to signals; delete `OnMount`/`OnUnmount`

1. Give each input a value signal and an error signal (created when the form is built, e.g. keyed by
   field ID). In `Render()`:

```go
in := Input(field.Type).Bind(field.value).
	On("input", func(Event) { field.err.Set(errText(field.input.Validate(field.value.Get()))) })
errSpan := Span(clsErr).BindText(field.err).
	BindClassFunc("tw-field-error--visible", func() bool { return field.err.Get() != "" }) // auto-tracked; no deps
```

   On each keystroke only `field.err`'s text node + class patch — **the `<input>` node is never
   touched**, so focus/cursor/IME composition (accents, CJK) are preserved. No whole-form re-render.

2. Submit: bind the `<form>` element `.On("submit", c.onSubmit)`. `onSubmit` syncs signal values into
   `f.data` (`SyncValues`), validates, toggles a `submitting *SignalBool` bound to the button's
   `disabled` (`BindAttrBool`) and label (`BindText`), calls `f.onSubmit`, and on success `reset()`
   sets each value signal back to its default (patches inputs).

3. **Delete `OnMount` and `OnUnmount`** from `mount.go`; remove their stubs from `mount_stub.go`. The
   form implements no lifecycle interface (optionally `Init(ctx)` only if it ever needs async defaults).

> Two-way `Bind` already guards against cursor jumps (skips patching the active input), so live
> validation no longer risks breaking IME — the original reason this plan thread started.

---

## Documentation (do FIRST)

- **`docs/INTERACTIVITY_AND_MOUNTING.md`**: this describes the **old** delegated-`OnMount` mounting
  model and is now obsolete — rewrite it to the signal-binding model (or fold into ARCHITECTURE.md and
  delete, updating the README index).
- **`docs/DESIGN.md`**: record the decision — per-field `Signal` binding replaces imperative
  delegation; rationale includes the IME/cursor-safety of two-way `Bind` (the issue that drove this).
- **`docs/API.md`**: update the public API to the `Bind`/`BindText`/`BindClass` surface; remove
  `OnMount`/`OnUnmount`.
- **`README.md`**: re-index `docs/` (every file linked) after the above.

## Tests — frequent use cases (`gotest`)

Stdlib assertions only; dual WASM/stdlib. Cover the everyday form patterns:

- **stdlib:** validation/sync logic over value signals (no DOM).
- **wasm (real DOM):**
  - **live validation:** typing into a bound input changes **only** the error text node + class
    (assert the `<input>` keeps identity and cursor position).
  - **submit lifecycle:** `submitting` toggles the button's disabled/label; on success `reset()`
    patches value signals back into inputs.
- **In-browser (tinywasm MCP):** live validation updates without input re-render; **IME check** —
  type `á`, `ñ`, and a CJK character; composition is not broken. `browser_get_errors` clean.

## Done When

- Form implements only `Render()` (+ optional `Init`); `OnMount`/`OnUnmount` deleted from `mount.go`
  and `mount_stub.go`. Live validation patches a single text node; inputs are never re-rendered.
- **Docs:** `INTERACTIVITY_AND_MOUNTING.md`, `DESIGN.md`, `API.md` updated; `README.md` re-indexed.
  **Tests:** use-case tests (incl. the IME check) pass under `gotest`.
