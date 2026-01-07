# Design & Architecture

## Philosophy
- **Minimalism**: Small binary size > Feature completeness.
- **TinyGo Optimized**: No maps, minimal allocations.
- **Convention over Configuration**: Field names drive behavior.

## Global Registry
Forms and inputs are registered globally using slices.
- `forms []*Form`: All form instances.
- `registeredInputs []input.Input`: Input templates for field matching.

## Input Matching
When `New()` processes a struct field:
1. Iterates `registeredInputs`.
2. Calls `input.Matches(fieldName)` which checks `htmlName` and `aliases`.
3. If match found, calls `input.Clone(parentID, fieldName)`.
4. If no match, returns error.

## Clone Pattern
Each Input implements `Clone(parentID, name string) Input`.
This eliminates switch cases and allows new inputs without modifying form.go.

```go
func (t *text) Clone(parentID, name string) Input {
    return Text(parentID, name)
}
```

## Event Handling
**Centralized listener** at Form level. One listener catches all input events and routes to the correct input via ID lookup.

See [INTERACTIVITY_AND_MOUNTING.md](INTERACTIVITY_AND_MOUNTING.md).

## Validation
Uses `Permitted` struct for whitelist-based validation:
- `Letters`, `Numbers`, `Tilde`: Character sets.
- `Characters []rune`: Specific allowed chars.
- `Min/Max`: Length constraints.
- `ExtraValidation`: Custom function.

See [input/permitted.go](../input/permitted.go).
