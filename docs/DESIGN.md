# Design & Architecture

## Philosophy
- **TinyGo first**: Zero `reflect`, flat slices over maps, minimal allocations.
- **Convention over configuration**: Field name → input type via aliases. `ormc` generated schema for fine-tuning.
- **Interface-driven**: Decoupled from concrete structs via `fmt.Fielder`.

## Core Layers

### 1. Registry (`registry.go`)
Global `registeredInputs []input.Input` slice. `New()` iterates it to find a match per field:
1. `lowercase(FieldName)` matches `input.HTMLName()` or any alias
2. `lowercase(StructName.FieldName)` matches any alias
3. Explicit override via `fmt.Field.Input`

### 2. State & Binding (`form.go`, `sync.go`)
- **Fielder-based**: `New()` uses `data.Schema()` and `data.Values()`.
- **Pointer-based sync**: `SyncValues(data)` uses `data.Pointers()` to write back values without reflection.

### 3. Validation (`validate.go`, `validate_struct.go`)
- `Form.Validate()` iterates `f.Inputs`, calls each `inp.ValidateField(GetSelectedValue())`.
- `ValidateData(action, data)` provides server-side or isomorphic validation.

### 4. Rendering (`render.go`)
`RenderHTML()` builds: `<form id="..."> [method action if SSR] inputs [submit if SSR] </form>`

### 5. WASM Interactivity (`mount.go`)
`OnMount()` attaches delegated listeners:
- `input`/`change` → Live sync and validation.
- `submit` → `SyncValues(f.data)` → `Validate()` → `OnSubmit(f.data)`.
