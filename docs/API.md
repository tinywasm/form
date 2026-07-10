# API Reference

See `README.md` for the consolidated API. This file contains additional detail.

## `form.New` — Widget Resolution Detail

```go
f, err := form.New("content", data) // data implements model.Fielder
// -> f.GetID() == "content." + resolveStructName(data)
// -> f.String() renders all fields that have a Widget in data.Schema()
```

For each field in `data.Schema()`:
1. `field.IsPK() && field.IsAutoInc()` → skip (auto-increment PKs not editable).
2. `field.Widget == nil` → skip (no UI binding).
3. `field.Widget.Clone(formID, fieldName).(input.Input)` → positioned input.
4. `field.NotNull` → `SetRequired(true)` on the input.
5. Current value bound via `fmt.ReadValues()` + `SetValues()`.

## `(*Form).Submit()` — Submit Pipeline

Runs the full submission pipeline programmatically:
1. `SyncValues(data)` — syncs input values to the struct.
2. `Validate()` — performs final validation.
3. If valid and `OnSubmit` is set:
   - Sets the `submitting` signal to `true`.
   - Fires the `onSubmit(data, done)` callback.
   - Resets the form on `done(nil)` unless `NoResetOnSuccess()` was called.

The DOM `submit` handler delegates to this method after `PreventDefault()`.
Returns the first validation error, or `nil` if submission was dispatched.

## `(*Form).Validate()` — Validation Detail

- Skips fields with `SkipValidation` set to true in the input.
- Pulls values from reactive signals.
- Calls `inp.Validate(val)` (promoted from `model.Kind`).
- Returns the **first** error encountered.

## `(*Form).SyncValues(data model.Fielder)` — Binding Detail

Synchronizes input values back to the struct pointers provided by `data.Pointers()`.
Supports `model.FieldText`, `model.FieldInt`, `model.FieldFloat`, and `model.FieldBool`.

## `(*Form).ValidateData(action byte, data model.Fielder)` — Server-side Validation

Validates the provided `data` using the form's input rules. Satisfies `crudp.DataValidator`.

## `fmt.Permitted` — Validation Engine

```go
type Permitted struct {
	Letters    bool     // a-z, A-Z (and ñ/Ñ)
	Tilde      bool     // á, é, í, ó, ú
	Numbers    bool     // 0-9
	Spaces     bool     // ' '
	BreakLine  bool     // '\n'
	Tab        bool     // '\t'
	Extra      []rune   // additional allowed characters
	NotAllowed []string // disallowed substrings
	Minimum    int      // minimum length
	Maximum    int      // maximum length
}
```

Error messages from `Permitted.Validate(name, text)`:
- `"{name} minimum {min} chars"` — value shorter than Minimum
- `"{name} maximum {max} chars"` — value longer than Maximum
- `"space not allowed"` — space when Spaces=false
- `"character {X} not allowed"` — disallowed character

## `form.Renderer` — Custom Input Markup

Optional capability interface for custom inputs that own their markup:

```go
type Renderer interface {
    RenderInput(value *dom.SignalString, onInput func(string)) *dom.Element
}
```

- **Contract**: The form still owns the field wrapper (`div.tw-field`), the error span, and validation.
- **Wiring**: The widget must call `onInput` with the new value whenever the user interacts with the control; the form then updates the value signal and runs live validation.
- **Location**: Defined in package `form` (not `input`) because it references `*dom.Element`.

## Namer Interface

Fielder types can optionally implement `Namer` to provide a custom form name:

```go
type Namer interface {
    FormName() string
}
```

If not implemented, defaults to `"form"`.
