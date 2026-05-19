# PLAN: tinywasm/form — Validación en tiempo real + errores near-field + submit lifecycle

## Contexto

El formulario de contacto en goflare-demo tiene tres problemas:
1. No hay feedback visual mientras el usuario escribe (no hay errores near-field).
2. El botón submit puede presionarse infinitas veces — no hay protección contra doble envío.
3. La firma del callback `func(fmt.Fielder) error` es incompatible con fetch async en WASM:
   el callback retorna `nil` antes de que el fetch complete, por lo que el form no puede
   saber cuándo rehabilitar el botón ni cuándo resetear los campos.

Este plan cubre:
1. Error near-field en el DOM (`base.go` + `mount.go`)
2. CSS de error usando tokens de `tinywasm/css`
3. Tag `notilde` para opt-out de tildes en `input.Text()`
4. Ciclo de vida del submit: callback async, loading state, reset por defecto

**Dependencias requeridas** (ya presentes en `go.mod`):
- `github.com/tinywasm/dom v0.9.4` — provee `Reference.SetValue`, `SetAttr`, `RemoveAttr`,
  `SetText` para mutación quirúrgica del DOM sin destruir event listeners.
- `github.com/tinywasm/fmt` — ya existente.

**Nota sobre `tinywasm/fmt`**: el fix de `hasPermittedRules`/`validateLength`/`validateChars`
ya fue publicado. Este plan solo cubre la capa `tinywasm/form`.

**Nota sobre CSS en WASM**: `form/ssr.go` usa `//go:build !wasm` — esto es CORRECTO e
INTENCIONAL. El CSS se emite como string durante el renderizado SSR y queda en el `<style>`
del HTML que carga el navegador. El código WASM no necesita inyectar CSS; ya está en el DOM.
No modificar este build tag.

---

## API de DOM disponible en v0.9.4

`dom.Get(id)` retorna `(Reference, bool)`. La interfaz `Reference` en v0.9.4 expone:

```go
type Reference interface {
    GetAttr(key string) string     // leer atributo
    Value() string                 // leer element.value
    SetValue(value string)         // escribir element.value — sin re-render
    SetAttr(key, value string)     // element.setAttribute — sin re-render
    RemoveAttr(key string)         // element.removeAttribute — sin re-render
    SetText(text string)           // element.textContent — sin re-render, sin XSS
    Checked() bool
    On(eventType string, handler func(Event))
    Focus()
}
```

**Regla**: usar siempre `dom.Get(id)` + métodos de `Reference` para mutar estado de
elementos existentes. Usar `dom.Render` solo para montar componentes nuevos, nunca para
actualizar estado (destruye event listeners via `cleanupChildren`).

---

## 1. Error near-field en el DOM

### Decisión: span SSR + mutación via Reference en WASM

`RenderInput()` en `base.go` emite un `<span>` de error junto al input:

```html
<!-- HTML emitido por RenderInput() -->
<input type="text" id="app.nombre" name="nombre" ...>
<span id="app.nombre.error" class="tw-field-error" aria-live="polite"></span>
```

- El span existe en el HTML inicial (SSR-compatible, sin layout shift).
- `min-height` en CSS reserva espacio para evitar saltos de layout.
- En WASM, `mount.go` muta su texto via `ref.SetText(msg)` — preserva `id`, `class`
  y `aria-live` sin destruir nada ni anidar elementos.

### Nuevo método en `base.go` — `ErrorID()`

```go
// ErrorID returns the ID of the associated error span for this input.
func (b *Base) ErrorID() string { return b.id + ".error" }
```

Todos los consumidores usan `inp.ErrorID()` — nunca concatenan strings manualmente.

### Cambio en `base.go` — `RenderInput()`

Al final de cada rama (textarea e input), construcción tipada via el builder de `dom`
(`base.go` ya importa `"github.com/tinywasm/dom"`):

