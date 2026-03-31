# PLAN: Remove Duplicate Permitted — Use fmt.Permitted in form/input

**Module:** `github.com/tinywasm/form` (covers subpackage `tinywasm/form/input`)
**Breaking change:** Yes — removes `input.Permitted` struct, renames fields on `Base`.
**Dependency:** `tinywasm/fmt` — NO changes needed. `fmt.Permitted` already has everything required.

---

## Context

`input/permitted.go` contains a duplicated `Permitted` struct with its own `Validate()` method and character maps. `fmt/permitted.go` already has the canonical version — optimized (ASCII ranges instead of maps), with the same fields under cleaner names, and a 2-param `Validate(field, text)` signature that includes the field name in errors.

### Goal

- Delete `input/permitted.go` entirely.
- `Base` embeds `fmt.Permitted` instead of the local `Permitted`.
- Rename field references across all 17 concrete inputs.
- Update all `Permitted.Validate(value)` calls to `Permitted.Validate(name, value)`.

---

## Development Rules

- **Standard Library Only:** No external assertion libraries. Use `testing`.
- **Testing Runner:** Use `gotest` (`go install github.com/tinywasm/devflow/cmd/gotest@latest`).
- **Max 500 lines per file.** If exceeded, subdivide by domain.
- **TinyGo Compatible:** No `fmt`, `strings`, `strconv`, `errors` from stdlib. Use `tinywasm/fmt`.
- **No `reflect` at runtime.**

---

## Part 1: Delete `input/permitted.go`

Delete the entire file. It contains:
- `Permitted` struct (duplicated from `fmt.Permitted`)
- `Validate(text string) error` method
- `MinMaxAllowedChars()` method (unused)
- `valid_letters`, `valid_tilde`, `valid_number` maps (replaced by ASCII ranges in fmt)
- Constants `tabulation`, `white_space`, `break_line`

---

## Part 2: Update `input/base.go`

### 2.1 Change embedded struct

```go
// Before:
type Base struct {
    // ...
    Permitted // anonymous embed: promotes Letters, Numbers, Validate(), etc.
}

// After:
type Base struct {
    // ...
    fmt.Permitted // anonymous embed: promotes Letters, Numbers, Validate(), etc.
}
```

### 2.2 Update `Base.Validate()` call

```go
// Before (line 136):
return b.Permitted.Validate(value)

// After:
return b.Permitted.Validate(b.name, value)
```

No other changes to base.go. `Type()`, `FieldName()`, `HTMLName()` are unaffected.

---

## Part 3: Rename fields in ALL concrete inputs

Field name mapping (input.Permitted → fmt.Permitted):

| Old Name | New Name |
|---|---|
| `WhiteSpaces` | `Spaces` |
| `Tabulation` | `Tab` |
| `Characters` | `Extra` |
| `TextNotAllowed` | `NotAllowed` |
| `SkipRules` | _(delete — see Part 4)_ |
| `ExtraValidation` | _(delete — unused)_ |

### Files that use renamed fields:

| File | Old Field | New Field |
|---|---|---|
| textarea.go | `WhiteSpaces`, `BreakLine`, `Characters` | `Spaces`, `BreakLine`, `Extra` |
| address.go | `WhiteSpaces`, `Characters` | `Spaces`, `Extra` |
| hour.go | `Characters` | `Extra` |
| email.go | `Characters` | `Extra` |
| text.go | `Characters` | `Extra` |
| password.go | `Characters` | `Extra` |
| phone.go | `Characters` | `Extra` |
| date.go | `Characters` | `Extra` |
| ip.go | `Characters` | `Extra` |
| filepath.go | `Characters` | `Extra` |
| rut.go | `Characters` | `Extra` |

Fields that keep the same name (no change needed): `Letters`, `Tilde`, `Numbers`, `BreakLine`, `Minimum`, `Maximum`.

---

## Part 4: Handle SkipRules removal

`SkipRules` was used by `checkbox.go` and `datalist.go`. Both types already **override** `Validate()` entirely on the concrete type — they never call `Permitted.Validate()`. The `SkipRules = true` line is therefore dead code.

**Action:** Delete `c.SkipRules = true` from checkbox.go and `dl.SkipRules = true` from datalist.go.

---

## Part 5: Update Validate calls in concrete types with overrides

These files call `Permitted.Validate(value)` internally and need the field name added:

| File | Before | After |
|---|---|---|
| date.go | `d.Permitted.Validate(value)` | `d.Permitted.Validate(d.name, value)` |
| filepath.go | `fp.Permitted.Validate(value)` | `fp.Permitted.Validate(fp.name, value)` |
| hour.go | `h.Permitted.Validate(value)` | `h.Permitted.Validate(h.name, value)` |
| ip.go | `i.Permitted.Validate(value)` | `i.Permitted.Validate(i.name, value)` |
| rut.go | `r.Permitted.Validate(value)` | `r.Permitted.Validate(r.name, value)` |

---

## Part 6: Update tests

### input/validation_test.go

If tests call `Permitted.Validate(value)` directly, update to `Permitted.Validate("field", value)`.

If tests reference old field names (`WhiteSpaces`, `Characters`, etc.), rename them.

### input/inputs_test.go, input/render_test.go

Check for any `Permitted{}` struct literals with old field names — rename accordingly.

### form/ tests (base.back_test.go, base.shared_test.go, etc.)

Check for references to old field names or `Permitted.Validate(value)` calls — update.

---

## Files to Modify

### `form/input/`
| File | Change |
|---|---|
| `permitted.go` | **DELETE** |
| `base.go` | Embed `fmt.Permitted` instead of `Permitted`; update `Validate` call |
| `email.go` | `Characters` → `Extra` |
| `text.go` | `Characters` → `Extra` |
| `password.go` | `Characters` → `Extra` |
| `textarea.go` | `WhiteSpaces` → `Spaces`, `Characters` → `Extra` |
| `phone.go` | `Characters` → `Extra` |
| `number.go` | No field renames needed |
| `date.go` | `Characters` → `Extra`; update `Permitted.Validate` call |
| `hour.go` | `Characters` → `Extra`; update `Permitted.Validate` call |
| `ip.go` | `Characters` → `Extra`; update `Permitted.Validate` call |
| `rut.go` | `Characters` → `Extra`; update `Permitted.Validate` call |
| `address.go` | `WhiteSpaces` → `Spaces`, `Characters` → `Extra` |
| `checkbox.go` | Delete `c.SkipRules = true` |
| `datalist.go` | Delete `dl.SkipRules = true` |
| `select.go` | No field renames needed |
| `radio.go` | No field renames needed |
| `filepath.go` | `Characters` → `Extra`; update `Permitted.Validate` call |
| `gender.go` | No field renames needed |

### `form/input/` tests
| File | Change |
|---|---|
| `validation_test.go` | Update `Validate` calls and field names if needed |
| `inputs_test.go` | Update field names if needed |
| `render_test.go` | Update field names if needed |

### Documentation (already updated)
README.md and docs/ were updated in prior session — verify `Permitted` section references `fmt.Permitted` field names after execution.

---

## go.mod

No changes needed — `tinywasm/fmt` is already a dependency at v0.21.1 which includes `fmt.Permitted`.

---

## Verification

After execution:
1. `gotest ./input/...` — all input tests pass
2. `gotest ./...` — all form tests pass
3. `grep -r "input.Permitted" .` — no results (only `fmt.Permitted` used)
4. `permitted.go` no longer exists in `input/`
