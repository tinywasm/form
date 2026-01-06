# Interactivity and Mounting Strategy

This document addresses how `tinywasm/form` integrates with `tinywasm/dom` for component mounting and event handling.

## 1. The Mounting Challenge
`tinywasm/dom` operates by:
1.  Rendering HTML string (`RenderHTML`).
2.  Injecting it into the DOM.
3.  Calling `OnMount()` if the component implements `dom.Mountable`.

**Problem**: The `Form` generates one large HTML string containing all inputs. Inputs are not mounted individually by the user; they are children of the Form.

## 2. Proposal: Hierarchical Mounting

We should use a delegated mounting strategy where the **Form** acts as the orchestrator.

### Structure
*   **Form**: Implements `dom.Mountable`.
*   **Input**: Implements `dom.Mountable`.

### Workflow
1.  **User** calls `dom.Mount("root", myForm)`.
2.  **DOM** renders `myForm.RenderHTML()` (which includes inputs' HTML).
3.  **DOM** calls `myForm.OnMount()`.
4.  **Form** iterates through its `Inputs` and manually calls `input.OnMount()`.

```go
// Form.OnMount
func (f *Form) OnMount() {
    for _, input := range f.Inputs {
        if mountable, ok := input.(dom.Mountable); ok {
            mountable.OnMount()
        }
    }
}
```

## 3. Event Handling (Input Interaction)

Each input is responsible for its own behavior ("Smart Components").

### Input.OnMount
Inside `OnMount`, the input attaches event listeners using its known ID.

```go
// Text.OnMount
func (t *text) OnMount() {
    // 1. Get Element from DOM
    el, found := dom.Get(t.id)
    if !found { return }

    // 2. Bind Validation
    el.AddEventListener("input", func(e dom.Event) {
        val := el.Value() // Get current value
        err := t.ValidateField(val)
        // Handle UI error feedback (toggle invalid class, etc.)
    })
}
```

## 4. Pros & Cons

### Pros
*   **Encapsulation**: Inputs manage their own validation logic and events. The Form doesn't need to know *how* a specific input works.
*   **Performance**: Event listeners are specific to elements.
*   **Consistency**: Follows `tinywasm/dom` lifecycle.

### Cons
*   **Boilerplate**: Every input type must implement `OnMount`. (Can be mitigated by a `BaseMountable` struct or just convention).
*   **Manual Delegation**: The Form *must* remember to call children's `OnMount`.

## 5. Configuration
*   **Auto-Mount**: Inputs are "auto-mounted" by the Form.
*   **Validators**: Configured via `Permitted` struct, executed inside the event listener.
