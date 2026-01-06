package form

import (
	"reflect"

	"github.com/tinywasm/fmt"
	"github.com/tinywasm/form/input"
)

// Form represents a form instance.
type Form struct {
	ID        string
	Value     any // The struct instance
	Inputs    []input.Input
	class     string // Default CSS class for inputs
	method    string // HTTP method (POST/GET)
	targetURI string // Action URL
}

// Global storage for forms to keep them alive (private)
var forms = make([]*Form, 0)

// SetGlobalClass sets a default CSS class for all subsequently created forms/inputs.
// This is a placeholder for global configuration.
func SetGlobalClass(class string) {
	// Implementation TODO: Store this in a package-level config variable
}

// New creates a new Form from a struct pointer.
// It uses reflection to discover fields and automatically create corresponding Inputs.
func New(id string, structPtr any) *Form {
	f := &Form{
		ID:     id,
		Value:  structPtr,
		Inputs: make([]input.Input, 0),
	}

	v := reflect.ValueOf(structPtr)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Name

		// Skip unexported fields
		if !fmt.HasUpperPrefix(fieldName) {
			continue
		}

		// Field Discovery Logic
		// 1. Check if type implements Input (not yet supported directly on struct fields)
		// 2. Map standard Go types to Form Inputs

		// TODO: Add more sophisticated type mapping (Email, RUT, etc.) based on name/type
		// For now, default everything to Text input
		inp := input.Text(id, fieldName)
		f.Inputs = append(f.Inputs, inp)
	}

	// Register in global scope
	forms = append(forms, f)

	return f
}

// RenderHTML renders the entire form.
func (f *Form) RenderHTML() string {
	out := fmt.GetConv()
	out.Write(`<form id="`).Write(f.ID).Write(`">`)
	for _, inp := range f.Inputs {
		out.Write(inp.RenderHTML())
	}
	out.Write("</form>")
	return out.String()
}
