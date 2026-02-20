# Design & Architecture

## Philosophy
- **TinyGo first**: No `reflect` on hot paths, flat slices over maps, minimal allocations.
- **Convention over configuration**: Field name → input type via aliases. Tags for fine-tuning.
- **Dual mode**: Same struct renders via WASM event delegation or full SSR.

## Core Layers

### 1. Registry (`registry.go`)
Global `registeredInputs []input.Input` slice. `New()` iterates it to find a match per field:
1. `lowercase(FieldName)` matches `input.HTMLName()` or any alias
2. `lowercase(StructName.FieldName)` matches any alias

`Clone(parentID, name)` is called on the matched template to produce a unique instance.

### 2. State & Binding (`form.go`)
- **One-way on creation**: `New()` copies struct values → inputs via `SetValues()`
- **Two-way on submit**: `SyncValues()` copies input values → struct fields

Supports: `string`, `[]string`, any type via `fmt.Convert(...).String()`.

### 3. Validation (`validate.go`)
`Form.Validate()` iterates `f.Inputs`, calls each `inp.ValidateField(GetSelectedValue())`.
Fields tagged `validate:"false"` are skipped.

`input.Permitted` provides whitelist validation (letters, numbers, special chars, min/max length).
Each input embeds `Permitted` and may add custom logic in `ValidateField`.

### 4. Rendering (`render.go`)
`RenderHTML()` builds: `<form id="..."> [method action if SSR] inputs [submit if SSR] </form>`

### 5. WASM Interactivity (`mount.go`)
`OnMount()` (build tag: `wasm`) attaches **one** delegated listener at the `<form>` element:
- `input`/`change` → `SetValues()` + `ValidateField()` per matching input
- `submit` → `PreventDefault()` → `SyncValues()` → `Validate()` → `OnSubmit` callback
