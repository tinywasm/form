---
PLAN: "feat: DirtyFields reports which fields the user actually changed"
TAG: v0.4.0
EXECUTOR: jules
REVIEWER: none
STATUS: review
SESSION: 2585114900977056736
PR: https://github.com/tinywasm/form/pull/22
---

> Este plan se despacha con el flujo CodeJob. Ver skill: agents-workflow.
>
> Forma parte de una ola: `docs/BULK_ACTIONS_MASTER_PLAN.md` en la raíz del
> monorepo. Este plan es **independiente** — no espera a ningún otro. Es
> **puerta** para el plan de `layout`.

# Plan — `DirtyFields`: qué campos tocó el usuario

## 0. Prerrequisito

```bash
go install github.com/tinywasm/devflow/cmd/gotest@latest
```

Los tests se ejecutan con `gotest`. **Nunca `go test`.**

## 1. Por qué

Se va a añadir edición masiva: el usuario marca varios registros, abre el
formulario, toca **un** campo (por ejemplo la fecha) y ese valor se aplica a
todos los marcados. Para eso hay que saber **cuáles** campos tocó, no sólo
*si* tocó alguno.

La información ya existe dentro del paquete. `form/form.go` tiene:

```go
func (f *Form) IsDirty() bool {
	for i, sig := range f.valueSignals {
		if sig.Get() != f.baseline[i] {
			return true
		}
	}
	return false
}
```

Compara `valueSignals[i]` contra `baseline[i]` índice a índice. O sea, el dato
por campo está ahí; lo único que falta es exponerlo. Este plan **no** cambia
cómo se calcula el estado sucio, sólo añade un lector.

## 2. Etapa 1 — `DirtyFields`

Fichero: **`form/form.go`**, inmediatamente después de `IsDirty`.

```go
// DirtyFields names the fields whose value differs from the baseline captured
// at the last load/reset — the per-field counterpart of IsDirty, which only
// answers whether ANY of them does.
//
// It exists for bulk edit: a host that applies one form to many records must
// write ONLY the fields the user actually touched, or it silently reverts
// every other column to whatever this form happened to be holding. Returns
// the names in schema order, and nil (not an empty non-nil slice) when
// nothing is dirty, so `len(f.DirtyFields()) == 0` and `!f.IsDirty()` always
// agree.
func (f *Form) DirtyFields() []string {
```

Requisitos exactos:

- Recorre `f.valueSignals` con el mismo índice que `IsDirty`.
- El **nombre** de cada campo sale de la misma fuente que ya usa el
  formulario para su esquema. Localízala leyendo cómo se construyen
  `valueSignals` y `baseline` (probablemente en `form.New` / `LoadValues`) y
  usa exactamente esa lista, en el mismo orden. **No** inventes una segunda
  fuente de nombres ni la deduzcas del DOM.
- Devuelve `nil` cuando no hay ninguno sucio.
- No muta nada: no llama a `MarkPristine`, no reordena, no toca señales.

**Invariante que debe cumplirse siempre** (y que la Etapa 2 comprueba):

```
len(f.DirtyFields()) > 0   ⟺   f.IsDirty() == true
```

Si alguna vez discrepan, uno de los dos está mal. Para que no puedan
divergir, **reescribe `IsDirty` en términos del nuevo recorrido** o extrae un
helper privado que ambos usen. No dejes dos bucles independientes que
comparan lo mismo: es exactamente la clase de duplicación que se desincroniza
en el siguiente refactor.

Sugerencia concreta (elige una y sé consistente):

```go
// helper privado, único sitio que compara contra la baseline
func (f *Form) dirtyIndexes() []int { ... }

func (f *Form) IsDirty() bool      { return len(f.dirtyIndexes()) > 0 }
func (f *Form) DirtyFields() []string { ... usa dirtyIndexes() ... }
```

## 3. Etapa 2 — Tests

Fichero: **`form/dirty_fields_test.go`** (o donde ya vivan los tests de
`IsDirty` — mira primero y ponlos al lado).

| Test | Comprueba |
|---|---|
| `TestDirtyFieldsIsEmptyOnAFreshForm` | Recién construido → `nil`, y `IsDirty()` es `false` |
| `TestDirtyFieldsNamesOnlyTheChangedField` | Cambia un campo → devuelve exactamente ese nombre, sólo uno |
| `TestDirtyFieldsNamesEveryChangedField` | Cambia dos → devuelve los dos, en orden de esquema |
| `TestDirtyFieldsClearsAfterMarkPristine` | Cambia uno, `MarkPristine()` → vuelve a `nil` |
| `TestDirtyFieldsClearsAfterReset` | Cambia uno, `Reset()` → vuelve a `nil` |
| `TestDirtyFieldsAgreesWithIsDirty` | Sobre varios escenarios, `len(DirtyFields()) > 0 == IsDirty()` en todos |
| `TestDirtyFieldsIgnoresAValueSetBackToItsOriginal` | Cambia un campo y devuélvelo a su valor inicial → `nil`. Es una comparación contra la baseline, no un registro de "se tocó" |

El último caso es el que distingue una implementación correcta de una que
lleva un flag "touched": el contrato es **valor distinto de la baseline**, no
"pasó el foco por aquí".

## 4. Etapa 3 — Documentar el uso previsto

Fichero: **`form/README.md`**.

Añade una sección corta, en inglés, bajo el apartado que ya hable del estado
sucio / auto-guardado. Debe decir:

- `IsDirty()` es para decidir *si* persistir (lo que ya hace el auto-guardado
  de `crudview`).
- `DirtyFields()` es para decidir *qué* columnas escribir cuando un mismo
  formulario se aplica a varios registros.
- Un aviso: aplicar un formulario a varios registros escribiendo todas las
  columnas revierte en silencio los campos que el usuario no tocó.

**Anti-footgun de idioma:** este repositorio es librería pública. El código,
los comentarios, los identificadores y el `README.md` van **en inglés**. Sólo
la prosa de este `PLAN.md` está en español.

## 5. Criterios de aceptación

- [ ] `gotest` en verde (vet, race, cover, wasm).
- [ ] `grep -n "func (f \*Form) DirtyFields" form/form.go` → una línea.
- [ ] `IsDirty` sigue existiendo y con la misma firma:
      `grep -n "func (f \*Form) IsDirty() bool" form/form.go` → una línea.
- [ ] No queda ningún bucle duplicado comparando contra `baseline`:
      `grep -cn "f.baseline\[i\]" form/form.go` → **1**.
- [ ] `form/README.md` menciona `DirtyFields`.
- [ ] Ningún consumidor existente rompe: este plan **sólo añade** API.

## 6. Etapas

| # | Etapa | Ficheros | Depende de |
|---|---|---|---|
| 1 | `DirtyFields` + helper compartido | `form/form.go` | — |
| 2 | Tests | `form/dirty_fields_test.go` | 1 |
| 3 | README | `form/README.md` | 1 |
