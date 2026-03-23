# PLAN: Fielder v2 Migration (tinywasm/form)

← [README](../README.md) | Depends on: [fmt PLAN_FIELDER_V2](../../fmt/docs/PLAN_FIELDER_V2.md)

## Development Rules

- **Standard Library Only:** No external assertion libraries. Use `testing`.
- **Testing Runner:** Use `gotest` (install: `go install github.com/tinywasm/devflow/cmd/gotest@latest`).
- **Max 500 lines per file.** If exceeded, subdivide by domain.
- **Flat hierarchy.** No subdirectories for library code.
- **TinyGo Compatible:** No `fmt`, `strings`, `strconv`, `errors` from stdlib. Use `tinywasm/fmt`.
- **Documentation First:** Update docs before coding.
- **Publishing:** Use `gopush 'message'` after tests pass and docs are updated.

## Prerequisite

Update `go.mod` to the new `tinywasm/fmt` version (Fielder v2, without `Values()`):

```bash
go get github.com/tinywasm/fmt@v0.19.0
```

## Context

The `form` library uses `Values()` in exactly 2 places:

1. **`form.go:61`** — `New()`: reads initial field values to bind them to form inputs.
2. **`validate_struct.go:10`** — `ValidateData()`: reads values for validation.

Both can read through `Pointers()` instead, using `fmt.Convert` to get the string
representation. `SyncValues()` in `sync.go` already uses only `Pointers()` + `Schema()`.

---

## Stage 1: Update `form.go` — `New()` function

**File:** `form.go`

### 1.1 Replace `values := data.Values()` with pointer-based reading

```go
// BEFORE (lines 60-61):
schema := data.Schema()
values := data.Values()

// AFTER:
schema := data.Schema()
ptrs := data.Pointers()
```

### 1.2 Update value binding (line 114)

```go
// BEFORE (line 113-115):
if setter, ok := inp.(interface{ SetValues(...string) }); ok {
    setter.SetValues(fmt.Convert(values[i]).String())
}

// AFTER:
if setter, ok := inp.(interface{ SetValues(...string) }); ok {
    setter.SetValues(readPtrAsString(ptrs[i], field.Type))
}
```

### 1.3 Add `readPtrAsString` helper

Add at bottom of `form.go` (or in a new `helpers.go` if form.go is near 500 lines):

```go
// readPtrAsString reads a value from a typed pointer and returns its string representation.
// Used by New() and ValidateData() to get field values without Values().
func readPtrAsString(ptr any, ft fmt.FieldType) string {
	switch ft {
	case fmt.FieldText:
		if p, ok := ptr.(*string); ok {
			return *p
		}
	case fmt.FieldInt:
		switch p := ptr.(type) {
		case *int64:
			return fmt.Convert(*p).String()
		case *int:
			return fmt.Convert(*p).String()
		case *int32:
			return fmt.Convert(*p).String()
		}
	case fmt.FieldFloat:
		switch p := ptr.(type) {
		case *float64:
			return fmt.Convert(*p).String()
		case *float32:
			return fmt.Convert(*p).String()
		}
	case fmt.FieldBool:
		if p, ok := ptr.(*bool); ok {
			if *p {
				return "true"
			}
			return "false"
		}
	}
	return ""
}
```

**Note:** For `FieldInt` and `FieldFloat`, `fmt.Convert(*p).String()` uses the pooled `Conv`,
so it's 0 extra allocs in backend (sync.Pool) and minimal in WASM.

---

## Stage 2: Update `validate_struct.go` — `ValidateData()` function

**File:** `validate_struct.go`

```go
// BEFORE:
func (f *Form) ValidateData(action byte, data fmt.Fielder) error {
	values := data.Values()
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
	schema := data.Schema()
	ptrs := data.Pointers()
	for i, inp := range f.Inputs {
		idx := f.fieldIndices[i]
		if idx < 0 || idx >= len(ptrs) {
			continue
		}
		if skipper, ok := inp.(interface{ GetSkipValidation() bool }); ok && skipper.GetSkipValidation() {
			continue
		}
		val := readPtrAsString(ptrs[idx], schema[idx].Type)
		if err := inp.ValidateField(val); err != nil {
			return err
		}
	}
	return nil
}
```

---

## Stage 3: Update tests

### 3.1 Remove `Values()` from test mocks

```bash
grep -rn "func.*Values().*\[\]any" *_test.go
```

Remove every `Values()` method from mock structs in test files.

### 3.2 Run tests

```bash
gotest
```

---

## Stage 4: Update documentation

### 4.1 Update `docs/SKILL.md` if it references `Values()`

```bash
grep -n "Values" docs/SKILL.md
```

### 4.2 Update `docs/API.md` or similar docs that show the Fielder interface

---

## Stage 5: Publish

```bash
gopush 'form: Fielder v2 migration — read through Pointers instead of Values'
```

---

## Summary

| Stage | File(s) | Action |
|-------|---------|--------|
| 1 | `form.go` | Replace `data.Values()` with `data.Pointers()` + `readPtrAsString` |
| 2 | `validate_struct.go` | Replace `data.Values()` with pointer-based reading |
| 3 | `*_test.go` | Remove Values() from mocks |
| 4 | `docs/` | Update references to Values() |
| 5 | — | `gotest` + `gopush` |
