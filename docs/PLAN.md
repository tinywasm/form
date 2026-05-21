# PLAN — `ssr.go` → split por extensión (`css.go`)

## Objetivo

Renombrar `form/ssr.go` a `form/css.go` para alinearse con la nueva convención
del motor de `assetmin`: los assets SSR se descubren por archivos con nombre de
extensión (`css.go`, `js.go`, `html.go`, `svg.go`), todos `//go:build !wasm`.
El nombre reservado `ssr.go` se elimina del ecosistema.

## Justificación

`ssr.go` es un nombre mágico que no comunica su contenido. `css.go` es
autoexplicativo y alinea con SRP (`core-principles`). Breaking change
coordinado a nivel monorepo — ver el stage homónimo en `assetmin/docs/PLAN.md`.

## Estado actual

`form/ssr.go` contiene una única función package-level:

- `RenderCSS() *css.Stylesheet` — CSS de los mensajes de error de campo
  (`.tw-field-error`).

Solo CSS, así que el archivo se renombra directo a `css.go`.

## Cambios

- Renombrar `form/ssr.go` → `form/css.go`. Contenido **literal**: mismo build
  tag, mismo package, misma función.
- Actualizar el comentario de doc de `RenderCSS()`: hoy dice *"Call from the
  project's ssr.go aggregate"*. Tras la migración no existe `ssr.go`; el texto
  debe referir al `css.go` agregador del proyecto.

## Precondición técnica

`assetmin` debe estar publicado con la whitelist `ssrSourceFiles`
(`css.go/js.go/svg.go/html.go`) y sin reconocer ya `ssr.go`. Aplicar este
renombrado en el **mismo PR coordinado** que el cambio de motor.

```bash
go list -m github.com/tinywasm/assetmin
```

## Tests y validación

- `go test ./...` verde en `tinywasm/form`.
- Verificar que `assetmin` sigue recogiendo `.tw-field-error` en el slot
  `middle` tras el renombrado.

## Stages

| # | Tarea | Done |
|---|---|---|
| 1 | Confirmar precondición: `assetmin` con whitelist `ssrSourceFiles` publicado | [ ] |
| 2 | Renombrar `form/ssr.go` → `form/css.go` (contenido literal) | [ ] |
| 3 | Actualizar el comentario de doc que menciona "project's ssr.go aggregate" | [ ] |
| 4 | `go test ./...` verde | [ ] |
| 5 | Verificar extracción de CSS de form vía `assetmin` sin regresiones | [ ] |
