# PLAN: input: Tag Unification — tinywasm/form

**Module:** `github.com/tinywasm/form` (covers subpackage `tinywasm/form/input`)
**Breaking change:** Yes — changes `input.Input` interface, removes registry auto-init.
**Execution order:** Requires `tinywasm/fmt` PLAN_WIDGET.md to be published first.
**Note on tag cleanup:** `tinywasm/orm` (PLAN.md) handles source file tag rewriting — it removes `form:` and `validate:` tags from `model.go`/`models.go` automatically when `ormc` runs. This plan does not need to address tag cleanup.

---

## Context

`tinywasm/form` is the UI layer of the tinywasm ecosystem. It builds HTML forms from `fmt.Model` schemas.
`tinywasm/form/input` defines the `Input` interface and all concrete input types (email, text, textarea, etc.).

### Current Problems

1. **`registry.go` `init()`** auto-registers all input types at startup. This is implicit — a field named "Email" magically becomes an email input via name matching. This breaks when field names don't follow conventions and makes the system non-deterministic.

2. **No connection to `fmt.Field`** — the form layer has no access to the input type through the schema; it guesses via `findInputForField()`.

### Goal

- Remove all magic name matching from `form`.
- `input.Input` embeds `fmt.Widget` — `Clone(parentID, name)` satisfies both Widget and positioned rendering.
- Add `New<Type>() fmt.Widget` constructors for use by ormc code generation.
- Migrate `form/input/` tests to `form/input/tests/` to follow orm convention.
- Remove `init()` auto-registration from `registry.go`.
- `form` uses `field.Widget` directly from the schema to render inputs.

---

## Development Rules

- **Standard Library Only:** No external assertion libraries. Use `testing`.
- **Testing Runner:** Use `gotest` (`go install github.com/tinywasm/devflow/cmd/gotest@latest`).
- **Max 500 lines per file.** If exceeded, subdivide by domain.
- **TinyGo Compatible:** No `fmt`, `strings`, `strconv`, `errors` from stdlib. Use `tinywasm/fmt`.
- **No `reflect` at runtime.**
---

## Part 1: `tinywasm/form/input` Changes

### 1.1 Update `input.Input` interface in `interface.go`

**Current state** (partially updated — `Clone` already renamed to `Build`):
```go
type Input interface {
    dom.Component
    HTMLName() string
    FieldName() string
    ValidateField(value string) error
    Build(parentID, name string) Input
}
```

**Target:**
```go
type Input interface {
    fmt.Widget    // Type(), Validate(), Clone(parentID, name) — semantic type contract
    dom.Component // GetID(), SetID(), RenderHTML(), Children()
}
```

**Remaining changes:**
- Embed `fmt.Widget` — promotes `Type()`, `Validate()`, `Clone(parentID, name string) Widget` into `Input`
- Remove `Build(parentID, name string) Input` — replaced by `fmt.Widget.Clone(parentID, name)` (same signature, returns `Widget` which is type-assertable to `Input`)
- Remove `HTMLName() string` — replaced by `fmt.Widget.Type()` (same value; kept on `Base` as internal method)
- Remove `FieldName() string` — internal rendering concern; name passed via `Clone(parentID, name)`
- Remove `ValidateField(value string) error` — replaced by `fmt.Widget.Validate()` (same behavior; kept on `Base` as implementation detail)

**Why remove `HTMLName`, `FieldName`, `ValidateField`, `Build` from the interface:** Interface Segregation — all replaced by `fmt.Widget` embedding. `HTMLName()` = `Type()`. `FieldName()` is internal. `ValidateField()` = `Validate()`. `Build()` = `Clone()`. Concrete types retain internal methods on `Base` for rendering use.

### 1.2 Update `base.go`

Add `Type()` and `Validate()` to `Base` so all concrete types embedding `Base` satisfy `fmt.Widget`:

```go
// Type satisfies fmt.Widget.Type(). Returns the semantic input type name.
func (b *Base) Type() string { return b.htmlName }

// Validate satisfies fmt.Widget.Validate(). Delegates to ValidateField.
func (b *Base) Validate(value string) error { return b.ValidateField(value) }
```

`Clone(parentID, name string) fmt.Widget` cannot live on `Base` (doesn't know concrete type) — each concrete type must implement it. See section 1.3.

`HTMLName()`, `FieldName()`, and `ValidateField()` remain on `Base` as internal methods but are no longer part of the `Input` interface.

Remove `Build` method from concrete types (not on `Base` directly) — replaced by `Clone`. Check `base.go` only for internal calls to `Build(`.

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
func (e *email) Build(parentID, name string) Input { return Email(parentID, name) }

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

// Clone satisfies fmt.Widget — Email() returns Input which implements Widget.
func (e *email) Clone(parentID, name string) fmt.Widget { return Email(parentID, name) }
```

Apply this pattern to ALL 17 concrete types.

### 1.4 Update `permitted.go` if needed

Check `permitted.go` for any reference to `Build` — update to `Clone`.

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
inp, ok := field.Widget.Clone(parentID, field.Name).(input.Input)
if !ok {
    continue // custom type not implementing input.Input — skip rendering
}
// inp is already positioned — ready for rendering
```

Fields without `Widget` are simply not rendered in the form — this is intentional. No magic fallback.

### 2.4 Keep `findInputByType()` if used elsewhere

Check if `findInputByType(htmlType string)` is used outside of `findInputForField`. If not, delete it too. If yes, keep it for explicit lookup by HTML type name.

---

## Files to Modify

### `form/input/`
| File | Change |
|---|---|
| `interface.go` | Embed `fmt.Widget`; remove `HTMLName()`, `FieldName()`, `ValidateField()`, `Build()` |
| `base.go` | Add `Type() string`, `Validate(value string) error`; remove `Build()` |
| `email.go` | Add `NewEmail()`; replace `Build` with `Clone(parentID, name string) fmt.Widget` |
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
| Form rendering files | Replace `findInputForField()` calls with `field.Widget.Clone(parentID, field.Name).(input.Input)` |

---

## Tests

### Test location

Tests currently live alongside source files. This is acceptable for `form/input/` since tests are already in:
- `form/input/inputs_test.go`
- `form/input/validation_test.go`
- `form/input/render_test.go`

Form-level tests are in:
- `form/base.front_test.go`
- `form/base.back_test.go`
- `form/base.shared_test.go`
- `form/setup_test.go`

Keep existing location — do NOT move to `tests/` subdirectory.

### `form/input/` test updates

Update existing tests that call `Build(parentID, name)` to call `Clone(parentID, name)`.

Add per-type tests in `inputs_test.go`:
- `New<Type>()` returns a non-nil `fmt.Widget`
- `Clone(parentID, name)` returns a non-nil `fmt.Widget` with correct `Type()`
- `Clone(parentID, name)` result is type-assertable to `input.Input`

Update existing validation tests in `validation_test.go`:
- Replace any `ValidateField()` calls with `Validate()`
- `Validate(validValue)` returns nil
- `Validate(invalidValue)` returns error

Update existing render tests in `render_test.go`:
- Replace `Build()` calls with `Clone()` for positioned instances
- Verify `Clone(parentID, name).(input.Input).RenderHTML()` produces correct output

### `form/` test updates

Update existing tests in `base.shared_test.go` or `base.back_test.go`:
- Form built from schema where `Widget == nil` does not render that field
- `field.Widget.Clone(parentID, field.Name).(input.Input)` produces correct HTML for each known type
- Replace any calls to `findInputForField()` with `field.Widget.Clone()` pattern

---

## go.mod Update

```bash
go get github.com/tinywasm/fmt@v0.21.1
go mod tidy
```


