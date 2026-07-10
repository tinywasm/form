# Input Types

This package contains all input implementations for `tinywasm/form`.
Each input implements `model.Kind` (Type, Validate, Clone) plus metadata getters.
Inputs are **render-free**: they carry no `dom`/`html` imports. Rendering is done by the
`form` package via `RenderInput(input.Input)`. All inputs use **only** `tinywasm/fmt` —
no `errors` or `strconv` from the standard library.

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
// Schema definition (ormc generates this).
// input.* kinds get a form input + validation; model.* base kinds
// (e.g. model.Text()) validate only and are never rendered.
var schema = []model.Field{
    {Name: "email", Type: input.Email(), NotNull: true},
}

// form.New calls Clone internally, for fields whose Type implements input.Input:
// field.Type.(input.Input).Clone(formID, fieldName) → positioned input with id, name, HTML attributes
```

## Creating a Custom Input (embedding Base)

All inputs share the same pattern: embed `Base`, configure `Permitted` rules, implement the `Input` interface. Custom inputs can live in your own package — `Base`, `InitBase` and all setters are exported.

```go
package myapp

import (
    "github.com/tinywasm/fmt"
    "github.com/tinywasm/form/input"
)

// myInput is a custom input that only allows lowercase hex characters.
type myInput struct {
    input.Base
}

// MyInput creates a prototype — no arguments.
func MyInput() input.Input {
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
    return m.Permitted.Validate(m.FieldName(), value)
}

// Clone creates a positioned copy preserving all configuration.
func (m *myInput) Clone(parentID, name string) input.Input {
    c := *m
    c.InitBase(parentID, name, "text")
    return &c
}
```

The custom input renders with the generic markup for its `htmlName`
(`"text"` above). Storage defaults by `htmlName` (`number` → int,
`checkbox` → bool, else text); override `Storage()` on your struct if your
kind needs a different mapping.

### Custom markup (Renderer)

If your input needs a special UI (e.g. a color picker, a file uploader, or a
composite widget), implement the `form.Renderer` interface in your struct.
The `form` package will call this instead of its default switch:

```go
import (
    "github.com/tinywasm/dom"
    "github.com/tinywasm/form"
)

func (m *myInput) RenderInput(value *dom.SignalString, onInput func(string)) *dom.Element {
    el := dom.NewElement("div").Class("my-custom-widget")
    // ... build your UI, bind 'value' signal for updates
    // ... call 'onInput(newValue)' when the user interacts
    return el
}
```

The `form.Renderer` interface lives in the root `form` package because it
references `dom.Element`; the `input` package remains dom-free by contract.

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
