# PLAN: input: Tag Unification — tinywasm/form

**Module:** `github.com/tinywasm/form` (covers subpackage `tinywasm/form/input`)
**Breaking change:** Yes — changes `input.Input` interface, removes registry auto-init.
**Execution order:** Requires `tinywasm/fmt` PLAN_FIELD_INPUT.md to be published first.
**Note on tag cleanup:** `tinywasm/orm` (PLAN_INPUT_TAG.md) handles source file tag rewriting — it removes `form:` and `validate:` tags from `model.go`/`models.go` automatically when `ormc` runs. This plan does not need to address tag cleanup.

---

## Context

`tinywasm/form` is the UI layer of the tinywasm ecosystem. It builds HTML forms from `fmt.Model` schemas.
`tinywasm/form/input` defines the `Input` interface and all concrete input types (email, text, textarea, etc.).

### Current Problems

1. **`registry.go` `init()`** auto-registers all input types at startup. This is implicit — a field named "Email" magically becomes an email input via name matching. This breaks when field names don't follow conventions and makes the system non-deterministic.

2. **`input.Input.Clone(parentID, name string) Input`** creates positioned instances for form rendering. There is no parameterless constructor for use as a schema template.

3. **No connection to `fmt.Field`** — the form layer has no access to the input type through the schema; it guesses via `findInputForField()`.

### Goal

- Remove all magic name matching from `form`.
- Add `Clone() fmt.Widget` (parameterless) to `input.Input` so templates can be stored in `fmt.Field.Widget`.
- Add `New<Type>() fmt.Widget` constructors for use by ormc code generation.
- Remove `init()` auto-registration from `registry.go`.
- `form` uses `field.Widget` directly from the schema to render inputs.

---

## Development Rules

- **Standard Library Only:** No external assertion libraries. Use `testing`.
- **Testing Runner:** Use `gotest` (`go install github.com/tinywasm/devflow/cmd/gotest@latest`).
- **Max 500 lines per file.** If exceeded, subdivide by domain.
- **TinyGo Compatible:** No `fmt`, `strings`, `strconv`, `errors` from stdlib. Use `tinywasm/fmt`.
- **No `reflect` at runtime.**
- **Publishing:** Use `gopush 'message'` after tests pass.

---

## Part 1: `tinywasm/form/input` Changes

### 1.1 Update `input.Input` interface in `interface.go`

**Before:**
```go
type Input interface {
    dom.Component
    HTMLName() string
    FieldName() string
    ValidateField(value string) error
    Clone(parentID, name string) Input
}
```

**After:**
```go
type Input interface {
    dom.Component                          // GetID(), SetID(), RenderHTML(), Children()
    Type() string                          // Semantic input type (e.g., "email", "textarea") — satisfies fmt.Widget.Type()
    HTMLName() string                      // HTML5 type attribute — same value as Type(), kept for dom/rendering context
    FieldName() string                     // Struct field name (without parent prefix)
    ValidateField(value string) error      // Semantic validation (existing — kept for input.Input consumers)
    Validate(value string) error           // Alias for ValidateField — satisfies fmt.Widget.Validate()
    Build(parentID, name string) Input     // Creates positioned instance for form rendering (renamed from Clone)
    Clone() fmt.Widget                     // Returns a fresh template instance — satisfies fmt.Widget.Clone()
}
```

**Changes:**
- `Clone(parentID, name string) Input` → renamed to `Build(parentID, name string) Input`
- New `Validate(value string) error` — added to interface, delegates to `ValidateField` (satisfies `fmt.Widget`)
- New `Clone() fmt.Widget` — added to interface, returns a template copy with no parentID/name (satisfies `fmt.Widget`)

**Why rename `Clone` to `Build`:** Go does not allow two methods with the same name and different signatures. `Clone()` (no params) is required by `fmt.Widget`. The existing positioned constructor is semantically a "build for position", so `Build` is more accurate.

**Why both `Validate` and `ValidateField`:** `ValidateField` is the existing method used by internal form rendering code. `Validate` is required by `fmt.Widget`. Both live in the interface so existing consumers of `input.Input` do not break, and `fmt.Widget` is satisfied. On `Base`, `Validate` delegates to `ValidateField`.

### 1.2 Update `base.go`

Add `Clone() fmt.Widget` to `Base` — all concrete types embedding `Base` automatically inherit it if they have a `New<Type>()` constructor.

Since `Base` doesn't know its concrete type, `Clone()` cannot be implemented on `Base` itself — each concrete type must implement it. See section 1.3.

Rename `Clone` usage in `Base` if any exist:
- Search `base.go` for any call to `Clone(` and update to `Build(`.

### 1.3 Update ALL concrete input types

For each type in `email.go`, `text.go`, `password.go`, `textarea.go`, `phone.go`, `number.go`, `date.go`, `hour.go`, `ip.go`, `rut.go`, `address.go`, `checkbox.go`, `datalist.go`, `select.go`, `radio.go`, `filepath.go`, `gender.go`:

**Pattern (example for `email`):**

