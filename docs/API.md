# API Specification

Public API for `tinywasm/form`.

## Global Functions

### `New(id string, structPtr any) (*Form, error)`
Creates a form from a struct. Returns error if any exported field has no matching input.
- **action**: Defaults to `/structname` (lowercase struct name).
- **method**: Defaults to `POST`.

```go
type Login struct {
    Username string
    Password string
}
f, err := form.New("login-form", &Login{})
```

### `RegisterInput(inputs ...input.Input)`
Registers input types. Uses `HtmlName()` and aliases for field matching.

```go
form.RegisterInput(
    input.Text("", ""),
    input.Email("", ""),
    input.Password("", ""),
)
```

### `SetGlobalClass(classes ...string)`
Sets default CSS classes for all forms.

### `SetModeSSR(enabled bool)`
Toggles Server-Side Rendering mode.
- **false (default)**: Minimal HTML for WASM/fetch usage.
- **true**: Full HTML with `method`, `action`, submit button for standard forms.

## Types

### `Form`
```go
type Form struct {
    ID     string
    Inputs []input.Input
}
```

**Methods:**
- `RenderHTML() string` - Generates HTML (respects SSR mode).
- `Validate() error` - Validates all inputs.
- `OnMount()` - Binds centralized event listener (WASM only).

### `input.Input` Interface
```go
type Input interface {
    dom.Component
    HtmlName() string
    ValidateField(value string) error
    Clone(parentID, name string) Input
}
```

### `input.Base` Fields
```go
type Base struct {
    Value       string
    Placeholder string
    Title       string
    Required    bool
    Disabled    bool
    Readonly    bool
}
```
