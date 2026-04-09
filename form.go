package form

import (
	"github.com/tinywasm/dom"
	"github.com/tinywasm/fmt"
	"github.com/tinywasm/form/input"
)

// Form represents a form instance.
type Form struct {
	id           string
	parentID     string // Parent element ID where the form is mounted
	data         fmt.Fielder
	Inputs       []input.Input
	fieldIndices []int                   // Pre-computed struct field index per Input (-1 if not found)
	class        string                  // CSS class(es)
	method       string                  // HTTP method (default POST)
	action       string                  // Form action URL (default: struct name)
	ssrMode      bool                    // Per-form SSR mode (default false)
	onSubmit     func(fmt.Fielder) error // WASM submit callback
	children     []dom.Component         // Cached dom components (zero-alloc)
}

// Children returns the form's input fields as dom components (O(1), zero-alloc).
func (f *Form) Children() []dom.Component {
	return f.children
}

// GetID returns the html id that group the form
func (f *Form) GetID() string {
	return f.id
}

// SetID sets the html id that group the form
func (f *Form) SetID(id string) {
	f.id = id
}

// ParentID returns the ID of the parent element.
func (f *Form) ParentID() string {
	return f.parentID
}

// OnSubmit sets the callback for form submission in WASM mode.
func (f *Form) OnSubmit(fn func(fmt.Fielder) error) *Form {
	f.onSubmit = fn
	return f
}

// Namer is optionally implemented by Fielder types to provide a custom name.
// If not implemented, the form derives the name from the first Schema field's context.
type Namer interface {
	FormName() string
}

func resolveStructName(data fmt.Fielder) string {
	if n, ok := data.(Namer); ok {
		return n.FormName()
	}
	return "form"
}

// New creates a new Form from a Fielder.
// parentID: ID of the parent DOM element where the form will be mounted.
// Returns an error if any exported field has no matching registered input.
func New(parentID string, data fmt.Fielder) (*Form, error) {
	schema := data.Schema()
	values := fmt.ReadValues(schema, data.Pointers())

	structName := resolveStructName(data)
	formID := parentID + "." + structName

	f := &Form{
		id:       formID,
		parentID: parentID,
		data:     data,
		Inputs:   make([]input.Input, 0, len(schema)),
		class:    globalClass,
		method:   "POST",
		action:   "/" + structName,
		ssrMode:  false,
		children: make([]dom.Component, 0, len(schema)),
	}

	for i, field := range schema {
		// Skip auto-increment PKs (not editable)
		if field.IsPK() && field.IsAutoInc() {
			continue
		}

		fieldName := field.Name

		// Resolve input: use Field.Widget directly.
		if field.Widget == nil {
			continue // skip fields with no UI binding
		}

		inp, ok := field.Widget.Clone(formID, fieldName).(input.Input)
		if !ok {
			continue // skip widgets that do not implement input.Input
		}

		// Apply constraint-based defaults from schema
		if field.NotNull {
			if b, ok := inp.(interface{ SetRequired(bool) }); ok {
				b.SetRequired(true)
			}
		}

		// Bind current value to input
		if setter, ok := inp.(interface{ SetValues(...string) }); ok {
			setter.SetValues(fmt.Convert(values[i]).String())
		}

		f.Inputs = append(f.Inputs, inp)
		f.children = append(f.children, inp)
		f.fieldIndices = append(f.fieldIndices, i)
	}

	forms = append(forms, f)
	return f, nil
}

// Input returns the input with the given field name, or nil if not found.
func (f *Form) Input(fieldName string) input.Input {
	for _, inp := range f.Inputs {
		if getter, ok := inp.(interface{ FieldName() string }); ok {
			if getter.FieldName() == fieldName {
				return inp
			}
		}
	}
	return nil
}

// SetOptions sets options for the input matching the given field name.
func (f *Form) SetOptions(fieldName string, opts ...fmt.KeyValue) *Form {
	inp := f.Input(fieldName)
	if inp != nil {
		if setter, ok := inp.(interface{ SetOptions(...fmt.KeyValue) }); ok {
			setter.SetOptions(opts...)
		}
	}
	return f
}

// SetValues sets values for the input matching the given field name.
func (f *Form) SetValues(fieldName string, values ...string) *Form {
	inp := f.Input(fieldName)
	if inp != nil {
		if setter, ok := inp.(interface{ SetValues(...string) }); ok {
			setter.SetValues(values...)
		}
	}
	return f
}
