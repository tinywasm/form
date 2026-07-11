# PLAN — Public `Submit()`, per-form `SetClass`, custom `Renderer`; finish moving tests to `tests/`

> This plan is dispatched via the CodeJob workflow. See skill: agents-workflow.

## Context (zero-context summary)

`tinywasm/form` generates HTML forms from Go structs implementing `model.Fielder`
(`Schema() []model.Field`, `Pointers() []any`, `Values() []any`). Fields whose
`model.Kind` also implements `input.Input` become form inputs. Each bound field
renders as `div.tw-field` containing the control + `span.tw-field-error`,
built by the unexported `fieldComponent` in `render_input.go`.

This plan closes four public-API gaps found while adopting the library in a
real project (mjosefa-cms):

1. **No programmatic submit** — the only submit trigger is the DOM event
   handler inside `Render()` (`render.go:56-76`). Tests and consumers must
   reach into unexported fields. → public `Submit() error`.
2. **No per-form CSS class** — `SetGlobalClass` exists (package-global), but a
   single module cannot style one form differently. → `SetClass(...string)`.
3. **No custom markup for custom inputs** — custom inputs (embedding
   `input.Base`) always render through `fieldComponent`'s fixed `HTMLName()`
   switch (input/textarea/select/radio/datalist). There is no escape hatch for
   e.g. a color-picker or composite widget. → optional `form.Renderer`
   capability interface.
4. Two test files still live at the module root because they need internals;
   after (1) and a black-box rewrite they move to `tests/`.

**Design doctrine (do not violate):** capability by interface, no registries,
no flags. `input.Input` itself works this way; `Renderer` follows the same
pattern.

**Ecosystem rules (mandatory):**
- WASM-compiled package — **no stdlib imports in library code** (`errors`,
  `strconv`, `strings`, `reflect` forbidden). Use `github.com/tinywasm/fmt`.
  (Test files under `tests/` may import stdlib — existing ones already do.)
- No `any`/`map` in public APIs. Value embedding only (never `*dom.Element`).
- **The `input` package must stay free of `dom` imports** (edge-safe contract)
  — this is why `Renderer` lives in package `form`, not `input`.
- Tests run with `gotest ./...` (native + wasm). Shared test logic uses the
  three-file pattern: `x.shared_test.go` (helpers), `x.back_test.go`
  (`//go:build !wasm`), `x.front_test.go` (`//go:build wasm`).
- Do NOT touch `docs/CHECK_PLAN.md` if present; never run `gopush`/`codejob`.

## Stage 0 — delete the `//ormc:storage` directives (no gate; run first)

The working tree carries `//ormc:storage <type>` comment directives above
every `input.*` constructor (unpublished phase-B work). That mechanism was
REJECTED (harness doctrine, `tinywasm/docs/ARNES_DE_CONSTRUCCION.md`): a
comment is prose the compiler cannot verify — it duplicates `Storage()` and
can silently contradict it. The ormc generator now resolves storage by
EXECUTING each kind's real `Storage()` through a temporary dependency probe
(see `ormc/docs/PLAN.md` stage 2) — no form-side artifact is needed at all.

Changes (comment/doc cleanup only — no behavior change):

- DELETE every `//ormc:storage` comment in `input/*.go`
  (`grep -rn "ormc:storage" .` must end empty).
- `docs/IMPLEMENTATION.md` ("Adding a New Built-in Input"): remove the
  directive step (current step 2) and renumber. Implementing `model.Kind`
  is all a kind needs — ormc extracts storage from the compiled contract.

## Stage 1 — public `Submit() error` in `form.go`; `Render()` delegates

The submit pipeline currently lives inline in the DOM event handler
(`render.go:56-76`). Extract it verbatim into a public method — the handler
must call the method so there is exactly one submit path.

Add to `form.go` (next to `Reset()`):

