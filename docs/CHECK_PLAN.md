# PLAN — Desacoplar metadata de widget del rendering (fuga de tamaño al edge). 

> Este plan se despacha vía el workflow CodeJob. Ver skill: `agents-workflow`.
> **Estado:** LISTO PARA DESPACHO.
> **Repo objetivo:** `github.com/tinywasm/form` (incluye el subpaquete `input`).
> **Impacto:** ~40–45 KB menos en binarios wasm que solo usan el `Schema()` del ORM (edge
> workers) — validado aguas abajo en `goflare-demo` (no es responsabilidad del agente).

## Prerequisito (PRIMERO — entorno del agente)

El agente NO puede probar en el navegador. Toda verificación es con `gotest`, que no viene
preinstalado en el entorno aislado del agente. Instalarlo antes de cualquier otra cosa:

```bash
go install github.com/tinywasm/devflow/cmd/gotest@latest
```

`gotest` corre los tests de backend (stdlib) **y** los de WASM con build tags, además de
`-vet`/`-race`/`-cover`. Usar `gotest` (sin argumentos para toda la suite, o
`gotest -run TestX`); **NO** usar `go test` directo.

## Problema

`github.com/tinywasm/form/input` arrastra `tinywasm/html` y `tinywasm/dom` (código de
**rendering**) a **cualquier** binario que solo necesite los **metadatos** del widget. Un
Cloudflare Worker (edge) que solo crea la tabla y valida —sin renderizar HTML— termina
incluyendo todo el stack de formularios + `regexp` (~44 KB de más).

### Causa raíz

1. `input/interface.go` — la interfaz `Input` **embebe `dom.Component`** además de `fmt.Widget`.
2. `input/base.go` importa `tinywasm/html` y `tinywasm/dom`: los métodos de render
   (`String()→RenderInput()`, `Children() []dom.Component`, `renderSelect()`, `renderRadio()`)
   viven en `Base`, que embeben TODOS los widgets (`type text struct{ Base }`).
3. El `model_orm.go` generado por `ormc` referencia `input.Text()/Email()/...` en `Schema()`.
   Como `Schema()` se usa **agnósticamente** (el edge lo llama para `CreateTable` + validación),
   importar el paquete `input` mete `dom`+`html`+`regexp` aunque solo se usen `Type()/Validate()`.

## Principio rector (NO romper)

> Código **agnóstico** (compila wasm **y** backend: schema, validación) **no debe** importar
> rendering (`dom`/`html`) ni `regexp`. Cada input valida con su propia lógica Go. El
> rendering es exclusivo de frontend.

## Decisión arquitectónica (resuelta)

**Separación de PAQUETE — NO build tags.** El edge (worker) y el frontend (browser) se
compilan **ambos** con `GOOS=js GOARCH=wasm`: ningún build tag (`wasm`/`!wasm`) los separa, y
el render del form ocurre EN wasm (browser), así que tampoco se puede tagear `!wasm`. La
única frontera que crea un grafo de imports sin `dom`/`html` en el edge es el **límite de
paquete**. `input` y `form` están en el **mismo módulo**, así que es un refactor de un repo.

**Mecanismo de render: tipos/funciones de render en el paquete `form`** (que sí importa
`dom`/`html`), que reciben la metadata del widget (vía getters de `input`) y producen el
HTML/DOM, despachando por `Type()` con el registry existente (`registry.go`,
`findInputByType`). Se aprovecha la costura que YA existe: `form.go` hace
`field.Widget.Clone(...).(input.Input)` y `render.go` hace `inp.String()`.

Contrato resultante:
- `input` → widgets que implementan **solo `fmt.Widget`** (`Type`, `Validate`, `Clone`) +
  getters de datos. **Sin importar `dom`/`html`.** `input.Text()` &co. siguen existiendo y
  devuelven `fmt.Widget` (la firma NO cambia → `ormc` no se toca).
- `form` → dueño del render (importa `dom`/`html`). Provee la capa que renderiza un
  `fmt.Widget` según su `Type()`.
- `input.Input` deja de embeber `dom.Component`.

## Pasos de ejecución

### Stage 1 — `input` agnóstico (sin `dom`/`html`)
1. Quitar de `input/base.go` los métodos de render (`String`, `Children`, `RenderInput`,
   `renderSelect`, `renderRadio` y cualquier helper que use `html`/`dom`) y sus imports
   `tinywasm/html` y `tinywasm/dom`.
