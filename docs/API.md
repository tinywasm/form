# API Reference

See `README.md` for the consolidated API. This file contains additional detail.

## `form.New` — Widget Resolution Detail

```go
f, err := form.New("content", data) // data implements fmt.Fielder
// -> f.GetID() == "content." + resolveStructName(data)
// -> f.RenderHTML() renders all fields that have a Widget in data.Schema()
```

For each field in `data.Schema()`:
1. `field.Widget == nil` → skip (no UI binding).
2. `field.Widget.Clone(formID, fieldName).(input.Input)` → positioned input.
3. `field.PK && field.AutoInc` → skip (auto-increment PKs not editable).
4. `field.NotNull` → `SetRequired(true)` on the input.
5. Current value bound via `fmt.ReadValues()` + `SetValues()`.

## `(*Form).Validate()` — Validation Detail

- Skips fields with `SkipValidation` set to true in the input.
- Uses `GetSelectedValue()` to get current value per input.
- Calls `inp.Validate(val)` (promoted from `fmt.Widget`).
- Returns the **first** error encountered.

## `(*Form).SyncValues(data fmt.Fielder)` — Binding Detail

Synchronizes input values back to the struct pointers provided by `data.Pointers()`.
Supports `fmt.FieldText`, `fmt.FieldInt`, `fmt.FieldFloat`, and `fmt.FieldBool`.

## `(*Form).ValidateData(action byte, data fmt.Fielder)` — Server-side Validation

Validates the provided `data` using the form's input rules. Satisfies `crudp.DataValidator`.

## `input.Permitted` — Validation Engine

```go
type Permitted struct {
    Letters         bool     // A-Z, a-z
    Tilde           bool     // accented chars
    Numbers         bool     // 0-9
    WhiteSpaces     bool     // space ' '
    BreakLine       bool     // '\n'
    Tabulation      bool     // '\t'
    Characters      []rune   // extra allowed chars e.g. []rune{'@', '.'}
    TextNotAllowed  []string // blacklisted substrings
    Minimum         int      // minimum length (0 = no limit)
    Maximum         int      // maximum length (0 = no limit)
    ExtraValidation func(string) error // custom logic
    StartWith       *Permitted // rules for first character only
}
```

Error messages from `Permitted.Validate()`:
- `"minimum N chars"` — value shorter than Minimum
- `"maximum N chars"` — value longer than Maximum
- `"space not allowed"` — space when WhiteSpaces=false
- `"character X not allowed"` — disallowed character

## Namer Interface

Fielder types can optionally implement `Namer` to provide a custom form name:

```go
type Namer interface {
    FormName() string
}
```

If not implemented, defaults to `"form"`.
