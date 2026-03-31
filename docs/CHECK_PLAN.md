# Plan: Prototype Pattern — Zero-Arg Constructors + True Clone + FieldDB Compat

## Depends on

- github.com/tinywasm/fmt v0.22.0 with FieldDB support 

## Problem

1. **Constructor mismatch**: All 17 input constructors require `(parentID, name string)` but `ormc` generates schema-level widgets where neither value is known: `Widget: input.Text()` fails to compile.
2. **Clone discards state**: `Clone()` calls the constructor from scratch (`return Email(parentID, name)`), losing any customization (Options, Placeholder, Min/Max overrides, Attributes) applied to the prototype.
3. **Aliases are dead code**: `Matches()` is never called. The old registry-based matching by alias is replaced by `Field.Widget` direct reference.

## Goal

Inputs become **stateless prototypes** at schema definition time. `Clone(parentID, name)` produces a **positioned copy** that preserves all prototype configuration.

## Changes

### 1. Remove aliases system from Base (dead code)

**In base.go**, remove:
- Field `aliases []string`
- Method `Matches(fieldName string) bool`
- Method `SetAliases(aliases ...string)`
- The `aliases ...string` variadic parameter from `InitBase` signature

`InitBase` changes from:
```go
func (b *Base) InitBase(parentID, name, htmlName string, aliases ...string)
```
To:
```go
func (b *Base) InitBase(parentID, name, htmlName string)
```

**In README.md**, remove alias references from the API table and matching algorithm section.

### 2. Constructors — remove `(parentID, name)` parameters + aliases

Every constructor changes from:

```go
func Text(parentID, name string) Input {
    t := &text{}
    // ... set defaults ...
    t.InitBase(parentID, name, "text", "name", "fullname", "username")
    return t
}
```

To:

```go
func Text() Input {
    t := &text{}
    // ... set defaults (unchanged) ...
    t.InitBase("", "", "text")
    return t
}
```

**Files** (17 constructors):
- text.go, email.go, password.go, number.go, phone.go, address.go
- checkbox.go, date.go, datalist.go, filepath.go, gender.go
- hour.go, ip.go, radio.go, rut.go, select.go, textarea.go

### 3. Clone — value copy instead of reconstruct

Every `Clone` changes from:

```go
func (e *email) Clone(parentID, name string) fmt.Widget {
    return Email(parentID, name)
}
```

To:

```go
func (e *email) Clone(parentID, name string) fmt.Widget {
    c := *e
    c.InitBase(parentID, name, "email")
    return &c
}
```

This preserves: Letters, Numbers, Extra, Minimum, Maximum, Placeholder, Title, Options, Attributes, Required, Disabled, Readonly, SkipValidation — everything set on the prototype.

**Same 17 files as above.**

### 4. Adapt form.go to FieldDB

Replace direct field access with helper methods from new fmt:

```go
// Before
if field.PK && field.AutoInc {
    continue
}

// After
if field.IsPK() && field.IsAutoInc() {
    continue
}
```

**Files**: form.go

### 5. Update tests

Update all test files:
- input/inputs_test.go, input/render_test.go, input/validation_test.go — constructor signature changes
- new_tests_test.go, setup_test.go — constructor changes + FieldDB struct literals

Constructor calls change:
- `input.Text("parent", "field")` → `input.Text().Clone("parent", "field").(input.Input)` or `input.Text()` for prototype-only tests

Schema literals in tests change:
```go
// Before
{Name: "id", Type: fmt.FieldText, PK: true, Widget: input.Text("", "")}

// After
{Name: "id", Type: fmt.FieldText, DB: &fmt.FieldDB{PK: true}, Widget: input.Text()}
```

### 6. Add Clone preservation test

New test: create a prototype with custom Options/Attributes, clone it, verify the clone has the same configuration.

## Execution Order

1. Wait for tinywasm/fmt with FieldDB to be published
2. Bump fmt dependency in go.mod
3. Remove aliases system from base.go
4. Modify all 17 constructors (zero-args)
5. Modify all 17 Clone methods (value copy)
6. Adapt form.go to use IsPK()/IsAutoInc()
7. Update all tests (constructors + FieldDB literals)
8. Update README.md
9. Run `go test ./...`

## Verification

- `go build ./...` passes
- `go test ./...` passes
- Clone preserves custom Options/Attributes (new test)
