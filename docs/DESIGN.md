# Design & Architecture

## Philosophy
- **Minimalism**: Small binary size > Feature completeness.
- **TinyGo Optimized**: Flat slices, minimal allocations, zero-allocation lookups.
- **Convention over Configuration**: Struct field names and [Tags](TAGS.md) drive behavior.

## Core Layers

### 1. Global Registry (`registry.go`)
Manages `registeredInputs` and `forms`. 
- **Matching**: `New()` searches for a match using:
  1. `FieldName` vs `htmlName` or `aliases`.
  2. `StructName.FieldName` vs `aliases` (allows field-specific specialized inputs).
- **Extensibility**: Anyone can `RegisterInput()` a component that implements the [Input Interface](API.md).

### 2. State & Binding Layer
- **Storage**: `input.Base` holds the state (`Values`, `Options`).
- **One-Way Binding (Creation)**: `New()` copies struct values to inputs.
- **Two-Way Sync (Interaction)**: `SyncValues()` reflects input changes back into the original struct.

### 3. Clone Pattern
Each Input implements `Clone(parentID, name string) Input`. This allows dynamic instantiation without huge switch-case blocks in the orchestrator.

### 4. Interactivity Strategy
Uses **event delegation** to minimize memory overhead. 
One listener at the Form root delegates to individual inputs. 
See [mount.go](../mount.go) and [Interactivity Strategy](INTERACTIVITY_AND_MOUNTING.md).

### 5. Validation Engine
Whitelist-based validation using the `Permitted` struct. 
- Fast character set checks.
- Length constraints.
- Custom logic via `ExtraValidation`.
- See [input/permitted.go](../input/permitted.go).
