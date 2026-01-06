# IMPLEMENTATION.md

This document serves as the central hub for the implementation of the `tinywasm/form` library. It outlines the goals, architectural decisions, and links to detailed specifications.

**Status**: DRAFT / RFC (Request for Comments)

## Goal

Create a minimalist, intuitive, and WASM-optimized (TinyGo) library for building web forms. It builds upon `tinywasm/dom` for DOM interaction and aims to solve strict memory/size constraints by avoiding heavy runtime overhead (e.g., minimizing map usage).

**Key Principles:**
1.  **Minimal API**: `New`, `Create`, `Add`.
2.  **Global Registry**: Forms are automatically registered upon creation.
3.  **Struct-Based**: Forms are defined by Go structs.
4.  **Tag-Less Configuration**: Render behavior is driven by field names/types and global configuration, *not* struct tags.
5.  **Performance**: Use slices instead of maps for storage to optimize for TinyGo.

## Documentation Map

*   **[API Specification](API.md)**: Detailed definition of public functions and interfaces.
*   **[Design & Architecture](DESIGN.md)**: Deep dive into internal decisions (Global registry, Reflection usage, etc.).

## 1. Core Architecture

### Global Handler
A global `formHandler` instance (singleton) will manage all form instances.
*   **Storage**: `[]*Form` (Slice, not Map).
*   **Lookup**: Linear search by Form Name/ID (acceptable for the number of forms usually present in a typical app).

### Form Creation (`New`)
*   **Signature**: `func New(formStruct any) *Form`
*   **Behavior**:
    1.  Reflects on `formStruct`.
    2.  Generates a unique ID (or uses struct name).
    3.  Registers the form in the global handler.
    4.  Returns the `*Form` instance for method chaining (e.g., `.Render()`).

### Field Rendering
*   **Discovery**: Uses `tinywasm/fmt.HasUpperPrefix` to identify exported fields.
*   **Mapping**:
    *   Fields are mapped to "Input Types" based on:
        1.  **Global Configuration**: Register custom types via `Set`.
        2.  **Naming Convention**: e.g., field `Email` -> `<input type="email">`, `Password` -> `<input type="password">`.
        3.  **Go Type**: `int` -> `<input type="number">`, `bool` -> `<input type="checkbox">`.
*   **No Tags**: We avoid `ctx:"ui"` tags. Configuration is decoupled from the struct definition.

## 2. Validation Architecture
The validation logic is a port of the robust `Permitted` system. 

### `Permitted` Struct
Defines the rules for a field.
[View Code](file:///home/cesar/Dev/Pkg/tinywasm/form/input/permitted.go)

### `Input` Interface
The polymorphic contract for all inputs.
[View Code](file:///home/cesar/Dev/Pkg/tinywasm/form/input/interface.go)
```go
type Input interface {
    dom.Component // Embeds ID() and RenderHTML()
    HtmlName() string
    ValidateField(value string) error
}
```

### Concrete Types
*   **Text**: Standard text input.
    *   **Factory**: `input.Text(formID, name)`
    *   **Struct**: Private `text`. Uses `fmt.Html` for rendering.
    [View Code](file:///home/cesar/Dev/Pkg/tinywasm/form/input/text.go)

## 4. Interactivity & Mounting
[Read Strategy Document](file:///home/cesar/Dev/Pkg/tinywasm/form/docs/INTERACTIVITY_AND_MOUNTING.md)
*   **Form**: Implements `dom.Mountable` and delegates `OnMount` to children.
*   **Inputs**: Implement `dom.Mountable` to bind events (validation on `input`).

*   **Internal Defaults**: The library comes pre-loaded with standard types (`Name`, `Email`, `RUT`, etc.) registered in `init()`.
*   **Extensibility**: Users can register new types via `form.RegisterType()`.
*   **Dynamic**: Global configuration allows application-wide settings (e.g. changing the underlying regex for "Phone" globally).

## 4. Derived Requirements (from Reference)
*   **HTML Generation**: Pure generation without DB coupling.
*   **Validation**:
    *   **Whitelist Approach**: Explicitly allow characters (`Letters`, `Numbers`) rather than blacklisting.
    *   **Validation.go Port**: The logic from `Archive/mono/validation.go` (handling sets of allowed runes) will be the core validator.

## 4. Decisions & API Refinement

### Naming Convention
*   **Creation**: Use `New` (e.g., `form.New(MyStruct{})`). checks semantic context.
*   **Global Registry**: Forms are automatically added to a private global handler.

### Contextual Intelligence ("Smart Fields")
The library automates configuration based on struct field names.
*   **Principle**: "Convention over Configuration".
*   **Behavior**: If a field is named `User.Name`, the library infers it is a "Name" type (Input: Text, Min: 2, Max: 100).
*   **Goal**: Minimize boilerplate for common types (Email, Age, Phone, RUT, etc.).
*   See **[Standard Types](STANDARD_TYPES.md)** for the full list of default mappings.

### Global Configuration
*   **Styling**: Use `SetGlobalClass(...string)` to define default CSS classes for all forms.
    *   Example: `form.SetGlobalClass("form-control", "p-2")`.
    *   Note: The Global Handler is private to prevent API pollution. Configuration is done via top-level functions.
