# Design & Architecture

## Philosophy
*   **Minimalism**: Small binary size > Feature completeness.
*   **TinyTime**: Optimize for WebAssembly/TinyGo.
*   **Convention over Configuration**: Reduce boilerplate (tags).

## Global Registry Pattern
Instead of managing form instances manually, the library maintains a global slice of forms.
*   **Why**: Simplifies the consumer code. Just call `New` and it's managed.
*   **Storage**: `[]*Form`. Slices are more memory-efficient and GC-friendly in TinyGo than Maps for small collections.

## Reflection Strategy
We uses `reflect` to inspect the struct fields.
*   **Performance**: Reflection is slower than generated code, but `tinywasm/form` is an initialization-time cost (mostly).
*   **Field Discovery**: We use `tinywasm/fmt.HasUpperPrefix` to filter only exported (public) fields, treating them as form inputs.

## Validation
Validation logic is decoupled from the struct itself.
*   **Shared Logic**: `Permitted` structs define rules (min/max/charset).
*   **Reusability**: Define a "Phone" validation rule once, apply it to any field named "Phone" or typed as "Phone".

## Input Interface Strategy (`form/input`)

To handle the variety of input types (Text, Email, RUT, etc.) without a monolithic switch statement, we use a polymorphic interface strategy.

### Interface Definition
```go
type Input interface {
    dom.Component // Embeds ID() and RenderHTML()
    HtmlName() string             // Standard HTML5 type (e.g., "text", "email")
    ValidateField(value string) error // Self-contained validation logic
}
```
*   **Optional Interfaces**: `dom.CSSRenderer` and `dom.JSRenderer` are supported via type assertion.

### Trade-off Analysis

**Pros (Advantages):**
*   **Decoupled Maintenance**: Each input type lives in its own file (e.g., `input/rut.go`, `input/email.go`). This fixes the "giant switch-case" problem from the legacy code.
*   **Encapsulation**: Validation logic is tightly coupled with the Input definition, not scattered.
*   **Extensibility**: Adding a new type just means creating a struct that implements `Input`. The core `form` package doesn't need changes.

**Cons (Challenges):**
*   **State Management**: If `ValidateField` depends on dynamic constraints (like a variable MinLength), the `Input` implementation must carry that state (e.g., `RUT` struct needs a `Permitted` field).
*   **Memory Overhead**: Creating distinct struct instances for every field might have a slightly higher memory footprint in TinyGo than a shared configuration map, but it acts as a cleaner abstraction.