```go
// Submit runs the full submit pipeline programmatically: syncs input values
// into the bound struct, validates, and (if valid) fires the OnSubmit
// callback. Returns the first validation error, or nil if the submission
// was dispatched. The async result of the submission itself is delivered
// through the OnSubmit callback's done function.
func (f *Form) Submit() error {
	// Sync all values from signals to struct
	f.SyncValues(f.data)

	// Validate all (final check)
	if err := f.Validate(); err != nil {
		return err
	}

	if f.onSubmit != nil {
		f.submitting.Set(true)
		f.onSubmit(f.data, func(err error) {
			f.submitting.Set(false)
			if err == nil && !f.noResetOnSuccess {
				f.reset()
			}
		})
	}
	return nil
}
```

Replace the body of the `el.On("submit", ...)` handler in `render.go` with:

```go
el.On("submit", func(e dom.Event) {
	e.PreventDefault()
	f.Submit()
})
```

No other behavior change. The validation error is intentionally ignored in the
DOM handler (per-field error spans already surface it to the user); the return
value exists for programmatic callers.

## Stage 2 — per-form `SetClass(...string) *Form`

In `form.go`. The `Form.class` field already exists (initialized from
`globalClass` in `New`); it only lacks a public setter. Semantics: **append**,
same accumulation style as `SetGlobalClass`, so a module adds its scoping
class on top of the app-wide default:

```go
// SetClass appends CSS classes to this form (on top of any global classes
// set via SetGlobalClass). Chainable.
func (f *Form) SetClass(classes ...string) *Form {
	for _, c := range classes {
		if f.class != "" {
			f.class += " "
		}
		f.class += c
	}
	return f
}
```

Note `render.go` already emits `el.Class(f.class)` when non-empty — no render
change needed.

## Stage 3 — `form.Renderer`: optional custom markup for inputs

In `render_input.go`, add the capability interface and check it before the
`HTMLName()` switch inside `fieldComponent.Render()`:

```go
// Renderer is an optional capability for custom inputs that own their markup.
// The form still owns the field wrapper (div.tw-field), the error span, and
// validation: the widget must call onInput with the new value on user input —
// the form updates the value signal and runs live validation. The value
// signal carries the initial value and programmatic updates (SetValues).
type Renderer interface {
	RenderInput(value *dom.SignalString, onInput func(string)) *dom.Element
}
```

In `fieldComponent.Render()` (render_input.go:39), replace the switch block
with:

```go
if r, ok := fc.Input.(Renderer); ok {
	container.Child(r.RenderInput(fc.value, func(v string) {
		fc.value.Set(v)
		fc.validate(v)
	}))
} else {
	htmlName := fc.Input.HTMLName()
	switch htmlName {
	case "radio":
		fc.renderRadio(container)
	case "select":
		fc.renderSelect(container)
	case "datalist":
		fc.renderDatalist(container)
	default:
		fc.renderInput(container)
	}
}
```

The error span appending after the switch stays untouched — custom widgets can
never break the field structure (wrapper, error span, ids remain uniform).

Rationale for the `onInput` callback (vs exposing the err signal): it mirrors
exactly what the built-in handler does (`value.Set` + `validate`), keeps
signal coupling inside the form, and gives the widget a single obvious thing
to call. Do NOT pass the error signal to the widget.

`Renderer` MUST live in package `form` (it references `*dom.Element`; the
`input` package is dom-free by contract).

## Stage 4 — rewrite submit tests over the public API, move to `tests/`

Delete the white-box internals usage and move the three files to
`tests/submit.back_test.go`, `tests/submit.front_test.go`,
`tests/submit.shared_test.go`, all `package form_test`.

The `submitStruct` fixture stays as-is (it only uses public `model`/`input`
API). The rewritten `runSubmitTests` must cover, using ONLY public API
(`form.New`, `OnSubmit`, `Submit`, `SetValues`, `Reset`, `Input`,
`NoResetOnSuccess`):

1. **Callback receives data and done** — `f.OnSubmit(cb)`, set a valid value,
   `err := f.Submit()` → `err == nil`, callback was called, the
   `model.Fielder` passed to the callback is the bound struct (check the
   synced field value on the struct, proving SyncValues ran).
2. **Validation failure returns error and does not fire callback** — leave the
   `NotNull` field empty, `f.Submit()` → non-nil error, callback NOT called.
   (This is new coverage: the old handler swallowed this error.)
