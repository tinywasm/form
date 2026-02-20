# Input Types

This package contains all input implementations for `tinywasm/form`.
Each input implements `ValidateField(value string) error` and `Clone(parentID, name string) Input`.
All inputs use **only** `tinywasm/fmt` — no `errors` or `strconv` from the standard library.

## Available Inputs

| Type | HTML type | Aliases | Validation rules |
|------|-----------|---------|-----------------|
| `Address` | `text` | `address`, `addr` | Letters, Numbers, `. , - # /`, Min: 5, Max: 200 |
| `Checkbox` | `checkbox` | `check`, `boolean`, `bool` | `true`, `false`, `on`, `1`, `0` or empty |
| `Date` | `date` | `fecha` | `YYYY-MM-DD` format, leap year + month/day range check |
| `Datalist` | `text` | `list`, `options` | Value must match one of the registered `Options.Key` |
| `Email` | `email` | `mail`, `correo` | Letters, Numbers, `@ . _ -`, Min: 5, Max: 100 |
| `Filepath` | `text` | `path`, `dir`, `file` | Letters, digits, `.\/- _`, no whitespace, Min: 1 |
| `Gender` | `radio` | `gender`, `sex` | `m`/`f` pre-wired options |
| `Hour` | `time` | `hour` | `HH:MM` format, digits + `:`, validates 24 not valid |
| `IP` | `text` | `ip`, `address` | IPv4 or IPv6 format; `0.0.0.0` rejected |
| `Number` | `number` | `num`, `amount`, `price`, `age` | Digits only (0-9), Min: 1, Max: 20 chars |
| `Password` | `password` | `pass`, `clave`, `pwd` | Any char, Min: 5, Max: 50 |
| `Phone` | `tel` | `phone`, `mobile`, `cell` | Digits, `+ ( ) -`, Min: 7, Max: 15 |
| `Radio` | `radio` | — | Value must match one of the registered `Options.Key` |
| `Rut` | `text` | `rut`, `run`, `dni` | Chilean RUT `XXXXXXX-D`, verifies check digit |
| `Select` | `select` | `select` | Value must match one of the registered `Options.Key` |
| `Text` | `text` | `name`, `fullname`, `username` | Letters, Numbers, `. , ( )`, Min: 2, Max: 100 |
| `Textarea` | `textarea` | `description`, `details`, `comments` | Wide char set incl. `\n`, Min: 5, Max: 2000 |

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

## Creating a Custom Input (embedding Base)

All inputs share the same pattern: embed `Base`, add a `Permitted` struct for rules, implement the `Input` interface.

```go
package input

import "github.com/tinywasm/fmt"

// myInput is a custom input that only allows lowercase hex characters.
type myInput struct {
    Base
    Permitted Permitted
}

// MyInput creates a new instance ready for use.
func MyInput(parentID, name string) Input {
    m := &myInput{
        Permitted: Permitted{
            Letters:    true,
            Numbers:    true,
            Characters: []rune{},
            Minimum:    1,
            Maximum:    40,
        },
    }
    // Args: parentID, fieldName, htmlType, ...aliases
    m.Base.InitBase(parentID, name, "text", "myhex", "hex")
    m.Base.SetPlaceholder("e.g. 3f4a1b")
    m.Base.SetTitle("Lowercase hex only")
    return m
}

func (m *myInput) HTMLName() string { return m.Base.HTMLName() }

func (m *myInput) ValidateField(value string) error {
    // Custom rule: no uppercase letters
    for _, c := range value {
        if c >= 'A' && c <= 'F' {
            return fmt.Err("Character", "Invalid")
        }
    }
    return m.Permitted.Validate(value)
}

func (m *myInput) RenderHTML() string  { return m.Base.RenderInput() }

func (m *myInput) Clone(parentID, name string) Input {
    return MyInput(parentID, name)
}
```

### Base Available Methods

| Method | Purpose |
|--------|---------|
| `InitBase(parentID, name, htmlName, aliases...)` | **Required** — sets ID, name, html type and aliases |
| `SetPlaceholder(string)` | HTML placeholder text |
| `SetTitle(string)` | HTML title (tooltip) |
| `SetOptions(...fmt.KeyValue)` | Options for select/radio/datalist |
| `AddAttribute(key, value string)` | Custom extra HTML attributes |
| `SetSkipValidation(bool)` | Skip validation entirely |
| `RenderInput()` | Generates `<input>`, `<textarea>`, or `<select>` tag |

## Composition Pattern (wrapping another input)

Reuse existing inputs to create semantic wrappers:

```go
func Gender(parentID, name string) Input {
    g := &gender{}
    g.Base.InitBase(parentID, name, "radio", "gender", "sex")
    g.Base.SetOptions(
        fmt.KeyValue{Key: "m", Value: "Male"},
        fmt.KeyValue{Key: "f", Value: "Female"},
    )
    return g
}
```

## Matching Logic

The form engine matches struct fields to inputs in this order:
1. `LowerCase(FieldName)` vs `htmlName` or `aliases`
2. `LowerCase(StructName.FieldName)` vs `aliases`

## Registering Custom Inputs

Register your own input type to make it available to `form.New`:

```go
form.RegisterInput(MyCustomInput("", ""))
```
