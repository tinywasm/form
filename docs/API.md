# API Specification

Public API for `tinywasm/form`.

## Global Functions

### `New(parentID string, structPtr any) (*Form, error)`
Creates a form from a struct. Returns error if any exported field has no matching input.
- **parentID**: ID of the parent element in the DOM.
- **action**: Defaults to `/structname` (lowercase struct name).
- **method**: Defaults to `POST`.

### `RegisterInput(inputs ...input.Input)`
Registers input types. Matching is based on `HTMLName()` and aliases.

### `SetGlobalClass(classes ...string)`
Sets default CSS classes for all forms.

### `SetSSR(enabled bool)`
Toggles Server-Side Rendering mode globally for the library (via `registry.go`).

## Types

### `Form`
Core struct managing fields and submission logic. See [form.go](../form.go).

**Methods:**
- `ID() string` - Returns the generated form ID.
- `RenderHTML() string` - Generates HTML (respects SSR mode).
- `Validate() error` - Validates all inputs.
- `OnMount()` - Binds centralized event listeners (WASM only). See [mount.go](../mount.go).
- `OnSubmit(func(any) error)` - Sets a callback for successful form submission in WASM.
- `SyncValues() error` - Synchronizes input values back to the underlying struct.

### `input.Input` Interface 
Minimal interface for form components. See [interface.go](../input/interface.go).

```go
type Input interface {
    dom.Component
    HTMLName() string
    FieldName() string
    ValidateField(value string) error
    Clone(parentID, name string) Input
}
```

### `input.Base`
Common state embedded in all inputs. See [base.go](../input/base.go).
- `Values []string`: Current value(s).
- `Options []fmt.KeyValue`: Available options (for select/radio/etc).
- `HTMLName()`: Access to the HTML type.
- `FieldName()`: Access to the struct field name.