3. **NoResetOnSuccess retains values** — `f.NoResetOnSuccess()`, set value,
   submit with callback calling `done(nil)` → input still holds the value
   (`f.Input("nombre")` + `GetValue()` type assertion, as today).
4. **Default reset on success** — WITHOUT `NoResetOnSuccess`, submit with
   `done(nil)` → input value is cleared (proves the done-wrapper reset path,
   previously only testable via internals).
5. **Reset clears values** — keep the existing `TestSubmit_ResetClearsValues`
   case (already public-API only).

Note: the old assertions on `f.onSubmit != nil` and manual `f.onSubmit(...)`
invocation disappear — behavior is asserted through observable effects only.

## Stage 5 — rewrite `render_input_test.go` as black-box, move to `tests/`

Delete root `render_input_test.go` (package `form`). Create
`tests/render_input_test.go` (`package form_test`) preserving all existing
render cases (checkbox, datalist ×3, radio ×3, select ×3, textarea ×2,
search, text ×2) but exercising them through the real pipeline:
`form.New` → `SetOptions`/`SetValues` → `f.String()` → substring check.

Replace the old `buildInput` + unexported `fieldComponent` machinery with a
single-field fixture whose pointer type matches the kind's storage:

```go
// kindFixture is a one-field Fielder used to render a single input kind
// through the real form pipeline.
type kindFixture struct {
	inp     input.Input
	valText string
	valInt  int64
	valBool bool
}

func (k *kindFixture) Schema() []model.Field {
	return []model.Field{{Name: "tfield", Type: k.inp}}
}

func (k *kindFixture) Pointers() []any {
	switch k.inp.Storage() {
	case model.FieldInt:
		return []any{&k.valInt}
	case model.FieldBool:
		return []any{&k.valBool}
	default:
		return []any{&k.valText}
	}
}

func (k *kindFixture) Values() []any {
	switch k.inp.Storage() {
	case model.FieldInt:
		return []any{k.valInt}
	case model.FieldBool:
		return []any{k.valBool}
	default:
		return []any{k.valText}
	}
}
```

Per test case:

```go
fx := &kindFixture{inp: input.Select()} // constructor per case
f, err := form.New("tid", fx)
// fatal on err
if len(c.opts) > 0 {
	f.SetOptions("tfield", c.opts...)
}
if len(c.values) > 0 {
	f.SetValues("tfield", c.values...)
}
html := f.String()
// assert html contains c.contain (same expected substrings as today)
```

Keep the compact `rc` case table and the shared option vars (`opts12`,
`optsGender`). The case table's `t` field maps to a constructor via a switch
(`"Select"` → `input.Select()`, etc.) — same mapping as the old `buildInput`,
minus the `Clone` call (form.New clones internally).

Watch for identifier collisions inside `tests/` (single package
`form_test`): `setup_test.go`, `render.shared_test.go` etc. already declare
fixtures — the new names (`kindFixture`, `rc`, `opts12`, `optsGender`) must
not clash with existing declarations; rename if needed.

The expected substrings (`type='checkbox'`, `<datalist`, `value='1'`,
`list='`, `<label>`, `value='m'`, `checked`, `<select`, `<option`,
`selected`, `<textarea`, rendered value text, `type='search'`, `<input`,
`type='text'`) are id-independent and must remain identical. If any case
fails only because of the form-id prefix in generated ids, adjust nothing in
the library — the assertions above don't depend on ids.

## Stage 6 — tests for `SetClass` and `Renderer` (in `tests/`)

All `package form_test`, public API only:

1. **SetClass** — `form.New(...).SetClass("cms-form")` → `f.String()`
   contains `class='` with `cms-form`; combined with `SetGlobalClass` both
   classes appear (beware: `SetGlobalClass` mutates package state — set it
   before `New` and don't leak it into other tests; if needed, keep this case
   minimal and rely on per-form class only).
2. **Renderer** — declare a custom input in the test file: embed `input.Base`,
   implement `form.Renderer` returning a distinctive element (e.g.
   `dom.NewElement("div").Class("my-widget")`). Bind it via a one-field
   fixture → `f.String()` must contain BOTH the custom markup
   (`class='my-widget'`) AND the standard error span (`tw-field-error`),
   proving the form still owns the field structure.
