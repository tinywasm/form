# PLAN: form — Field v3 Migration + ValidateData con action byte

← [README](../README.md) | Depende de: [fmt PLAN.md](../../fmt/docs/PLAN.md)

## Development Rules

- **Standard Library Only:** No external assertion libraries. Use `testing`.
- **Testing Runner:** Use `gotest` (install: `go install github.com/tinywasm/devflow/cmd/gotest@latest`).
- **Max 500 lines per file.** If exceeded, subdivide by domain.
- **Flat hierarchy.** No subdirectories for library code.
- **TinyGo Compatible:** No `fmt`, `strings`, `strconv`, `errors` from stdlib. Use `tinywasm/fmt`.
- **Documentation First:** Update docs before coding.


## Prerequisite

```bash
go get github.com/tinywasm/fmt@latest  # versión con ValidateFields(action, f)
```

## Contexto

Con Field v3:
- `Field.Input` ya no existe → form resuelve input type solo por heurística de nombre.
- `input.Permitted` se reemplaza por `fmt.Permitted` (sin maps, ASCII ranges).
- `ValidateData(action, data)` actualmente ignora `action` y valida per-input.
  Con `fmt.ValidateFields(action, data)` la validación es unificada y action-aware.

---

## Stage 1: `form.go` — Eliminar uso de `Field.Input`

**File:** `form.go`

### 1.1 Eliminar skip por `Field.Input`

```go
// ANTES (línea 84):
if field.Input == "-" {
    continue
}

// DESPUÉS: eliminar este bloque. Solo skip PK:
if field.PK {
    continue
}
```

### 1.2 Eliminar override explícito por `Field.Input`

```go
// ANTES (líneas 93-96):
var template input.Input
if field.Input != "" {
    template = findInputByType(field.Input)
}
if template == nil {
    template = findInputForField(fieldName, structName)
}

// DESPUÉS:
template := findInputForField(fieldName, structName)
```

### 1.3 Fallback a text cuando no hay match

```go
// ANTES (líneas 100-101):
if template == nil {
    return nil, fmt.Err("field", fieldName, "no matching input registered")
}

// DESPUÉS:
if template == nil {
    template = findInputByType("text")
}
```

### 1.4 Cambiar skip de PK

```go
// ANTES (líneas 78-81):
if field.PK && field.AutoInc {
    continue
}

// DESPUÉS:
if field.PK {
    continue
}
```

---

## Stage 2: `input/permitted.go` — Eliminar, reemplazar con `fmt.Permitted`

### 2.1 Eliminar `input/permitted.go`

El archivo `form/input/permitted.go` se elimina completo.
Toda la lógica de validación vive ahora en `fmt.Permitted`.

### 2.2 Actualizar `input/base.go` — embed `fmt.Permitted`

```go
// ANTES (línea 24):
Permitted      // anonymous embed

// DESPUÉS:
fmt.Permitted  // anonymous embed
```

### 2.3 Actualizar `Base.ValidateField`

`fmt.Permitted.Validate` toma `(field, text string)`:

```go
// ANTES (línea 129-131):
func (b *Base) ValidateField(value string) error {
    return b.Permitted.Validate(value)
}

// DESPUÉS:
func (b *Base) ValidateField(value string) error {
    return b.Permitted.Validate(b.FieldName(), value)
}
```

### 2.4 Actualizar inputs que llaman `Permitted.Validate` directamente

**Files:** `hour.go`, `ip.go`, `rut.go`, `date.go`, `filepath.go`

Patrón para cada uno:
```go
// ANTES:
if err := x.Permitted.Validate(value); err != nil {

// DESPUÉS:
if err := x.Permitted.Validate(x.FieldName(), value); err != nil {
```

### 2.5 Actualizar constructores de inputs

Cambiar `Permitted{...}` → `fmt.Permitted{...}` en todos los constructores.

**Files:** `email.go`, `phone.go`, `rut.go`, `ip.go`, `date.go`, `hour.go`,
`filepath.go`, `address.go`, `text.go`, `number.go`, `textarea.go`, etc.

---

## Stage 3: `validate_struct.go` — Delegar a `fmt.ValidateFields(action, data)`

**File:** `validate_struct.go`

```go
// ANTES:
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

// DESPUÉS:
func (f *Form) ValidateData(action byte, data fmt.Fielder) error {
    return fmt.ValidateFields(action, data)
}
```

---

## Stage 4: Actualizar tests

### 4.1 Actualizar Field literals en tests

Eliminar `Input:` de todos los `fmt.Field` en tests. Agregar `Permitted: fmt.Permitted{...}` donde corresponda.

### 4.2 Actualizar tests de `ValidateData`

Verificar comportamiento action-aware:
- `'c'` create: NotNull requerido, PK+AutoInc omitido, PK sin AutoInc requerido
- `'u'` update: PK requerido, NotNull requerido
- `'d'` delete: solo PK requerido

### 4.3 Actualizar tests de inputs

Verificar que `ValidateField` funciona con `fmt.Permitted` embebido en `Base`.

```bash
gotest
```

---

## Stage 5: Actualizar documentación

### 5.1 `docs/SKILL.md`

- Eliminar referencia a `Field.Input` en sección "Field Matching" (línea 29: punto 3)
- Actualizar `ValidateData` descripción: "delega a `fmt.ValidateFields(action, data)`"
- Eliminar mención de `input.Permitted` — ahora es `fmt.Permitted`

### 5.2 `docs/DESIGN.md`

- Eliminar punto 3 del "Registry" (línea 13): "Explicit override via `fmt.Field.Input`"
- Actualizar sección "Validation" (línea 21-22):

```
// ANTES:
- `ValidateData(action, data)` provides server-side or isomorphic validation.

// DESPUÉS:
- `ValidateData(action, data)` delegates to `fmt.ValidateFields(action, data)` — unified, action-aware validation.
```

### 5.3 `input/README.md`

Si referencia `Permitted` local, actualizar a `fmt.Permitted`.

---

## Stage 6: Limpiar planes obsoletos

Eliminar archivos que ya no aplican:
- `docs/CHECK_PLAN.md` — plan v3 original, reemplazado por este
- `docs/PLAN_UPDATE.md` — Fielder v2 migration, ya ejecutado
- `docs/PLANT.md` — plan v3 parcial, reemplazado por este
- `docs/PLAN_VALIDATE_ACTION.md` — fusionado en este plan

---

## Stage 7: Publish

```bash
gopush 'form: Field v3 + action byte — delete Permitted, remove Input, ValidateData usa fmt.ValidateFields(action, data)'
```

---

## Resumen

| Stage | File(s) | Acción |
|-------|---------|--------|
| 1 | `form.go` | Eliminar `Field.Input`, skip all PKs, fallback a text |
| 2 | `input/permitted.go` → delete, `input/base.go` + inputs | Reemplazar `input.Permitted` con `fmt.Permitted`, actualizar `Validate` signature |
| 3 | `validate_struct.go` | Reemplazar loop per-input con `fmt.ValidateFields(action, data)` |
| 4 | `*_test.go` | Actualizar literals, tests action-aware, tests de inputs |
| 5 | `docs/SKILL.md`, `docs/DESIGN.md`, `input/README.md` | Eliminar `Field.Input`, documentar validación unificada |
| 6 | `docs/` | Eliminar planes obsoletos |
| 7 | — | `gotest` + `gopush` |
