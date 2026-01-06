# API Specification

This document defines the public API for the `tinywasm/form` library.

## Global Functions

## Global Functions

### `New[T any](structInstance T) *Form`
Creates a new form instance from a struct and registers it in the internal global handler.
*   **Params**: A struct instance (value or pointer).
*   **Returns**: A pointer to the created `Form`.
*   **Side Effect**: Registers the form globally.
*   **Example**:
    ```go
    type Login struct {
        User     string
        Password string
    }
    f := form.New(Login{})
    ```

### `SetGlobalClass(classes ...string)`
Configures the default CSS classes applied to all form elements globally.
*   **Semantics**: Appends/Updates the global class configuration.
*   **Example**: `form.SetGlobalClass("my-input-class", "mb-3")`

### `SetType(name string, config InputConfig)` (Proposed)
Registers or overrides a "Smart Type" definition.
*   **Example**: 
    ```go
    form.SetType("ZipCode", InputConfig{
        Type: "number",
        Validation: func(v string) error { ... }
    })
    ```

## Types

### `Form`
Represents a form instance.
```go
type Form struct {
    // Hidden internal state
}
```

#### `(*Form) Render() string`
Generates the HTML string for the form.

#### `(*Form) Bind()`
Binds the form to the DOM (attaches event listeners).

### `InputConfig`
Configuration for an input type.
```go
type InputConfig struct {
    Type            string // "text", "password", "email", etc.
    Class           string // CSS class
    Placeholder     string // Default placeholder
    Validation      func(string) error
}
```
