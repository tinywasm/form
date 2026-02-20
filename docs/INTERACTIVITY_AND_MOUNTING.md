# Interactivity and Mounting

> Summary in `README.md` under "WASM Event Flow". This file has additional details.

## `OnMount()` — WASM Only (`mount.go`)

Called automatically by `tinywasm/dom` after the form is injected into the DOM.

```
dom.Mount("root-id", myForm)
  1. dom calls myForm.RenderHTML() → injects HTML
  2. dom calls myForm.OnMount()
     → dom.Get(f.GetID()) → gets <form> element
     → el.On("input",  onInput)   ← live sync per keystroke
     → el.On("change", onInput)   ← for select/radio/checkbox
     → el.On("submit", onSubmit)  ← intercepts form submit
```

## Event Handlers

**onInput** (input + change events):
1. `e.TargetID()` → find matching `inp` in `f.Inputs`
2. `inp.SetValues(e.TargetValue())`
3. `inp.ValidateField(val)` — immediate feedback

**onSubmit**:
1. `e.PreventDefault()`
2. `f.SyncValues()` — flush input values to struct
3. `f.Validate()` — full validation (returns on first error)
4. `f.onSubmit(f.Value)` — user callback with populated struct

## Why Centralized?

One listener per form (not per input) = fewer closures = smaller WASM binary.