```go
errSpan := dom.Span("").ID(b.ErrorID()).Class("tw-field-error").Attr("aria-live", "polite")
out.Write(errSpan.RenderHTML())
```

Para radio/select/datalist: el span va después del contenedor del grupo.

---

## 2. Actualización del error en mount.go

Usar `dom.Get` + `ref.SetText` — nunca `dom.Render` sobre el span de error:

```go
// Antes (incorrecto — anida <span> dentro del span existente)
// dom.Render(errID, dom.Span(msg).Class("tw-field-error--visible"))

// Después (correcto — muta textContent en el span existente, preserva id y aria-live)
errID := inp.ErrorID()
if ref, ok := dom.Get(errID); ok {
    if err := inp.Validate(val); err != nil {
        ref.SetText(err.Error())
        ref.SetAttr("class", "tw-field-error tw-field-error--visible")
    } else {
        ref.SetText("")
        ref.SetAttr("class", "tw-field-error")
    }
}
```

### En onSubmit: mostrar el primer error near-field

```go
if err := f.Validate(); err != nil {
    // mostrar error en el primer campo fallido — iteración completa va en v2
    return
}
```

---

## 3. CSS de errores — fuente única de verdad vía tinywasm/css

### Nuevo archivo: `form/ssr.go`

```go
//go:build !wasm

package form

// IMPORTANTE: este build tag es correcto. El CSS se emite en SSR y el navegador
// lo recibe en el <style> del HTML inicial. El código WASM no necesita CSS.

import "github.com/tinywasm/css"

type formCSS struct{}

func (f *formCSS) String() string {
    return `.tw-field-error {
  display: block;
  font-size: ` + css.TextSm.Var() + `;
  color: ` + css.ColorError.Var() + `;
  min-height: 1.2em;
  margin-top: ` + css.Space1.Var() + `;
}
.tw-field-error--visible {
  font-weight: ` + css.FontWeightMedium.Var() + `;
}`
}

// RenderCSS returns the default CSS for form validation errors.
// Call from the project's ssr.go to include this CSS in the page.
func RenderCSS() *formCSS { return &formCSS{} }
```

`.tw-field-error` siempre `display:block` con `min-height` — reserva espacio para
evitar layout shift. `.tw-field-error--visible` añade peso visual cuando hay error.

### Cadena completa

```
tinywasm/css   → --color-error: #E34F26  (token + fallback)
tinywasm/form  → .tw-field-error { color: var(--color-error, #E34F26) }
goflare-demo   → puede redefinir --color-error en su ssr.go (opcional)
```

---

## 4. Tag `notilde` — opt-out en input.Text()

### Motivación

En español las tildes son normativas: `María`, `Andrés`, `Diseño`.
`input.Text()` las permite por defecto (`Tilde: true`). Campos técnicos pueden desactivarlas.

### Sintaxis de tag (opt-out)

```go
Nombre   string `input:"required,min=2"`          // tildes permitidas por defecto
Username string `input:"required,min=3,notilde"`  // opt-out explícito
```

### Implementación

**`tinywasm/form/input/text.go`** — añadir `SetTilde(bool)`:

```go
func (t *text) SetTilde(v bool) { t.Tilde = v }
```

No se expande la interfaz `Input`. El parser usa duck-typing.

**`tinywasm/form/validate_struct.go`** — parsear tag `notilde`:

```go
if hasTag(inputTag, "notilde") {
    if nt, ok := widget.(interface{ SetTilde(bool) }); ok {
        nt.SetTilde(false)
    }
}
```

`SetTilde(false)` modifica la instancia clonada por campo, no el template global de `Text()`.

**Documentación en `text.go`**:

```go
// Text creates a standard text input.
// Tildes (á, é, í, ó, ú, Á, É, Í, Ó, Ú, Ñ) are allowed by default (natural for Spanish).
// To disallow tildes for a specific field, use the struct tag: `input:"...,notilde"`.
// Example:
//   Name   string `input:"required,min=2"`         // tildes allowed
//   UserID string `input:"required,min=3,notilde"` // tildes rejected
func Text() Input { ... }
```

