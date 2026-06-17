# form — PLAN: Eliminar dependencia i18n del core (opt-in real)

## Contexto

`tinywasm/fmt` v0.24.1 movió `Translate`, `RegisterWords`, `DictEntry` a `fmt/lang` (opt-in).
`form` actualmente:
1. Llama `fmt.Translate(...)` en `input/base.go`, `render.go`, `form.go` — ya no compila.
2. Tiene un `init()` en `words.go` que llama `fmt.RegisterWords(...)` — arrastra el
   diccionario completo a cualquier binario que importe `form`.

**No** se debe importar `fmt/lang` desde `form`: eso volvería a arrastrar el diccionario
automáticamente, rompiendo el opt-in que el refactor de `fmt` acaba de lograr.

## Solución

Eliminar la dependencia i18n de `form`. Las cadenas UI (`"Submit"`, `"Optional"`, nombres de
campos) se emiten como texto crudo. Quien quiera traducción importa `fmt/lang` en su app
(o en un paquete de words propio) — esto es consistente con el principio opt-in.

## Archivos a cambiar (solo en `form/`)

### `input/base.go`

```go
// Antes:
b.Placeholder = fmt.Translate(name).String()
b.Title       = fmt.Translate(name).String()

// Después (texto crudo — la app puede traducir en su capa de presentación):
b.Placeholder = name
b.Title       = name
```

### `render.go` y `form.go`

```go
// Antes:
label = fmt.Translate("Submit").String()

// Después:
label = "Submit"
```

Actualizar también los comentarios en `form.go` que mencionan `Translate("Submit")`:
- línea ~20: `// Submit button label (empty = use Translate("Submit"))` → `// Submit button label (empty = "Submit")`
- línea ~67: `// If never called, the button shows Translate("Submit") (locale-aware).` → `// If never called, the button shows "Submit".`

### `words.go` — ELIMINAR el archivo

El `init()` que registra palabras en el diccionario global arrastra i18n sin que nadie lo pida.
Eliminar `words.go` por completo. Si una app necesita traducir "Submit"/"Optional" al idioma del
usuario, lo hace en su propia capa (importando `fmt/lang` explícitamente).

### `base.shared_test.go`

El comentario de línea ~65 dice `// Placeholder and title default to the translated field name via fmt.Translate.`
Actualizar a: `// Placeholder and title default to the raw field name.`
El comportamiento esperado (`"Name"`) no cambia — el test pasa sin modificar asserts.

## Verificación

```bash
# Sin referencias a Translate/RegisterWords/DictEntry en el core de form (excl. tests):
grep -rn 'Translate\|RegisterWords\|DictEntry' *.go input/*.go | grep -v _test

# Tests verdes:
gotest
```

## Invariantes

- `form` NO importa `fmt/lang`. Solo importa el root `fmt`.
- La traducción de etiquetas de UI es responsabilidad del consumidor (app/componente), no de `form`.
- No cambiar la API pública de `form` — solo eliminar la traducción automática silenciosa.
