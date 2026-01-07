# Interactivity and Mounting Strategy

## Overview
Forms use a **Centralized Event Listener** at the Form level. One listener captures all input events and delegates to the correct Input based on `event.target.id`.

## Why Centralized?
*   **Memory Efficient**: One closure per form (not N closures).
*   **Smaller Binary**: Fewer closures = less WASM code.
*   **Control**: All event logic lives in `Form.OnMount`.

## Implementation

### Form.OnMount
```go
func (f *Form) OnMount() {
    formEl, _ := dom.Get(f.ID)
    
    formEl.AddEventListener("input", func(e dom.Event) {
        targetID := e.TargetID()
        for _, input := range f.Inputs {
            if input.ID() == targetID {
                err := input.ValidateField(e.TargetValue())
                // Handle error UI feedback here
                break
            }
        }
    })
}
```

### Input Role
Inputs provide:
*   `ID()` for identification.
*   `ValidateField(value)` for validation logic.

Inputs do **not** attach their own listeners (unless they are special cases like RichText editors).

## Mounting Flow
1.  User calls `dom.Mount("root", myForm)`.
2.  DOM injects `myForm.RenderHTML()`.
3.  DOM calls `myForm.OnMount()`.
4.  Form attaches centralized listener.