---

## 5. Ciclo de vida del submit — callback async + loading + reset

### Nueva firma del callback — breaking change

```go
// Antes
f.OnSubmit(func(data fmt.Fielder) error { ... })

// Después
f.OnSubmit(func(data fmt.Fielder, done func(error)) {
    fetch.Post(apiURL).Send(func(resp *fetch.Response, err error) {
        done(err) // el form sabe cuándo terminó la operación async
    })
})
```

`done(nil)` = éxito → form se resetea + botón rehabilitado.
`done(err)` = error → botón rehabilitado, form conserva datos para corrección.

### Estado del botón submit

El botón necesita un ID. En `render.go`:

```go
// Antes
out.Write(`<button type="submit">`).Write(label).Write(`</button>`)

// Después
out.Write(`<button type="submit" id="`).Write(f.id).Write(`.submit">`).Write(label).Write(`</button>`)
```

**En `mount.go`** — usar `dom.Get` + `Reference` para mutar el botón sin destruir listeners:

```go
onSubmit := func(e dom.Event) {
    e.PreventDefault()
    f.SyncValues(f.data)

    if err := f.Validate(); err != nil {
        return
    }

    // Deshabilitar botón + mostrar label de carga via Reference (sin re-render)
    submitID := f.GetID() + ".submit"
    loadingLabel := f.submitLoadingLabel
    if loadingLabel == "" {
        loadingLabel = f.resolveSubmitLabel() + "..."
    }
    if btnRef, ok := dom.Get(submitID); ok {
        btnRef.SetAttr("disabled", "")
        btnRef.SetText(loadingLabel)
    }

    if f.onSubmit != nil {
        f.onSubmit(f.data, func(err error) {
            if err == nil && !f.noResetOnSuccess {
                f.reset()
            }
            // Rehabilitar botón via Reference (sin re-render)
            if btnRef, ok := dom.Get(submitID); ok {
                btnRef.RemoveAttr("disabled")
                btnRef.SetText(f.resolveSubmitLabel())
            }
        })
    }
}
```

### Reset del form

`f.reset()` usa `dom.Get` + `Reference` — nunca `dom.Render` sobre inputs existentes:

```go
func (f *Form) reset() {
    submitID := f.GetID() + ".submit"
    for _, inp := range f.Inputs {
        // Limpiar valor via SetValue — preserva el listener "input" del form
        if ref, ok := dom.Get(inp.GetID()); ok {
            ref.SetValue("")
        }
        // Limpiar span de error via SetText — preserva id y aria-live
        if ref, ok := dom.Get(inp.ErrorID()); ok {
            ref.SetText("")
            ref.SetAttr("class", "tw-field-error")
        }
        // Limpiar valor en el struct interno
        if setter, ok := inp.(interface{ SetValues(...string) }); ok {
            setter.SetValues("")
        }
    }
    _ = submitID // el botón ya fue rehabilitado por el caller (done callback)
}

// Reset es la versión pública para uso manual desde el caller.
func (f *Form) Reset() { f.reset() }
```

### Configuración

```go
f.SubmitLoadingLabel("Enviando...") // texto del botón mientras carga (default: label + "...")
f.NoResetOnSuccess()                // opt-out del reset automático en éxito
f.Reset()                           // reset manual (público)
```

Nuevos campos en `Form`: `submitLoadingLabel string`, `noResetOnSuccess bool`.

---

## Archivos a modificar