```go
// Before:
func Email(parentID, name string) Input {
    e := &email{}
    // ... init ...
    return e
}
func (e *email) Clone(parentID, name string) Input { return Email(parentID, name) }

// After:
// NewEmail returns a template instance for use in fmt.Field.Widget (no position).
// Used by ormc-generated schema code.
func NewEmail() fmt.Widget {
    return Email("", "")
}

func Email(parentID, name string) Input {
    e := &email{}
    // ... init (unchanged) ...
    return e
}

// Build creates a positioned instance for form rendering.
func (e *email) Build(parentID, name string) Input { return Email(parentID, name) }

// Clone returns a fresh template instance implementing fmt.Widget.
func (e *email) Clone() fmt.Widget { return NewEmail() }
```

Apply this pattern to ALL 17 concrete types.

### 1.4 Add `Type() string` to `Base`

`fmt.Widget` requires `Type() string`. `Base` already has `HTMLName() string` returning the same value. Add a `Type()` method on `Base` that delegates to `htmlName`:

```go
// In base.go
func (b *Base) Type() string { return b.htmlName }
```

This satisfies `fmt.Widget.Type()` for all types embedding `Base`. `HTMLName()` is kept unchanged — it is used by the rendering layer (`dom` context) where the HTML attribute name is explicit.

### 1.5 Satisfy `fmt.Widget.Validate()` on concrete types

`fmt.Widget` requires `Validate(value string) error`. Each concrete type already has `ValidateField(value string) error` on `Base`. Add a `Validate` alias on `Base`:

```go
// In base.go
func (b *Base) Validate(value string) error { return b.ValidateField(value) }
```

`ValidateField` is kept for backwards compatibility within the `input.Input` interface.

### 1.6 Update `permitted.go` if needed

Check `permitted.go` for any reference to `Clone` — update to `Build`.

---

## Part 2: `tinywasm/form` Changes

### 2.1 Remove `init()` from `registry.go`

Delete the entire `init()` block that calls `RegisterInput(...)` with all default inputs.

The registry and `RegisterInput()` function can be kept for projects that want runtime-registered custom inputs via old API — but it is no longer called automatically.

**Before:**
```go
func init() {
    RegisterInput(
        input.Text("", ""),
        input.Email("", ""),
        // ... 15 more ...
    )
}
```

**After:** Delete `init()` entirely. `RegisterInput` stays as exported function but is a no-op until called explicitly by the project.

### 2.2 Remove `findInputForField()` from `registry.go`

Delete `findInputForField(fieldName, structName string) input.Input` — this is the magic name-matching function. It is replaced by direct use of `field.Widget`.

### 2.3 Update form rendering to use `field.Widget` directly

Wherever `form` calls `findInputForField(field.Name, structName)` to get an input, replace with:

```go
if field.Widget == nil {
    continue // field has no UI binding — skip
}
inp, ok := field.Widget.Clone().(input.Input)
if !ok {
    continue // custom type not implementing input.Input — skip rendering
}
rendered := inp.Build(parentID, field.Name)
```

Fields without `Widget` are simply not rendered in the form — this is intentional. No magic fallback.

### 2.4 Keep `findInputByType()` if used elsewhere

Check if `findInputByType(htmlType string)` is used outside of `findInputForField`. If not, delete it too. If yes, keep it for explicit lookup by HTML type name.

---

## Files to Modify

### `form/input/`
| File | Change |
|---|---|
| `interface.go` | Rename `Clone → Build`, add `Validate(value string) error`, add `Clone() fmt.Widget` |
| `base.go` | Add `Type() string`, `Validate(value string) error`, rename any `Clone(` calls to `Build(` |
| `email.go` | Add `NewEmail()`, rename `Clone → Build`, add `Clone()` |
| `text.go` | Same pattern |
| `password.go` | Same pattern |
| `textarea.go` | Same pattern |
| `phone.go` | Same pattern |
| `number.go` | Same pattern |
| `date.go` | Same pattern |
| `hour.go` | Same pattern |
| `ip.go` | Same pattern |
| `rut.go` | Same pattern |
| `address.go` | Same pattern |
| `checkbox.go` | Same pattern |
| `datalist.go` | Same pattern |
| `select.go` | Same pattern |
| `radio.go` | Same pattern |
| `filepath.go` | Same pattern |
| `gender.go` | Same pattern |

### `form/`
| File | Change |
|---|---|
| `registry.go` | Delete `init()`, delete `findInputForField()` |
| Form rendering files | Replace `findInputForField()` calls with `field.Widget.Clone().(input.Input).Build(...)` |

---

## Tests

### `form/input/` tests

Update existing tests that call `Clone(parentID, name)` to call `Build(parentID, name)`.

Add per-type tests for:
- `New<Type>()` returns a non-nil `fmt.Widget`
- `Clone()` returns a non-nil `fmt.Widget` with correct `Name()`
- `Validate(validValue)` returns nil
- `Validate(invalidValue)` returns error
- `Build(parentID, name)` returns a correctly positioned `Input`

### `form/` tests

Add test that verifies a form built from a schema where `Widget == nil` does not render that field.
Add test that verifies `field.Widget.Clone().(input.Input).Build(...)` produces correct HTML for each known type.

---

## go.mod Update

```bash
go get github.com/tinywasm/fmt@v0.21.0
go mod tidy
```


