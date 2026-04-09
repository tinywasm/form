# Plan de Implementación: Compatibilidad de form.Form con dom.Component

## Contexto y Justificación
Durante la implementación de proyectos frontend que utilizan `github.com/tinywasm/form` para generar formularios y `github.com/tinywasm/dom` para el montaje, se ha descubierto un error de compilación: `*form.Form does not implement dom.Component (missing method Children)`.

La raíz del problema es que la interfaz `dom.Component` requiere el método `Children() []dom.Component`. La estructura `*form.Form` posee `GetID`, `SetID` y `RenderHTML`, pero omite `Children`. Dado que los `Input` internos de un formulario ya implementan `dom.Component`, es logico, limpio y arquitectónicamente correcto hacer que el formulario devuelva esos inputs.

**Pros:**
1. Mantener la consistencia arquitectónica del ecosistema `tinywasm`.
2. Permitir el uso directo de formularios en `dom.Render()` sin necesidad de escribir componentes envoltorios (wrappers).
3. Habilitar posibles funcionalidades de recorrido del DOM virtual (re-bound handlers, actualizaciones parciales).

**Contras:**
 - Ninguno, es una extensión esperable y obligatoria del contrato `dom.Component`.

## Pasos de Ejecución

1. **Modificar `form.go`**
   - Asegurar la importación de `github.com/tinywasm/dom`.
   - Implementar el método `Children() []dom.Component` para el tipo `*Form`.
   - La implementación debe crear un slice de `[]dom.Component` con la misma longitud que `f.Inputs`, iterar sobre `f.Inputs` (los cuales implementan `Input` que incluye `dom.Component`) y retornarlos.

```go
// Children returns the form's input fields as dom components.
func (f *Form) Children() []dom.Component {
	children := make([]dom.Component, 0, len(f.Inputs))
	for _, inp := range f.Inputs {
		children = append(children, inp)
	}
	return children
}
```

2. **Verificar Compilación y Tests**
   - Asegurar de que no se rompan pruebas unitarias dentro del propio módulo `form`.
   - Confirmar que `*form.Form` ahora se puede instanciar o asignar satisfactoriamente a interfaces `dom.Component` sin fallar en compilación.
