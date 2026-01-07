# Interactivity and Mounting Strategy

## Overview
Forms use a **Centralized Event Listener** at the Form root level. One listener captures all events (`input`, `change`, `submit`) and delegates to the correct Input based on `event.target.id`.

## Why Centralized?
*   **Memory Efficient**: One closure/listener per form instead of N.
*   **Smaller Binary**: Fewer closures reduce WASM code size.
*   **Deterministic State**: All synchronization logic is in one place.

## Implementation Details
The actual implementation resides in [mount.go](../mount.go).

### Key steps in `OnMount()`:
1.  **Event Delegation**: Captures `input` and `change` events.
2.  **Live Binding**: Matches `event.target.id` against `f.Inputs` and calls `SetValues()`.
3.  **Automatic Validation**: Triggers `ValidateField()` immediately on change.
4.  **Submission Control**: Prevents default browser submit, calls `SyncValues()` to update the source struct, and triggers the `OnSubmit` callback.

## Data Synchronization
State flows both ways:
1.  **Mounting**: Struct values are copied to Inputs.
2.  **Interaction**: User input updates Input state.
3.  **Sync/Submit**: `SyncValues()` reflects all Input states back into the original struct.

## Mounting Flow
1.  User calls `dom.Mount("root", myForm)`.
2.  DOM injects `myForm.RenderHTML()`.
3.  DOM calls `myForm.OnMount()` (if component implements `dom.Mountable`).
4.  Form attaches centralized listeners.
