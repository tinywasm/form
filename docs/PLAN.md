# PLAN: tinywasm/form — Migración a html + rename String()

## Repositorio
`github.com/tinywasm/form` — path local: `tinywasm/form/`

## Dependencias de ejecución
```bash
go install github.com/tinywasm/devflow/cmd/gotest@latest
```

## Prerequisitos (ejecutar ANTES de este plan)
1. `tinywasm/dom` publicado con `String()` en lugar de `RenderHTML()`, builders eliminados
2. `tinywasm/html` publicado con `Span`, `Div`, `Label`, `Button`, etc.

---

## Objetivo

Migrar `tinywasm/form` para:
1. Usar `tinywasm/html` para el único builder que usa: `dom.Span()`
2. Renombrar `RenderHTML() string` → `String() string` en `Form` y `Base`
3. Actualizar la interface `input.Input` (que embede `dom.Component`)
4. Actualizar todos los tests que llaman `.RenderHTML()`

Es un **break change limpio** — sin aliases, sin compatibilidad hacia atrás.

---

## Actualizar go.mod

```
module github.com/tinywasm/form

go 1.25

require (
    github.com/tinywasm/fmt  v<version-actual>
    github.com/tinywasm/dom  v<nueva-version>   // aún necesario: Component, Event, etc.
    github.com/tinywasm/css  v<version-actual>
    github.com/tinywasm/html v<nueva-version>   // NUEVO: Span, Label, etc.
)
```

---

## Cambio 1: `input/base.go` — Span builder + rename RenderHTML

### Imports — agregar html

**Buscar en `input/base.go`:**
```go
import (
    "github.com/tinywasm/dom"
```

**Reemplazar con:**
```go
import (
    . "github.com/tinywasm/html"  // Span, Label, etc.
    "github.com/tinywasm/dom"
```

### Reemplazar builder

**Buscar (3 ocurrencias en líneas ~234, ~262, ~289):**
```go
errSpan := dom.Span("").ID(b.ErrorID()).Class("tw-field-error").Attr("aria-live", "polite")
```

**Reemplazar con:**
```go
errSpan := Span("").ID(b.ErrorID()).Class("tw-field-error").Attr("aria-live", "polite")
```

### Rename RenderHTML → String en Base

**Buscar en `input/base.go`:**
```go
// RenderHTML renders the input to HTML.
func (b *Base) RenderHTML() string {
```

**Reemplazar con:**
```go
// String serializes the input to its HTML string representation.
func (b *Base) String() string {
```

También actualizar las llamadas internas a `errSpan.RenderHTML()`:

**Buscar (3 ocurrencias):**
```go
out.Write(errSpan.RenderHTML())
```

**Reemplazar con:**
```go
out.Write(errSpan.String())
```

---

## Cambio 2: `render.go` — rename RenderHTML + actualizar llamadas

### Rename método de Form

**Buscar en `render.go`:**
```go
// RenderHTML renders the form based on its SSR mode.
func (f *Form) RenderHTML() string {
```

**Reemplazar con:**
```go
// String serializes the form to its HTML string representation.
func (f *Form) String() string {
```

### Actualizar llamada interna a inputs

**Buscar en `render.go`:**
```go
out.Write(inp.RenderHTML())
```

**Reemplazar con:**
```go
out.Write(inp.String())
```

---

## Cambio 3: `input/interface.go` — Input interface

La interface `Input` embebe `dom.Component`. Después del rename en dom, `dom.Component` expone `String() string` en lugar de `RenderHTML() string`. No se necesita cambiar la interface explícitamente — se hereda automáticamente del embed.

Verificar que el archivo compile sin cambios adicionales:
```bash
cd tinywasm/form
go build ./input/...
```

---

## Cambio 4: Tests

### `render.shared_test.go`

**Buscar y reemplazar todas las ocurrencias:**
```go
// ANTES:
html := f.RenderHTML()

// DESPUÉS:
html := f.String()
```

También actualizar nombre de test si contiene "RenderHTML" en el nombre (opcional — no afecta compilación):
```go
// ANTES:
func TestForm_RenderHTML_SSR_Shared(t *testing.T) {
// DESPUÉS (opcional):
func TestForm_String_SSR_Shared(t *testing.T) {
```

Y en `base.shared_test.go`:
```go
// ANTES:
html := f.RenderHTML()
t.Run("RenderHTML_SSR", TestForm_RenderHTML_SSR_Shared)

// DESPUÉS:
html := f.String()
t.Run("String_SSR", TestForm_String_SSR_Shared)
```

### `input/render_test.go`

**Buscar y reemplazar:**
```go
// ANTES:
html := inp.RenderHTML()
t.Errorf("RenderHTML() missing %q\ngot: %s", c.contain, html)

// DESPUÉS:
html := inp.String()
t.Errorf("String() missing %q\ngot: %s", c.contain, html)
```

### `input/inputs_test.go`

**Buscar y reemplazar:**
```go
// ANTES:
html := cloned.RenderHTML()

// DESPUÉS:
html := cloned.String()
```

---

## Documentación a Actualizar

### `form/README.md`

Actualizar cualquier ejemplo de código que muestre `.RenderHTML()` → `.String()`.

### `form/input/README.md`

Actualizar ejemplos de `RenderHTML()` → `String()`. Líneas que referencian `dom.Span` → `html.Span`.

### `form/docs/`

Revisar `API.md`, `DESIGN.md`, `IMPLEMENTATION.md` — actualizar referencias a `RenderHTML()` → `String()` y al import de `dom` para builders → `html`.

---

## Verificación

```bash
cd tinywasm/form
go build ./...
gotest
```

Todos los tests deben pasar. Si algún test falla por `RenderHTML`, es migración incompleta.

Ver `tinywasm/docs/MASTER_PLAN.md` para el orden global de ejecución.
