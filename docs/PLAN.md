---
PLAN: "feat: LoadValues — populate a form from a record (inverse of SyncValues)"
TAG: v0.2.14
---

> This plan is dispatched via the CodeJob workflow. See skill: agents-workflow.
> Fase **B** de la ola CRUD Harness:
> https://github.com/tinywasm/app/blob/main/docs/CRUD_HARNESS_MASTER_PLAN.md
> Gate: requiere `github.com/tinywasm/model v0.0.14` publicado.

# Plan — `Form.LoadValues`: rellenar el formulario desde un registro

## El problema

`form` sabe ir **del formulario al struct** (`sync.go`):

```go
func (f *Form) SyncValues(data model.Fielder) error // inputs → data.Pointers()
```

…pero **no sabe hacer el camino inverso**. Lo único disponible es campo a campo:

```go
func (f *Form) SetValues(fieldName string, values ...string) *Form
```

Eso obliga a cualquier consumidor con un CRUD ("el usuario pincha un registro de la lista →
el formulario se rellena con ese registro") a **convertir cada campo a string a mano**, uno
por uno, en cada módulo y en cada app. Es exactamente el boilerplate sin reflexión que este
ecosistema existe para eliminar, y es el hueco que hoy bloquea `layout/crudview`.

La lógica ya existe y está probada: **`New()` ya rellena los inputs desde el struct** al
construir el formulario (`form.go:116` → `model.ReadValues(schema, data.Pointers())`, y luego
`fmt.Convert(values[i]).String()` por input). Lo que falta es **exponerla como operación
repetible** sobre un formulario ya construido.

## Paso 1 — bump de dependencia

`go.mod`: `github.com/tinywasm/model` de `v0.0.12` a **`v0.0.14`**.

En v0.0.14, `model.Model` pasa a ser el contrato completo (`Fielder` + `ModuleNaming` +
`Encodable` + `Decodable`). **Este repo no cambia sus firmas por eso**: `form` legítimamente
solo necesita `model.Fielder` (esquema + punteros). No "actualices" `New`, `SyncValues`,
`ValidateData` ni `LoadValues` a `model.Model` — pedir más de lo que usas es lo contrario del
arnés. Deja `model.Fielder` en las cuatro.

## Paso 2 — `load.go` (archivo nuevo)

Crea `load.go` — hermano simétrico de `sync.go`:

```go
package form

import (
	"github.com/tinywasm/fmt"
	"github.com/tinywasm/model"
)

// LoadValues populates every input from data, the inverse of SyncValues.
// It is the operation a CRUD view needs when the user selects a record: one call,
// no per-field string conversion at the call site.
//
// A nil data (including a typed-nil pointer inside the interface) resets the form —
// that is the "new record" case, not an error.
func (f *Form) LoadValues(data model.Fielder) error {
	if model.IsNil(data) {
		f.reset()
		return nil
	}

	values := model.ReadValues(data.Schema(), data.Pointers())

	for i, inp := range f.Inputs {
		idx := f.fieldIndices[i]
		if idx < 0 || idx >= len(values) {
			continue
		}

		val := fmt.Convert(values[idx]).String()

		// Signal is the source of truth in WASM mode.
		f.valueSignals[i].Set(val)
		f.errorSignals[i].Set("") // loading a record clears stale validation errors

		// Keep input internal state in sync for SSR mode.
		if setter, ok := inp.(interface{ SetValues(...string) }); ok {
			setter.SetValues(val)
		}
	}

	return nil
}
```

Notas de implementación que **no** debes cambiar:

- `f.fieldIndices[i]` mapea el input `i` al índice de esquema. Es el mismo mapeo que usa
  `SyncValues` — reúsalo, no lo recalcules por nombre de campo.
- `model.IsNil(data)` (`model/codec.go`) cubre el puntero nil dentro de la interfaz, que un
  `data == nil` a secas **no** detecta. Es el caso real: la caché del módulo devuelve
  `(*ServiceItem)(nil)` cuando el ID no está.
- `reset()` (minúscula, `form.go`) es el método interno ya existente. No dupliques su lógica.
- Devuelve `error` aunque hoy no falle: simetría con `SyncValues` y espacio para validación
  futura sin romper a los consumidores.

## Paso 3 — `New` debe fallar RUIDOSAMENTE, no devolver un formulario vacío

Hoy hay una divergencia entre el doc y el código de `New` (`form.go:110-147`). El doc dice:

> *"Returns an error if any exported field has no matching registered input."*

…pero el código hace `continue // skip fields with no UI binding`. Resultado real, comprobado
con `veltylabs/service_catalog`: sus campos están declarados con `model.Text()` / `model.Int()`
(Kinds sin widget), que **no** satisfacen `input.Input` — así que `form.New(&ServiceItem{})`
devuelve un `*Form` **con cero inputs y `err == nil`**. Un formulario vacío que no falla.

Eso viola de frente el arnés de construcción
(https://github.com/tinywasm/app/blob/main/docs/CONSTRUCTION_HARNESS.md):
*"What the compiler can't catch becomes a loud development warning, never a silent failure."*

**Arreglo mínimo y sin ambigüedad:** al final de `New`, antes del `return`:

```go
if len(f.Inputs) == 0 {
	return nil, fmt.Errorf("form.New: %s has no renderable field — every Field.Type is a "+
		"plain model.Kind, not a form input.Input. Declare the widget in the model "+
		"Definition (input.Text(), input.Number(), …) instead of model.Text()/model.Int()",
		structName)
}
```

Usa `tinywasm/fmt`, y el mensaje **literal** de arriba: dice qué pasó, por qué, y cuál es la
línea que el consumidor tiene que cambiar.

**No conviertas el `continue` por campo en un error.** Un modelo tiene campos que
legítimamente no se pintan (`tenant_id`, `updated_at`, PKs autoinc — ya se saltan a
propósito). Distinguir "oculto a propósito" de "widget olvidado" campo a campo requiere una
marca nueva en el modelo, y eso es trabajo de `ormc`, **fuera del alcance de este plan**. Un
formulario con **cero** inputs, en cambio, nunca es intencional: ese es el caso que se cierra
aquí, y captura ruidosamente el fallo real que motivó esta ola.

## Paso 4 — tests (`load_test.go`)

Sin tag de build (lógica pura, sin DOM). Sobre un `Fielder` de prueba del repo:

1. **Rellena**: `New` → `LoadValues(&X{Name: "ACME", Price: 1500})` → los `valueSignals`
   correspondientes valen `"ACME"` y `"1500"`.
2. **Round-trip**: `LoadValues(a)` → `SyncValues(b)` → `b` es igual a `a` campo a campo. Es
   la garantía de que los dos sentidos son inversos de verdad.
3. **Reemplazo, no acumulación**: `LoadValues(a)` → `LoadValues(b)` → ningún campo conserva
   el valor de `a` (un campo vacío en `b` debe vaciar el input, no dejar el anterior).
4. **Nil resetea**: `LoadValues(a)` → `LoadValues(nil)` → todos los inputs vacíos.
5. **Nil tipado resetea**: `LoadValues((*X)(nil))` → todos los inputs vacíos, **sin panic**.
   Este test es el que justifica `model.IsNil`; si lo borras, el CRUD peta al deseleccionar.
6. **Limpia errores**: fuerza un error de validación en un input, `LoadValues(a)`, el
   `errorSignal` de ese input queda vacío.

7. **`New` con un modelo sin widgets falla**: un `Fielder` cuyos `Field.Type` sean todos
   `model.Text()`/`model.Int()` → `New` devuelve `err != nil` (no un form vacío).
8. **`New` con widgets funciona**: el mismo modelo con `input.Text()`/`input.Number()` →
   `err == nil` y `len(f.Inputs) > 0`.

## Paso 5 — documentación

`README.md`: en la sección de API, documenta el par simétrico junto con el ciclo CRUD real:

```go
f.LoadValues(record) // registro → formulario (el usuario selecciona)
f.Validate()         // el usuario edita y guarda
f.SyncValues(record) // formulario → registro (listo para enviar)
```

## Anti-footguns

- **Cero stdlib.** Este paquete compila a WASM: `tinywasm/fmt`, nunca `strconv`/`strings`/`errors`.
- **No toques `SyncValues`, `New`, `Reset` ni `SetValues`.** `LoadValues` se **añade**; no
  sustituye a nadie. `SetValues(fieldName, …)` sigue siendo válido para el caso de un solo campo.
- **No refactorices `New` para que llame a `LoadValues`.** Parece tentador (hacen algo
  parecido), pero `New` construye los inputs mientras lee, y `LoadValues` opera sobre inputs
  ya construidos. Fusionarlos mezcla dos ciclos de vida y rompe el orden de `fieldIndices`.
  Está considerado y descartado.
- **No añadas un `map`** de nombre de campo → input para acelerar la búsqueda. Regla del
  ecosistema: cero `map`. El recorrido por índice ya es O(n) sobre n pequeño.

## Criterios de aceptación

1. `Form.LoadValues(model.Fielder) error` existe en `load.go`.
2. `New` devuelve error cuando no vincula ni un solo input; el doc-comment de `New` y su
   código **ya no se contradicen**.
3. Los ocho tests pasan, incluido el del puntero nil tipado y el del modelo sin widgets.
4. `grep -rn "strconv\|strings\.\|\"errors\"" load.go` → vacío.
5. `New`, `SyncValues`, `ValidateData` siguen aceptando `model.Fielder` (no `model.Model`).
6. `gotest ./...` verde (stdlib + navegador).

## Tabla de etapas

| Etapa | Archivo(s) | Acción | Gate |
|---|---|---|---|
| B1 | `go.mod` | `model` → v0.0.14 | model publicado |
| B2 | `load.go` | `LoadValues` (inverso de `SyncValues`) | B1 |
| B3 | `form.go` | `New` falla ruidosamente con cero inputs vinculados | B1 |
| B4 | `load_test.go` | 8 tests: round-trip, nil tipado, modelo sin widgets | B2, B3 |
| B5 | `README.md` | documentar el ciclo `LoadValues`/`Validate`/`SyncValues` | B4 |
