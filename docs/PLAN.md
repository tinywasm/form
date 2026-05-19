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

Nota: el fix de `tinywasm/fmt` (`hasPermittedRules`/`validateLength`/`validateChars`)
ya fue publicado. Este plan solo cubre la capa `tinywasm/form`.

---

## 1. Error near-field en el DOM

### Decisión: span SSR + actualización via dom.Render en WASM

`RenderInput()` en `base.go` emite un `<span>` de error junto al input:

```html
<!-- Antes -->
<input type="text" id="app.nombre" name="nombre" ...>

<!-- Después -->
<input type="text" id="app.nombre" name="nombre" ...>
<span id="app.nombre.error" class="tw-field-error" aria-live="polite"></span>
```

- El span existe en el HTML inicial (SSR-compatible, sin layout shift).
- Empieza vacío. `min-height` en CSS reserva espacio para evitar saltos de layout.
- En WASM, `mount.go` lo actualiza usando `dom.Render` (única API disponible en `tinywasm/dom`
  para manipular el DOM; `Reference` solo expone `GetAttr`, `Value`, `Checked`, `On`, `Focus`).

### Nuevo método en `base.go` — `ErrorID()`

`Base` encapsula la convención del ID del span de error, igual que `GetID()` encapsula
`b.id`. Todos los consumidores (`RenderInput`, `mount.go`) usan este método en lugar de
concatenar strings manualmente:

```go
// ErrorID returns the ID of the associated error span for this input.
func (b *Base) ErrorID() string { return b.id + ".error" }
```

### Cambio en `base.go` — `RenderInput()`

Al final de cada rama (textarea e input), construcción tipada via el builder de `dom`
(`base.go` ya importa `"github.com/tinywasm/dom"` para `Children() []dom.Component`):

```go
errSpan := dom.Span("").ID(b.ErrorID()).Class("tw-field-error").Attr("aria-live", "polite")
out.Write(errSpan.RenderHTML())
```

Para radio/select/datalist: el span va después del contenedor del grupo.

---

## 2. Actualización del error en mount.go

```go
// Antes
inp.Validate(val)

// Después
errID := inp.ErrorID()
if err := inp.Validate(val); err != nil {
    dom.Render(errID, dom.Span(err.Error()).Class("tw-field-error", "tw-field-error--visible"))
} else {
    dom.Render(errID, dom.Span("").Class("tw-field-error"))
}
```

`dom.Render` reemplaza el contenido del elemento con el ID dado — es la API correcta
para actualizar nodos existentes sin métodos imperativos de DOM.

### En onSubmit: mostrar el primer error near-field

```go
if err := f.Validate(); err != nil {
    // mostrar error en el primer campo fallido — iteración completa va en v2
    return
}
```

En v1 basta con el primer error. Iteración completa de todos los campos queda para v2.

---

## 3. CSS de errores — fuente única de verdad vía tinywasm/css

### Principio

`tinywasm/form` consume tokens de `tinywasm/css`. No hardcodea colores.
Los proyectos sobreescriben el token `--color-error` en su propio `ssr.go`.

### Nuevo archivo: `form/ssr.go`

```go
//go:build !wasm

package form

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
func RenderCSS() *formCSS { return &formCSS{} }
```

`.tw-field-error` siempre `display:block` con `min-height` — reserva espacio para
evitar layout shift. `.tw-field-error--visible` solo añade peso visual; no cambia display.

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
`input.Text()` las permite por defecto (`Tilde: true` en su `Permitted`).
Campos técnicos (usernames, códigos) deben poder desactivarlas sin crear un widget nuevo.

### Sintaxis de tag (opt-out)

```go
// Tilde permitida por defecto — no se requiere ninguna etiqueta
Nombre string `input:"required,min=2"`

// Opt-out explícito para campos donde las tildes no aplican
Username string `input:"required,min=3,notilde"`
```

### Implementación

**En `tinywasm/form/input/text.go`** — método `SetTilde(bool)`:

```go
func (t *text) SetTilde(v bool) { t.Tilde = v }
```

