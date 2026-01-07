# Implementation Plan

Central document for `tinywasm/form`.

## Goal
Minimalist, WASM-optimized form library using `tinywasm/dom`.

**Principles:**
1. Minimal API: `New`, `RegisterInput`.
2. Global Registry: Forms auto-register.
3. Struct-Based: Forms defined by Go structs.
4. No Tags: Config via field names and global settings.
5. Performance: Slices over maps.
6. **Dual Render Mode**: WASM/fetch vs SSR.

## Documentation
- [API Specification](API.md)
- [Design & Architecture](DESIGN.md)
- [Interactivity Strategy](INTERACTIVITY_AND_MOUNTING.md)
- [Input Types](../input/README.md)

## Architecture

### Render Modes

| Mode | SetModeSSR | Renders |
|------|------------|---------|
| **WASM/Fetch** | `false` (default) | id, class, inputs only |
| **SSR** | `true` | + method, action, submit |

Default endpoint: `/structname` (lowercase).

### Form Creation
```go
f, err := form.New("form-id", &MyStruct{})
```
1. Reflects on struct fields.
2. Matches each field to registered inputs via `Matches()`.
3. Returns error if no match found.
4. Sets `action` to `/structname`.

### Input Interface
```go
type Input interface {
    dom.Component
    HtmlName() string
    ValidateField(value string) error
    Clone(parentID, name string) Input
}
```

## Event Handling
**Centralized listener** at Form level. See [INTERACTIVITY_AND_MOUNTING.md](INTERACTIVITY_AND_MOUNTING.md).

## Validation
Uses `Permitted` struct. `Form.Validate()` iterates all inputs.