3. **Renderer validation wiring** (WASM-independent part only): the `onInput`
   callback is exercised in the browser; from native tests, assert that a
   `Renderer` input still participates in `f.Validate()` (set an invalid
   value via `SetValues`, expect error) — validation must not depend on the
   rendering path.

## Stage 7 — documentation

- `README.md` "Form Methods" table: add rows for
  `Submit() error` — "Runs sync + validate + OnSubmit callback
  programmatically; returns first validation error" — and
  `SetClass(...string) *Form` — "Appends CSS classes to this form (on top of
  SetGlobalClass)".
- `README.md` "Styling" section: mention `SetClass` next to `SetGlobalClass`
  (point 3 of that section).
- `README.md` "Custom Inputs" section: add one sentence — custom markup is
  possible by implementing `form.Renderer`; link `input/README.md`.
- `docs/API.md`: add `## (*Form).Submit()` (pipeline order SyncValues →
  Validate → OnSubmit, return contract, DOM handler delegates to it) and
  `## form.Renderer` (contract: form owns wrapper/error span/validation;
  widget calls `onInput`; lives in `form` because `input` is dom-free).
- `input/README.md` ("Creating a Custom Input" section): append a short
  "Custom markup" subsection showing a minimal `form.Renderer`
  implementation on a custom input.

## Code quality checklist (mandatory)

- No hardcoded repeated strings in library code; the only new literals are doc
  comments. Test field name `"tfield"` may be a local `const` in the test
  file.
- No stdlib in library changes (none needed).
- No `any`/`map` in the new public API (`Submit() error`,
  `SetClass(...string) *Form`, `Renderer` — all compliant).
- Single submit path: `Render()`'s handler MUST call `f.Submit()`; do not
  leave duplicated pipeline code in `render.go`.
- `Renderer` in package `form` only — adding a `dom` import to `input/` is a
  hard failure.
- `gotest ./...` green (native + wasm) before finishing.
- Breaking change: none (pure additions + test relocation). Minor version
  bump.

## Acceptance criteria

1. `form.Submit() error` exists; `render.go`'s submit handler is a two-line
   delegation to it.
2. `form.SetClass(...string) *Form` appends per-form classes; rendered
   `<form>` carries global + per-form classes.
3. `form.Renderer` exists in package `form`; a custom input implementing it
   renders its own markup inside `div.tw-field` with the standard error span;
   inputs NOT implementing it render exactly as before; `input/` has no `dom`
   import.
4. `submit.*_test.go` and `render_input_test.go` no longer exist at the module
   root; their coverage lives in `tests/` as `package form_test`, using only
   exported API.
5. New coverage: validation failure returns error without firing the callback;
   default reset-on-success observable via public API; SetClass and Renderer
   cases from Stage 6.
6. `grep -rn "package form$" *.go` at module root matches only non-test files.
7. `gotest ./...` green (native + wasm); README, docs/API.md and
   input/README.md updated per Stage 7.
8. Stage 0: `grep -rn "ormc:storage" .` empty (code and
   `docs/IMPLEMENTATION.md`); no other file changed by that stage.

## Stages

| Stage | File(s) | Action |
|---|---|---|
| 0 | `input/*.go`, `docs/IMPLEMENTATION.md` | delete every `//ormc:storage` comment (mechanism rejected; ormc probes `Storage()` directly) |
| 1 | `form.go`, `render.go` | add `Submit() error`; handler delegates |
| 2 | `form.go` | add `SetClass(...string) *Form` (append semantics) |
| 3 | `render_input.go` | `Renderer` interface + capability check in `fieldComponent.Render()` |
| 4 | `submit.{shared,back,front}_test.go` → `tests/` | rewrite over public API |
| 5 | `render_input_test.go` → `tests/render_input_test.go` | black-box rewrite via `form.New` + `String()` |
| 6 | `tests/` | SetClass + Renderer coverage |
| 7 | `README.md`, `docs/API.md`, `input/README.md` | document Submit, SetClass, Renderer |
