# SKILL — tinywasm/form

Compact reference for LLMs. Contains everything needed to use and extend this library.

## Rules (Non-Negotiable)
- **No stdlib imports**: Only `github.com/tinywasm/fmt`. Never `errors`, `strconv`, `strings`.
- **No maps in inputs**: Use slices. Maps increase WASM binary size.
- **`fmt.Err("Noun", "Adjective")` for errors**: e.g. `fmt.Err("Format", "Invalid")`.
- **`fmt.Convert(s).Int()`** instead of `strconv.Atoi`. **`fmt.Contains(s, sub)`** instead of `strings.Contains`.
- **`_ "github.com/tinywasm/fmt/dictionary"`** must only be imported in `_test.go` files.

## Core Concept

`form.New("parentID", &MyStruct{})` reflects exported fields → matches each to a registered `Input` by name/alias → builds a `*Form` ready to render and validate.

```go
type Order struct {
    Name  string `placeholder:"Customer name"`
    Email string
    Total string `options:"s:Small,m:Medium,l:Large" validate:"false"`
}

f, err := form.New("content", &Order{Name: "Ana"})
// err != nil → some field has no matching input
html := f.RenderHTML()            // <form id="content.order">...</form>
f.SetSSR(true).RenderHTML()       // + method="POST" action="/order" + <button>
f.Validate()                      // validates all inputs (skips validate:"false")
f.SyncValues()                    // copies input values back to struct
f.OnSubmit(func(v any) error { ... })
```

## Field Matching (in `New`)

Order per field:
1. `lowercase(FieldName)` vs Input's `htmlName` + `aliases`
2. `lowercase(StructName.FieldName)` vs Input's `aliases`

## Struct Tags

| Tag | Example | Effect |
|-----|---------|--------|
| `placeholder` | `placeholder:"Enter name"` | Sets HTML placeholder |
| `title` | `title:"Tooltip text"` | Sets HTML title (tooltip) |
| `options` | `options:"k1:Label1,k2:Label2"` | Sets options for select/radio/datalist |
| `validate` | `validate:"false"` | Skips `ValidateField` for this field |

## Form Methods

```go
f.GetID()                          // "parentID.structname"
f.SetSSR(true)                     // enable SSR render mode
f.RenderHTML()                     // generate HTML string
f.Validate()                       // error on first invalid input
f.SyncValues()                     // struct ← inputs
f.Input("FieldName")               // get input by field name
f.SetOptions("Field", opts...)     // set options programmatically
f.SetValues("Field", "val")        // set value programmatically
f.OnSubmit(fn)                     // WASM: callback on submit
f.OnMount()                        // WASM: called by dom after render
```

## Input Interface

```go
type Input interface {
    dom.Component  // GetID(), SetID(), RenderHTML(), Children()
    HTMLName() string
    FieldName() string
    ValidateField(value string) error
    Clone(parentID, name string) Input
}
```

## All Registered Inputs (alphabetical)

| Constructor | htmlType | Aliases | Validation |
|-------------|----------|---------|-----------|
| `Address(p,n)` | `text` | `address`, `addr` | Letters+Numbers+`. , - # /`, 5–200 |
| `Checkbox(p,n)` | `checkbox` | `check`, `boolean`, `bool` | `true/false/on/1/0` or empty |
| `Date(p,n)` | `date` | `fecha` | `YYYY-MM-DD`, leap year, day range |
| `Datalist(p,n)` | `text` | `list`, `options` | Value must match an `Options.Key` |
| `Email(p,n)` | `email` | `mail`, `correo` | Letters+Numbers+`@._-`, 5–100 |
| `Filepath(p,n)` | `text` | `path`, `dir`, `file` | Letters+Numbers+`./\-_`, no space, 1–200 |
| `Gender(p,n)` | `radio` | `gender`, `sex` | pre-wired `m`/`f` options |
| `Hour(p,n)` | `time` | `hour` | `HH:MM`, digits+`:`, 24:xx rejected |
| `IP(p,n)` | `text` | `ip`, `address` | IPv4 or IPv6, `0.0.0.0` rejected |
| `Number(p,n)` | `number` | `num`, `amount`, `price`, `age` | Digits only, 1–20 chars |
| `Password(p,n)` | `password` | `pass`, `clave`, `pwd` | Any char, 5–50 |
| `Phone(p,n)` | `tel` | `phone`, `mobile`, `cell` | Digits+`+()-`, 7–15 |
| `Radio(p,n)` | `radio` | — | Value must match an `Options.Key` |
| `Rut(p,n)` | `text` | `rut`, `run`, `dni` | `NNNNNNN-D`, verifies check digit |
| `Select(p,n)` | `select` | — | Value must match an `Options.Key` |
| `Text(p,n)` | `text` | `name`, `fullname`, `username` | Letters+Numbers+`. ,()`, 2–100 |
| `Textarea(p,n)` | `textarea` | `description`, `details`, `comments` | Wide charset+`\n`, 5–2000 |

## Creating a Custom Input

```go
package input

import "github.com/tinywasm/fmt"

type myInput struct {
    Base
    Permitted Permitted
}

func MyInput(parentID, name string) Input {
    m := &myInput{
        Permitted: Permitted{
            Letters: true, Numbers: true,
            Minimum: 2, Maximum: 50,
        },
    }
    m.Base.InitBase(parentID, name, "text", "alias1", "alias2")
    return m
}

func (m *myInput) HTMLName() string                    { return m.Base.HTMLName() }
func (m *myInput) ValidateField(v string) error        { return m.Permitted.Validate(v) }
func (m *myInput) RenderHTML() string                  { return m.Base.RenderInput() }
func (m *myInput) Clone(p, n string) Input             { return MyInput(p, n) }
```

Then register in `registry.go` `init()`:
```go
input.MyInput("", ""),
```

## `Permitted` Fields

```go
Permitted{
    Letters: true,      // A-Z a-z Ñ ñ
    Tilde: true,        // Á É Í Ó Ú á é í ó ú
    Numbers: true,      // 0-9
    WhiteSpaces: true,  // ' '
    BreakLine: true,    // '\n'
    Characters: []rune{'@', '.'},  // extra chars
    Minimum: 5,         // min length (0 = no min)
    Maximum: 100,       // max length (0 = no max)
    ExtraValidation: func(s string) error { ... }, // custom logic
}
```

`Permitted.Validate(s)` errors: `"minimum N chars"`, `"maximum N chars"`, `"space not allowed"`, `"character X not allowed"`.

## WASM Event Flow

```
dom.Mount("root", f)
  → f.RenderHTML() injected
  → f.OnMount():
      el.On("input",  fn) → SetValues + ValidateField per input
      el.On("change", fn) → same
      el.On("submit", fn) → PreventDefault → SyncValues → Validate → onSubmit(f.Value)
```

## File Map

| File | Owns |
|------|------|
| `form.go` | `Form` struct, `New()`, tag parsing, `SyncValues`, `Input`, `SetOptions`, `SetValues` |
| `registry.go` | `registeredInputs`, `RegisterInput`, `findInputForField`, `SetGlobalClass` |
| `render.go` | `RenderHTML`, `SetSSR` |
| `validate.go` | `Validate` |
| `mount.go` | `OnMount`, `OnUnmount` (wasm only) |
| `input/base.go` | `Base` struct |
| `input/permitted.go` | `Permitted` validation engine |
| `input/interface.go` | `Input` interface |
| `input/inputs_test.go` | Table-driven tests for all inputs |
