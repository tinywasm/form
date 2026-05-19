# PLAN: tinywasm/form — Validación en tiempo real + errores near-field + tag `notilde`

## Contexto

El formulario de contacto en goflare-demo muestra un 422 silencioso al enviar datos válidos
con tildes. Además, no hay feedback visual mientras el usuario escribe. Este plan cubre:

1. Fix de validación en tiempo real (mount.go)
2. Error near-field en el DOM (render.go + base.go)
3. CSS de error usando tokens de tinywasm/css
4. Tag `notilde` para opt-out de tildes en `input.Text()`

---

## 1. Error near-field en el DOM

### Decisión: opción A (render SSR + actualización WASM)

`RenderInput()` en `base.go` emite un `<span>` de error junto al input:

```html
<!-- Antes -->
<input type="text" id="app.nombre" name="nombre" ...>

<!-- Después -->
<input type="text" id="app.nombre" name="nombre" ...>
<span id="app.nombre.error" class="tw-field-error" aria-live="polite"></span>
```

- El span existe en el HTML inicial (SSR-compatible).
- Empieza vacío (`aria-live="polite"` para accesibilidad).
- En WASM, `mount.go` lo llena/vacía según el resultado de `Validate`.

### Cambio en `base.go` — `RenderInput()`

Al final de cada rama de renderizado (textarea, input), agregar:

```go
out.Write(`<span id="`).Write(b.id).Write(`.error" class="tw-field-error" aria-live="polite"></span>`)
```

Para radio/select/datalist: el span va después del contenedor del grupo, no dentro.

---

## 2. Actualización del error en mount.go

```go
// Antes
inp.Validate(val)

// Después
errID := inp.GetID() + ".error"
if err := inp.Validate(val); err != nil {
    dom.SetText(errID, err.Error())
    dom.AddClass(errID, "tw-field-error--visible")
} else {
    dom.SetText(errID, "")
    dom.RemoveClass(errID, "tw-field-error--visible")
}
```

`dom.SetText` y `dom.AddClass`/`dom.RemoveClass` ya existen en `tinywasm/dom`.
Si no existen, se usan los métodos disponibles equivalentes.

### En onSubmit: mostrar todos los errores

```go
if err := f.Validate(); err != nil {
    // Mostrar el primer error near-field (Validate ya retorna el primer campo fallido)
    // La iteración completa es opcional en v1
    return
}
```

En v1 basta con mostrar el primer error. Una iteración completa puede venir en v2.

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
  display: block;
}`
}

// RenderCSS returns the default CSS for form validation errors.
func RenderCSS() *formCSS { return &formCSS{} }
```

### Cadena completa

```
tinywasm/css     → --color-error: #E34F26  (token + fallback)
tinywasm/form    → .tw-field-error { color: var(--color-error, #E34F26) }
goflare-demo     → puede redefinir --color-error en su ssr.go (opcional)
                   resultado: el color de error cambia automáticamente
```

---

## 4. Tag `notilde` — opt-out en input.Text()

### Motivación

En español las tildes son normativas: `María`, `Andrés`, `Diseño`.
`input.Text()` las permite por defecto (`Tilde: true` en su `Permitted`).
Campos técnicos (usernames, códigos) deben poder desactivarlas sin cambiar el widget.

### Sintaxis de tag (opt-out)

```go
// Tilde permitida por defecto — no se requiere ninguna etiqueta
Nombre string `input:"required,min=2"`

// Opt-out explícito para campos donde las tildes no aplican
Username string `input:"required,min=3,notilde"`
```

### Implementación

**En `tinywasm/form/input/text.go`**: exponer método `SetTilde(bool)`:

```go
func (t *text) SetTilde(v bool) { t.Tilde = v }
```

No se expande la interfaz `Input`. El parser usa duck-typing via type-assert.

**En `tinywasm/form/registry.go` o `validate_struct.go`**: al registrar campos desde struct,
parsear la etiqueta `input:"..."`:

```go
tags := strings.Split(inputTag, ",")
for _, tag := range tags {
    if tag == "notilde" {
        if nt, ok := widget.(interface{ SetTilde(bool) }); ok {
            nt.SetTilde(false)
        }
        break
    }
}
```

Esto permite que `input.Text()` desactive tildes sin necesidad de crear un widget nuevo.

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

### Relación con tinywasm/fmt

`tinywasm/fmt.Permitted` expone el campo `Tilde bool`. El widget `text` lo usa directamente.
Cuando `form/registry` parsea `notilde`, llama `widget.SetTilde(false)` que modifica el
`Permitted` embebido.

---

## Archivos a modificar

| Archivo | Cambio |
|---------|--------|
| `form/input/base.go` — `RenderInput()` | Emitir `<span id="X.error" class="tw-field-error">` junto a cada input |
| `form/mount.go` — `OnMount()` | Capturar error de `Validate` y actualizar el span near-field |
| `form/ssr.go` (nuevo) | Exportar `RenderCSS()` con `.tw-field-error` usando tokens de `tinywasm/css` |
| `form/input/text.go` | Agregar `SetTilde(bool)` |
| `form/registry.go` o `validate_struct.go` | Parsear tag `notilde` y llamar `SetTilde(false)` |
| `form/go.mod` | Agregar dependencia `github.com/tinywasm/css` |

---

## Tests a agregar

| Test | Verifica |
|------|----------|
| `render_test.go` | `RenderInput()` emite el span `id="X.error"` |
| `notilde_test.go` | Campo con `notilde` rechaza `á`; campo sin tag acepta `á` |
| `mount_test.go` (si hay harness WASM) | Error near-field aparece al escribir carácter inválido |

---

## Orden de ejecución

1. `tinywasm/fmt`: fix `hasPermittedRules` + `validateLength`/`validateChars` → publicar
2. `form/input/base.go`: agregar span de error en `RenderInput`
3. `form/ssr.go`: crear con `RenderCSS()` usando tokens css
4. `form/mount.go`: actualizar span on-input
5. `form/input/text.go`: agregar `SetTilde(bool)`
6. `form/registry.go`: parsear tag `notilde`
7. Agregar tests
8. Publicar via `gopush`
9. Actualizar `goflare-demo`: agregar `RenderCSS()` del form en su `ssr.go`
