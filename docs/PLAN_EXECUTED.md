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

---

# form — PLAN: migrar a `widget` v0.4.0

**Ejecutado** en [#19](https://github.com/tinywasm/form/pull/19). Se registra aquí
porque el resto de la suite sí necesitó un plan de migración y `form` no: sin este
apunte no se distingue "no hacía falta trabajo" de "nadie lo revisó".

## Contexto

`tinywasm/widget` v0.4.0 es una release cerrada y rompedora. Rompió los ocho
paquetes de `tinywasm/components` y los tres de `tinywasm/layout`.

**No rompió `form`.**

## Por qué

`form` **nunca importa `widget/style`** — el paquete que cambió. Usa solo el
paquete raíz `widget`, y solo estos símbolos:

| Símbolo | Usos | ¿Cambió en v0.4.0? |
|---|---|---|
| `widget.NameField` | 7 | no |
| `widget.PartInput` | 3 | no |
| `widget.PartLabel`, `PartError`, `PartRadioGroup` | 1 c/u | no |
| `widget.State`, `Locked`, `Invalid` | 1 c/u | no |

Toda la superficie rompedora de v0.4.0 vive dentro de `widget/style`: las escalas,
las superficies, los primitivos de flujo y el constructor de hojas. El contrato de
identidad del paquete raíz (`Name`, `Part`, `Class`, `State`, `Cue`) se dejó
intacto precisamente para que consumidores como este no tuvieran que moverse.

## Verificación

Comprobado, no supuesto, tras `go get github.com/tinywasm/widget@v0.4.0 && go mod tidy`:

```bash
go build ./...   # limpio, sin editar una sola línea de código
go test ./...    # ok  github.com/tinywasm/form
                 # ok  github.com/tinywasm/form/input
                 # ok  github.com/tinywasm/form/tests
```

El cambio completo fueron `go.mod` y `go.sum`.

## Invariantes

- `form` NO importa `widget/style`. Si algún día lo hiciera, deja de estar exento
  de las migraciones de esa API y necesitará su propio plan.
- Los `Part*` y `NameField` que consume son del paquete raíz, compartidos con
  `components/fieldset`, que es quien los estiliza.

## Estado de la suite tras esta migración

| Repo | Necesitó plan | Estado |
|---|---|---|
| `css` | sí | ejecutado y en `main` |
| `widget` | sí | publicado v0.4.0 |
| `form` | **no** | este apunte |
| `ssr` | sí | plan en `main`, pendiente de ejecutar |
| `components` | sí | plan en PR, 8 paquetes, 1 revienta en runtime |
| `layout` | sí | plan en PR, bloqueado por `components` |
