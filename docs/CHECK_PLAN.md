# PLAN: Form v3 — Field v3 Migration + Unified Validation (tinywasm/form)

← [README](../README.md) | Depends on: [fmt PLAN.md](../../fmt/docs/PLAN.md), [orm PLAN.md](../../orm/docs/PLAN.md)

## Development Rules

- **Standard Library Only:** No external assertion libraries. Use `testing`.
- **Testing Runner:** Use `gotest` (install: `go install github.com/tinywasm/devflow/cmd/gotest@latest`).
- **Max 500 lines per file.** If exceeded, subdivide by domain.
- **Flat hierarchy.** No subdirectories for library code.
- **TinyGo Compatible:** No `fmt`, `strings`, `strconv`, `errors` from stdlib. Use `tinywasm/fmt`.
- **Documentation First:** Update docs before coding.
- **Publishing:** Use `gopush 'message'` after tests pass and docs are updated.

## Prerequisite

Update `go.mod` to the new `tinywasm/fmt` version (Field v3):

```bash
go get github.com/tinywasm/fmt@v0.19.0
```

## Context

With Field v3:
- `Field.Input` is removed → form resolves input type **only** by field name heuristic.
- `Field` embeds `fmt.Permitted` → validation rules live in the schema.
- `Field.Validate(value)` is a method on Field — validates using embedded Permitted.
- `fmt.ValidateFielder(data)` is the generic validation function.
- `input.Permitted` is deleted → replaced by `fmt.Permitted` (no maps, ASCII ranges).

### Current validation flow (v2):
```
form.ValidateData() → for each input → input.ValidateField(stringValue)
                                         └→ input.Permitted.Validate(value)  (form/input/permitted.go, uses maps)
```
Each input has its own `Permitted` config with map-based character validation.

### New validation flow (v3):
```
form.ValidateData() → fmt.ValidateFielder(data)
                        └→ for each field → field.Validate(value)
                                             └→ field.Permitted.Validate(field.Name, value)  (fmt/permitted.go, ASCII ranges)

input.ValidateField() → field.Permitted.Validate(field.Name, value)  (client-side UX only)
```

Validation is driven by the **schema** (Field embeds Permitted), not by input instances.
Form inputs keep `ValidateField()` for client-side WASM UX feedback only.

---

## Stage 1: Remove `Field.Input` usage from `form.go`

**File:** `form.go`

### 1.1 Remove explicit Input type override

```go
// BEFORE (lines 84-99):
if field.Input == "-" {
    continue
}
// ...
var template input.Input
if field.Input != "" {
    template = findInputByType(field.Input)
}
if template == nil {
    template = findInputForField(fieldName, structName)
}

// AFTER:
var template input.Input
template = findInputForField(fieldName, structName)
```

### 1.2 Handle fields with no matching input

Default to Text input for unmatched fields:

```go
if template == nil {
    template = findInputByType("text") // fallback to text input
}
```

### 1.3 Skip mechanism without `Input: "-"`

```go
// BEFORE:
if field.PK && field.AutoInc {
    continue
}
if field.Input == "-" {
    continue
}

// AFTER:
if field.PK {
    continue // PKs are never editable in forms
}
```

**Note:** If a dev needs a PK in a form (rare), they use `ormc:formonly` struct without PK flag.

---

## Stage 2: Delete `input/permitted.go` — replaced by `fmt.Permitted`

### 2.1 Delete the file

`form/input/permitted.go` is completely replaced by `fmt/permitted.go`.
Delete it.

### 2.2 Update `input/base.go` — change embed

```go
// BEFORE (line 24):
Permitted      // anonymous embed: promotes Letters, Numbers, Validate(), etc.

// AFTER:
fmt.Permitted  // anonymous embed: promotes Letters, Numbers, Validate(), etc.
```

### 2.3 Update `Base.ValidateField` signature

`fmt.Permitted.Validate` now takes `(field, text string)` instead of just `(text string)`:

```go
// BEFORE (base.go:128-130):
func (b *Base) ValidateField(value string) error {
    return b.Permitted.Validate(value)
}

// AFTER:
func (b *Base) ValidateField(value string) error {
    return b.Permitted.Validate(b.FieldName(), value)
}
```

### 2.4 Update all inputs that call `Permitted.Validate`

**Files:** `hour.go`, `ip.go`, `rut.go`, `date.go`, `filepath.go`

All follow the same pattern — add field name as first arg:

```go
// BEFORE:
if err := x.Permitted.Validate(value); err != nil {

// AFTER:
if err := x.Permitted.Validate(x.FieldName(), value); err != nil {
```

### 2.5 Update input constructors

Each input constructor configures a `Permitted` struct. Change type:

```go
// BEFORE:
return &EmailInput{Base: Base{Permitted: Permitted{Letters: true, ...}}}

// AFTER:
return &EmailInput{Base: Base{Permitted: fmt.Permitted{Letters: true, ...}}}
```

**Note:** With Field v3, input Permitted configs should match what ormc generates in the schema.
The constructors still exist for standalone form usage (without ormc), but for ormc-generated
models the Permitted config comes from the schema, not the input.

---

## Stage 3: Update `ValidateData` to use `fmt.ValidateFielder`

**File:** `validate_struct.go`

```go
// BEFORE:
func (f *Form) ValidateData(action byte, data fmt.Fielder) error {
    values := fmt.ReadValues(data.Schema(), data.Pointers())
    for i, inp := range f.Inputs {
        idx := f.fieldIndices[i]
        if idx < 0 || idx >= len(values) {
            continue
        }
        if skipper, ok := inp.(interface{ GetSkipValidation() bool }); ok && skipper.GetSkipValidation() {
            continue
        }
        val := fmt.Convert(values[idx]).String()
        if err := inp.ValidateField(val); err != nil {
            return err
        }
    }
    return nil
}

// AFTER:
func (f *Form) ValidateData(action byte, data fmt.Fielder) error {
    return fmt.ValidateFielder(data)
}
```

**Impact:** Validation is now driven by the schema (Field.Validate), not by form inputs.
Same rules apply whether data comes from form, JSON, or API.

---

## Stage 4: Clean up registry

**File:** `registry.go`

### 4.1 Simplify text fallback

`findInputByType` is now only used for the text fallback. Inline it:

```go
var textFallback = input.Text("", "")

// In form.go:
if template == nil {
    template = textFallback
}
```

### 4.2 Optionally remove `findInputByType`

If no other code calls it, delete the function.

---

## Stage 5: Update tests

### 5.1 Update test Field literals

Remove `Input:` from all test `fmt.Field` structs:

```go
// BEFORE:
{Name: "email", Type: fmt.FieldText, Input: "email"}

// AFTER:
{Name: "email", Type: fmt.FieldText, Permitted: fmt.Permitted{Letters: true, Numbers: true}}
```

### 5.2 Update `ValidateData` tests

Tests should verify `ValidateData` calls `fmt.ValidateFielder` (schema-driven).

### 5.3 Add test: form with no matching input → falls back to text

```go
func TestUnmatchedFieldFallsBackToText(t *testing.T) {
    // Field named "custom_xyz" matches no registered input
    // Should get Text input instead of error
}
```

### 5.4 Add test: Permitted from fmt works in inputs

Verify that `input.Base` with `fmt.Permitted` validates correctly.

### 5.5 Run tests

```bash
gotest
```

---

## Stage 6: Update documentation

### 6.1 Update `docs/SKILL.md`

- Remove `Field.Input` references.
- Document that form resolves inputs by field name only.
- Document the text fallback for unmatched fields.
- Document new validation flow: `ValidateData` → `fmt.ValidateFielder` → `Field.Validate`.

### 6.2 Update `docs/DESIGN.md`

- Update architecture diagram showing validation flow.
- Remove `Input` from Field diagram.
- Show `fmt.Permitted` embedded in Field and in input.Base.

---

## Stage 7: Publish

```bash
gopush 'form: Field v3 — delete Permitted (now fmt.Permitted), remove Input, ValidateData uses fmt.ValidateFielder'
```

---

## Summary

| Stage | File(s) | Action |
|-------|---------|--------|
| 1 | `form.go` | Remove `Field.Input` usage, fallback to text, skip all PKs |
| 2 | `input/permitted.go` → delete, `input/base.go` + all inputs | Replace `input.Permitted` with `fmt.Permitted`, update Validate signature |
| 3 | `validate_struct.go` | Replace per-input validation with `fmt.ValidateFielder(data)` |
| 4 | `registry.go` | Simplify input lookup, inline text fallback |
| 5 | `*_test.go` | Update Field literals with Permitted, validation tests |
| 6 | `docs/` | Update SKILL.md, DESIGN.md |
| 7 | — | `gotest` + `gopush` |

## Execution Order (cross-package)

```
1. fmt   → Field v3 (Permitted structure extraction, Field.Validate, ValidateFielder)
2. orm   → Ormc v3 (Permitted in schema, composite standalone Validate, validate tags)
3. json  → JSON v3 (Field.Name as key, OmitEmpty, post-decode ValidateFielder)
4. form  → Form v3 (delete Permitted, fmt.Permitted in inputs, ValidateData → ValidateFielder)
```

Each step requires publishing before the next can begin (`gopush`).