No se expande la interfaz `Input`. El parser usa duck-typing para no acoplar la
interfaz base a una característica opt-out de un solo widget.

**En `tinywasm/form/validate_struct.go`** — parsear tag `notilde`:

```go
if hasTag(inputTag, "notilde") {
    if nt, ok := widget.(interface{ SetTilde(bool) }); ok {
        nt.SetTilde(false)
    }
}
```

`validate_struct.go` ya parsea otras etiquetas (`required`, `min`) — es el lugar correcto.
`SetTilde(false)` modifica la instancia clonada por campo, no el template global de `Text()`.

**Documentación (comentario en text.go)**:

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

### Problema actual

La firma `func(fmt.Fielder) error` es síncrona, pero `fetch.Send` en WASM es async.
El callback retorna `nil` antes de que llegue la respuesta del servidor. El form no
puede saber cuándo terminó la operación ni cuándo rehabilitar el botón.

### Nueva firma del callback — breaking change

```go
// Antes
f.OnSubmit(func(data fmt.Fielder) error { ... })

// Después
f.OnSubmit(func(data fmt.Fielder, done func(error)) {
    fetch.Post(apiURL).Send(func(resp *fetch.Response, err error) {
        done(err) // ← el form sabe que terminó
    })
})
```

`done(nil)` = éxito → form se resetea + botón rehabilitado.
`done(err)` = error → botón rehabilitado, form conserva datos para corrección.

### Estado del botón submit

El botón necesita un ID para ser direccionable via `dom.Render`.
En `render.go`, el botón pasa de:

```go
out.Write(`<button type="submit">`).Write(label).Write(`</button>`)
```

a:

```go
out.Write(`<button type="submit" id="`).Write(f.id).Write(`.submit">`).Write(label).Write(`</button>`)
```

**En `form.go`**: nuevo campo `submitLoadingLabel string`. Configurable con:

```go
f.SubmitLoadingLabel("Enviando...") // default vacío = usa submitLabel + "..."
```

**En `mount.go`** — flujo completo:

```go
onSubmit := func(e dom.Event) {
    e.PreventDefault()
    f.SyncValues(f.data)

    if err := f.Validate(); err != nil {
        // mostrar error near-field del primer campo fallido
        return
    }

    // Deshabilitar botón + mostrar label de carga
    submitID := f.GetID() + ".submit"
    loadingLabel := f.submitLoadingLabel
    if loadingLabel == "" {
        loadingLabel = f.resolveSubmitLabel() + "..."
    }
    dom.Render(submitID,
        dom.Button(loadingLabel).
            Attr("type", "submit").
            Attr("disabled", "true").
            ID(submitID),
    )

    if f.onSubmit != nil {
        f.onSubmit(f.data, func(err error) {
            if err == nil && !f.noResetOnSuccess {
                f.reset()
            }
            // Rehabilitar botón
            dom.Render(submitID,
                dom.Button(f.resolveSubmitLabel()).
                    Attr("type", "submit").
                    ID(submitID),
            )
        })
    }
}
```

### Reset del form

`f.reset()` (privado, llamado por `done(nil)`) + `f.Reset()` (público, para uso manual):

```go
func (f *Form) Reset() {
    // 1. Limpiar valores en inputs
    for _, inp := range f.Inputs {
        if setter, ok := inp.(interface{ SetValues(...string) }); ok {
            setter.SetValues("")
        }
        // 2. Limpiar span de error near-field
        errID := inp.ErrorID()
        dom.Render(errID, dom.Span("").Class("tw-field-error"))
    }
    // 3. Re-render de cada input con valor vacío se hace via dom.Render
    //    (requiere que cada input tenga ID propio — ya lo tiene via base.go)
}
```

### Configuración

```go
f.SubmitLoadingLabel("Enviando...")  // texto del botón mientras carga
f.NoResetOnSuccess()                 // opt-out del reset automático en éxito
f.Reset()                            // reset manual (público)
```

Nuevos campos en `Form`:
- `submitLoadingLabel string`
- `noResetOnSuccess bool`

---