| Archivo | Cambio |
|---------|--------|
| `form/input/base.go` | `ErrorID()` + span de error tipado en `RenderInput()` |
| `form/render.go` | Añadir `id="formID.submit"` al botón submit |
| `form/mount.go` | Callback async `done`, loading via `Reference`, error near-field via `Reference` |
| `form/form.go` | Campos `submitLoadingLabel`, `noResetOnSuccess`; métodos `Reset()`, `SubmitLoadingLabel()`, `NoResetOnSuccess()`, `reset()` |
| `form/ssr.go` (nuevo) | `RenderCSS()` con `//go:build !wasm` — correcto e intencional |
| `form/input/text.go` | `SetTilde(bool)` + doc de la etiqueta `notilde` |
| `form/validate_struct.go` | Parsear tag `notilde`, llamar `SetTilde(false)` via type-assert |
| `form/go.mod` | Agregar `github.com/tinywasm/css` si no está |

---

## Tests a agregar

### Instalación de gotest y wasmbrowsertest

```bash
# Instalar gotest (una sola vez)
go install github.com/tinywasm/devflow/cmd/gotest@latest

# gotest instala wasmbrowsertest automáticamente si no está disponible.
# También puede instalarse manualmente:
# go install github.com/tinywasm/wasmbrowsertest@latest
```

Uso:
```bash
gotest              # suite completa: vet + race + cover + wasm en navegador + badges
gotest -no-cache    # forzar re-ejecución aunque el código no haya cambiado
gotest -run TestFoo # ejecutar un test específico
```

### Patrón de tests en tinywasm/form

| Archivo | Build tag | Propósito |
|---------|-----------|-----------|
| `X.shared_test.go` | ninguno | lógica de test compartida (`runXTests(t)`) |
| `X.back_test.go` | `//go:build !wasm` | invoca `runXTests` en entorno nativo |
| `X.front_test.go` | `//go:build wasm` | invoca `runXTests` en navegador via `wasmbrowsertest` |

`gotest` detecta `//go:build wasm` automáticamente y lanza `wasmbrowsertest`.

### Tests nuevos a crear

**`render.shared_test.go`** + back + front:

| Test | Verifica |
|------|----------|
| `TestRenderInput_EmitsErrorSpan` | `RenderHTML()` contiene `id="X.error"` y `class="tw-field-error"` |
| `TestRender_SubmitButtonHasID` | Botón submit tiene `id="formID.submit"` |
| `TestRender_ErrorIDMethod` | `inp.ErrorID()` retorna `inp.GetID() + ".error"` |

**`notilde.shared_test.go`** + back + front:

| Test | Verifica |
|------|----------|
| `TestNotilde_RejectsAccent` | Campo con tag `notilde` rechaza `á` |
| `TestNotilde_AllowsNormal` | Campo con tag `notilde` acepta `a` |
| `TestText_AllowsAccentByDefault` | Campo sin tag acepta `á` |

**`submit.shared_test.go`** + back + front:

| Test | Verifica |
|------|----------|
| `TestSubmit_CallbackReceivesDone` | Nueva firma `func(Fielder, func(error))` invocada con `done` |
| `TestSubmit_NoResetOnSuccess` | `NoResetOnSuccess()` evita reset al llamar `done(nil)` |
| `TestSubmit_ResetClearsValues` | `Reset()` vacía inputs y spans de error |

---

## Orden de ejecución

1. `form/form.go`: campos + métodos nuevos (`SubmitLoadingLabel`, `NoResetOnSuccess`, `reset`, `Reset`)
2. `form/render.go`: añadir ID al botón submit
3. `form/input/base.go`: `ErrorID()` + span de error en `RenderInput()`
4. `form/ssr.go`: crear con `RenderCSS()` + `//go:build !wasm`
5. `form/mount.go`: callback async, loading via `Reference`, error near-field via `Reference`
6. `form/input/text.go`: `SetTilde(bool)` + doc
7. `form/validate_struct.go`: parsear tag `notilde`
8. `form/go.mod`: agregar `github.com/tinywasm/css` si no está
9. Agregar tests siguiendo el patrón shared/back/front
10. Ejecutar `gotest` — suite completa debe estar verde
11. Publicar via `gopush`
12. Actualizar `goflare-demo`: nueva firma de callback + `form.RenderCSS()` en `ssr.go`
