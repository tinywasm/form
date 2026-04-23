# Input Types

This package contains all input implementations for `tinywasm/form`.
Each input implements `fmt.Widget` (Type, Validate, Clone) and `dom.Component` (RenderHTML).
All inputs use **only** `tinywasm/fmt` — no `errors` or `strconv` from the standard library.

## Available Inputs

| Type | HTML type | Validation rules |
|------|-----------|-----------------|
| `Address` | `text` | Letters, Numbers, `. , - # /`, Min: 5, Max: 200 |
| `Checkbox` | `checkbox` | `true`, `false`, `on`, `1`, `0` or empty |
| `Date` | `date` | `YYYY-MM-DD` format, leap year + month/day range check |
| `Datalist` | `text` | Value must match one of the registered `Options.Key` |
| `Email` | `email` | Letters, Numbers, `@ . _ -`, Min: 5, Max: 100 |
| `Filepath` | `text` | Letters, digits, `.\/- _`, no whitespace, Min: 1 |
| `Gender` | `radio` | `m`/`f` pre-wired options |
| `Hour` | `time` | `HH:MM` format, digits + `:`, validates 24h range |
| `IP` | `text` | IPv4 or IPv6 format; `0.0.0.0` rejected |
| `Number` | `number` | Digits only (0-9), Min: 1, Max: 20 chars |
| `Password` | `password` | Any char, Min: 5, Max: 50 |
| `Phone` | `tel` | Digits, `+ ( ) -`, Min: 7, Max: 15 |
| `Radio` | `radio` | Value must match one of the registered `Options.Key` |
| `Rut` | `text` | Chilean RUT `XXXXXXX-D`, verifies check digit |
| `Search` | `search` | Letters, Numbers, Spaces, Min: 0, Max: 100 — optional |
| `Select` | `select` | Value must match one of the registered `Options.Key` |
| `Text` | `text` | Letters, Numbers, `. , ( )`, Min: 2, Max: 100 |
| `Textarea` | `textarea` | Wide char set incl. `\n`, Min: 5, Max: 2000 |

## No Standard Library

> **Rule**: All input files must import only `github.com/tinywasm/fmt`. No `errors`, `strconv`, or `strings`.

Use the `tinywasm/fmt` equivalents:

```go
// Instead of strconv.Atoi:
val, err := fmt.Convert("42").Int()

// Instead of errors.New:
return fmt.Err("Field", "Invalid")

// Instead of strings.ToLower:
lower := fmt.Convert(s).ToLower().String()

// Instead of strings.Contains:
found := fmt.Contains(haystack, needle)
```

## Prototype Pattern

Constructors take **zero arguments** and return stateless prototypes. The form layer calls `Clone(parentID, name)` to create positioned instances that preserve all configuration.

```go
// Schema definition (ormc generates this)
var schema = []fmt.Field{
    {Name: "email", Type: fmt.FieldText, Widget: input.Email()},
}

// form.New calls Clone internally:
// field.Widget.Clone(formID, fieldName) → positioned input with id, name, HTML attributes
```

## Creating a Custom Input (embedding Base)

All inputs share the same pattern: embed `Base`, configure `Permitted` rules, implement the `Input` interface.

```go
package input

import "github.com/tinywasm/fmt"

// myInput is a custom input that only allows lowercase hex characters.
type myInput struct {
    Base
}

// MyInput creates a prototype — no arguments.
func MyInput() Input {
    m := &myInput{}
    m.Letters = true
    m.Numbers = true
    m.Minimum = 1
    m.Maximum = 40
    m.InitBase("", "", "text")
    m.SetPlaceholder("e.g. 3f4a1b")
    m.SetTitle("Lowercase hex only")
    return m
}

func (m *myInput) Validate(value string) error {
    // Custom rule: no uppercase letters
    for _, c := range value {
        if c >= 'A' && c <= 'F' {
            return fmt.Err("Character", "Invalid")
        }
    }
    return m.Permitted.Validate(m.name, value)
}

// Clone creates a positioned copy preserving all configuration.
func (m *myInput) Clone(parentID, name string) fmt.Widget {
    c := *m
    c.InitBase(parentID, name, "text")
    return &c
}
```

### Base Available Methods

| Method | Purpose |
|--------|---------|
| `InitBase(parentID, name, htmlName)` | **Required** — sets ID, name, and HTML type |
| `SetPlaceholder(string)` | HTML placeholder text |
| `SetTitle(string)` | HTML title (tooltip) |
| `SetOptions(...fmt.KeyValue)` | Options for select/radio/datalist |
| `AddAttribute(key, value string)` | Custom extra HTML attributes |
| `SetRequired(bool)` | HTML required attribute |
| `SetSkipValidation(bool)` | Skip validation entirely |
| `RenderInput()` | Generates `<input>`, `<textarea>`, or `<select>` tag |

## Composition Pattern (wrapping another input)

Reuse existing inputs to create semantic wrappers:

```go
func Gender() Input {
    g := &gender{}
    g.InitBase("", "", "radio")
    g.SetOptions(
        fmt.KeyValue{Key: "m", Value: "Male"},
        fmt.KeyValue{Key: "f", Value: "Female"},
    )
    return g
}
```
