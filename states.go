package form

import "github.com/tinywasm/widget"

// Pares (clave, valor) de los estados que el campo publica en el DOM. widget.State.Attr()
// devuelve exactamente el par sobre el que selecciona la hoja de estilos, de modo que markup
// y CSS coinciden por construcción y no por convención.
var (
	attrInvalid = widget.Invalid.Attr() // data-invalid = "true"
	attrLocked  = widget.Locked.Attr()  // data-locked  = "true"
)