2. Conservar en `input` toda la metadata/validación + getters que el render necesitará
   (htmlName/`Type()`, placeholder, title, options, required, value(s), atributos, errorID,
   reglas `fmt.Permitted`). Exponer getters que falten para que el render externo arme el HTML.
3. Quitar `dom.Component` de la interfaz `Input` en `input/interface.go`. Si algún consumidor
   necesita un contrato de render, ese contrato vive en `form`, no en `input`.
4. `input` solo puede importar `tinywasm/fmt` (y stdlib mínimos permitidos en wasm).

### Stage 2 — render en `form`
5. Mover la lógica de render migrada del Stage 1 a `form` (p.ej. `form/render_input.go`),
   como funciones que reciben la metadata de `input` (vía getters/una interfaz de solo-lectura
   expuesta por `input`) y devuelven el HTML/`dom` equivalente al actual.
6. Reconectar el paquete `form` para que use esa capa en lugar de los métodos que estaban en
   `input`: `form.go` (el `.(input.Input)` / construcción de `children []dom.Component`),
   `render.go` (`inp.String()`), `mount.go`, `registry.go`. `form` importa `dom`/`html` (es
   frontend: permitido).
7. Preservar el comportamiento observable: mismo HTML generado, mismos IDs (`b.id`,
   `ErrorID()`), mismo registry por `Type()`.

### Stage 3 — verificación
8. `gotest` verde (backend + WASM). Ajustar tests existentes que asumían render en `input`.
9. Confirmar que `input` ya no importa `dom`/`html` (comando en §Verificación).

## Verificación (repo-local, ejecutable por el agente)

```bash
# 1. El subpaquete input NO importa rendering (criterio de aceptación principal):
GOOS=js GOARCH=wasm go list -deps ./input | grep -E 'tinywasm/(html|dom)' && echo "FALLA: input aún arrastra rendering" || echo "OK: input agnóstico"

# 2. input solo depende de tinywasm/fmt (entre las libs tinywasm):
GOOS=js GOARCH=wasm go list -f '{{.Imports}}' ./input | tr ' ' '\n' | grep tinywasm
#    → solo github.com/tinywasm/fmt

# 3. Tests verdes (backend + WASM, render del frontend intacto):
gotest
```

> Validación de tamaño aguas abajo (NO la hace el agente): en `goflare-demo`,
> `GOOS=js GOARCH=wasm go list -deps ./edge | grep -E 'tinywasm/(html|dom)|^regexp$'` debe
> quedar vacío y `edge.wasm` bajar ~40–45 KB.

## Checklist de calidad (obligatorio)

- **Sin strings hardcodeados repetidos:** nombres de tipo HTML, clases CSS (`tw-field-error`,
  etc.), `aria-*` y prefijos de ID que se repitan → constantes nombradas en el paquete. Nada
  de literales duplicados en la lógica.
- **Sin duplicación lógica:** la lógica de render se MUEVE (no se copia). No dejar copias en
  `input` y `form`.
- **Reglas tinywasm:**
  - Nada de stdlib pesado en código wasm: usar `tinywasm/fmt` (no `errors`/`strconv`/`strings`).
  - Embebido por valor (no punteros) para tipos `dom`.
  - `input` debe quedar libre de `dom`/`html`/`regexp`/`reflect`/`sync`.

## Tabla de stages

| Stage | Objetivo | Entregable | Criterio de salida |
|---|---|---|---|
| 1 | `input` agnóstico | `input/base.go` y `input/interface.go` sin `dom`/`html`; getters de metadata | `go list -deps ./input` sin `tinywasm/html|dom` |
| 2 | Render en `form` | render migrado a `form` + rewire (`form.go`, `render.go`, `mount.go`, `registry.go`) | compila; mismo HTML/IDs/registry |
| 3 | Verificación | tests ajustados | `gotest` verde |

## Nota relacionada (fuera de este repo)

En `goflare-demo`, `modules/contact` también importa `tinywasm/fetch` (cliente HTTP de
browser) en el edge. Se corrige aparte, por **paquete** (no build tags: edge y cliente son
ambos wasm). No es parte de este plan.
