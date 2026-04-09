# Plan de Optimización: Caching de Children en Form

## Contexto y Justificación
Durante la implementación de la interfaz `dom.Component` para `*form.Form`, se introdujo el método `Children()` el cual crea y retorna un nuevo *slice* `[]dom.Component` en cada llamada:
```go
func (f *Form) Children() []dom.Component {
	children := make([]dom.Component, 0, len(f.Inputs))
	// ... append y return
}
```
Esta aproximación genera asignación de memoria dinámica en el *heap* del motor (en WebAssembly/TinyGo) con cada ciclo de renderizado o recorrido del DOM virtual. Puesto que los *inputs* del formulario se determinan en el momento de su inicialización (en la función `New`), podemos pre-calcular este slice y guardarlo de forma estática en la estructura para evitar re-asignaciones, logrando cero allocs por llamada.

## Ejecución

### 1. Modificar la estructura `Form`
En el archivo `form.go`, añadir un campo privado `children` a la declaración de `Form` para almacenar las referencias previas:
```go
type Form struct {
	// ... campos actuales ...
	Inputs       []input.Input
	// ...
	onSubmit     func(fmt.Fielder) error
	children     []dom.Component // Cache para evitar allocations dinámicos en el heap
}
```

### 2. Modificar la inicialización en `New()`
Dentro de la función `New()` en `form.go`, inicializar `f.children` con capacidad pre-alojada y llenarlo mientras se agregan elementos a `f.Inputs`:
```go
	f := &Form{
		// ...
		Inputs:   make([]input.Input, 0, len(schema)),
		children: make([]dom.Component, 0, len(schema)),
		// ...
	}

	for i, field := range schema {
		// ... validaciones ...
		
		f.Inputs = append(f.Inputs, inp)
		f.children = append(f.children, inp) // Almacenar el input como componente dom
		f.fieldIndices = append(f.fieldIndices, i)
	}
```

### 3. Modificar el método `Children()`
Actualizar el método en `form.go` para que devuelva directamente el campo cacheado, operando de forma Inmediata (O(1)):
```go
// Children returns the form's input fields as dom components (O(1), zero-alloc).
func (f *Form) Children() []dom.Component {
	return f.children
}
```

## Resultados Esperados
Este cambio reducirá la presión en el GC y asegurará un alto rendimiento en la traversing de los hijos del formulario de parte de librerías como `dom.Render` con una asignación de memoria en el Heap equivalente a `0 allocs/op` por llamada.