## Archivos a modificar

| Archivo | Cambio |
|---------|--------|
| `form/input/base.go` — `RenderInput()` | Emitir `<span id="X.error" class="tw-field-error">` junto a cada input |
| `form/render.go` | Añadir `id="formID.submit"` al botón submit |
| `form/mount.go` — `OnMount()` | Callback async con `done func(error)`, loading state, reset via `dom.Render` |
| `form/form.go` | Nuevos campos: `submitLoadingLabel`, `noResetOnSuccess`; nuevo método `Reset()`, `SubmitLoadingLabel()`, `NoResetOnSuccess()` |
| `form/ssr.go` (nuevo) | Exportar `RenderCSS()` con `.tw-field-error` usando tokens de `tinywasm/css` |
| `form/input/text.go` | Agregar `SetTilde(bool)` + doc de la etiqueta `notilde` |
| `form/validate_struct.go` | Parsear tag `notilde` y llamar `SetTilde(false)` via type-assert |
| `form/go.mod` | Agregar dependencia `github.com/tinywasm/css` (build `!wasm` only) |

---

## Tests a agregar

### Instalación de gotest y wasmbrowsertest

`gotest` es la herramienta estándar del ecosistema tinywasm que ejecuta automáticamente:
vet, race detection, cobertura, **tests WASM en el navegador** (via `wasmbrowsertest`) y
badges de estado. Es el único comando necesario para validar la librería completa.

```bash
# Instalar gotest (una sola vez)
go install github.com/tinywasm/devflow/cmd/gotest@latest

# gotest instala wasmbrowsertest automáticamente si no está disponible:
# go install github.com/tinywasm/wasmbrowsertest@latest
```

Uso:
```bash
gotest              # suite completa: vet + race + cover + wasm + badges
gotest -no-cache    # forzar re-ejecución aunque el código no haya cambiado
gotest -run TestFoo # ejecutar un test específico
gotest -all         # incluye tests de integración, timeout 60s
```

### Patrón de tests en tinywasm/form

La librería usa tres archivos para cada grupo de tests, siguiendo el patrón ya establecido
(`base.back_test.go`, `base.front_test.go`, `base.shared_test.go`):

| Archivo | Build tag | Propósito |
|---------|-----------|-----------|
| `X.shared_test.go` | ninguno | lógica de test compartida (`runXTests(t)`) |
| `X.back_test.go` | `//go:build !wasm` | invoca `runXTests` en entorno nativo |
| `X.front_test.go` | `//go:build wasm` | invoca `runXTests` en el navegador via `wasmbrowsertest` |

`gotest` detecta automáticamente la presencia de archivos con `//go:build wasm` y lanza
`wasmbrowsertest` para ejecutar los tests del navegador sin configuración adicional.

### Tests nuevos a crear

**`render.shared_test.go`** (sin build tag) + `render.back_test.go` + `render.front_test.go`:

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
| `TestSubmit_CallbackReceivesDone` | La nueva firma `func(Fielder, func(error))` es invocada con `done` |
| `TestSubmit_NoResetOnSuccess` | `NoResetOnSuccess()` evita el reset al llamar `done(nil)` |
| `TestSubmit_ResetClearsValues` | `Reset()` vacía todos los inputs y spans de error |

---

## Orden de ejecución

1. `form/form.go`: nuevos campos + métodos (`SubmitLoadingLabel`, `NoResetOnSuccess`, `Reset`)
2. `form/render.go`: añadir ID al botón submit
3. `form/input/base.go`: agregar span de error en `RenderInput`
4. `form/ssr.go`: crear con `RenderCSS()` usando tokens css
5. `form/mount.go`: callback async con `done`, loading state, reset on success
6. `form/input/text.go`: agregar `SetTilde(bool)` + doc
7. `form/validate_struct.go`: parsear tag `notilde`
8. `form/go.mod`: agregar `github.com/tinywasm/css`
9. Agregar tests
10. Publicar via `gopush`
11. Actualizar `goflare-demo`: nueva firma de callback + `form.RenderCSS()` en `ssr.go`
