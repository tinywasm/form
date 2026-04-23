# PLAN: Add `Search` Input to `tinywasm/form/input`

## Context

Module: `github.com/tinywasm/form`
Sub-package: `input` — located at `tinywasm/form/input/`
Go version: 1.25.2
Final WASM compiler: TinyGo

This module provides form input widgets used by `tinywasm/components` and the `tinywasm/form` form layer.
Each input type is a standalone file implementing the `Input` interface.

## Objective

Add a new `Search()` input factory to `tinywasm/form/input/` for use as a search box inside UI components (not as a validated form field). It renders `<input type="search">`.

## Constraints

- **No standard library imports**: use only `github.com/tinywasm/fmt`. No `errors`, `strconv`, or `strings`.
- All input files follow the same pattern: embed `Base`, configure `Permitted` rules, implement `Clone(parentID, name string) fmt.Widget`.
- Constructors take **zero arguments** and return a prototype (stateless). The form layer calls `Clone(parentID, name)` to produce positioned instances.
- The `Input` interface is defined in `tinywasm/form/input/interface.go` — do not modify it.
- Do not modify `Base` in `base.go`.
- `RenderInput()` in `base.go` already handles `type="search"` correctly when `htmlName` is set to `"search"` — no custom `RenderHTML()` needed.

## Existing Pattern to Follow

File: `tinywasm/form/input/text.go`

```go
package input

import "github.com/tinywasm/fmt"

type text struct{ Base }

func Text() Input {
    t := &text{}
    t.Letters = true
    t.Tilde = true
    t.Numbers = true
    t.Spaces = true
    t.Extra = []rune{'.', ',', '(', ')'}
    t.Minimum = 2
    t.Maximum = 100
    t.InitBase("", "", "text")
    return t
}

func (t *text) Clone(parentID, name string) fmt.Widget {
    c := *t
    c.InitBase(parentID, name, "text")
    return &c
}
```

## Task

Create the file `tinywasm/form/input/search.go` with the following spec:

### Type name
`search_` (unexported, matches the pattern of other types: `text`, `select_`, etc.)

### Constructor `Search() Input`
- `Letters = true`
- `Numbers = true`
- `Spaces = true`
- `Minimum = 0` (search boxes are optional by default)
- `Maximum = 100`
- `htmlName = "search"` (passed to `InitBase("", "", "search")`)
- No `Required` by default

### `Clone(parentID, name string) fmt.Widget`
Same pattern as all other inputs: copy the struct value, call `InitBase(parentID, name, "search")`, return pointer.

### No custom `Validate()`
The default `Base.Validate()` is sufficient — search inputs are optional.

## File to Create

`tinywasm/form/input/search.go`

```go
package input

import "github.com/tinywasm/fmt"

type search_ struct{ Base }

func Search() Input {
    s := &search_{}
    s.Letters = true
    s.Numbers = true
    s.Spaces = true
    s.Minimum = 0
    s.Maximum = 100
    s.InitBase("", "", "search")
    return s
}

func (s *search_) Clone(parentID, name string) fmt.Widget {
    c := *s
    c.InitBase(parentID, name, "search")
    return &c
}
```

## Test to Add

Add a test case to `tinywasm/form/input/inputs_test.go` (or `render_test.go` — follow the existing file's pattern).

The test must verify:
1. `Search()` returns an `Input` (not nil).
2. `Search().Type()` returns `"search"`.
3. `Search().Clone("myform", "q").RenderHTML()` contains `type="search"` and `id="myform.q"`.
4. `Search().Validate("")` returns `nil` (optional field).

## README Update

Add a row to the table in `tinywasm/form/input/README.md`:

| `Search` | `search` | Letters, Numbers, Spaces, Min: 0, Max: 100 — optional |

Insert it alphabetically between `Rut` and `Select`.
