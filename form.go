package form

import (
	"reflect"

	"github.com/tinywasm/fmt"
	"github.com/tinywasm/form/input"
)

// Form represents a form instance.
type Form struct {
	ID      string
	Value   any
	Inputs  []input.Input
	class   string // CSS class(es)
	method  string // HTTP method (default POST)
	action  string // Form action URL (default: struct name)
	ssrMode bool   // Per-form SSR mode (default false)
}

// New creates a new Form from a struct pointer.
// parentID: ID of the parent DOM element where the form will be mounted.
// Returns an error if any exported field has no matching registered input.
func New(parentID string, structPtr any) (*Form, error) {
	v := reflect.ValueOf(structPtr)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	t := v.Type()
	structName := fmt.Convert(t.Name()).ToLower().String()

	// Generate form ID from parent and struct name
	formID := parentID + "." + structName

	f := &Form{
		ID:      formID,
		Value:   structPtr,
		Inputs:  make([]input.Input, 0),
		class:   globalClass,
		method:  "POST",
		action:  "/" + structName,
		ssrMode: false,
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Name

		if !fmt.HasUpperPrefix(fieldName) {
			continue
		}

		template := findInputForField(fieldName)
		if template == nil {
			return nil, fmt.Err("field", fieldName, "no matching input registered")
		}

		inp := template.Clone(formID, fieldName)

		// Parse options tag: `options:"key1:text1,key2:text2"`
		if opts, ok := GetTagOptions(string(field.Tag)); ok && len(opts) > 0 {
			inp.SetOptions(opts...)
		}

		// Bind struct field value to input
		fieldValue := v.Field(i)
		switch fieldValue.Kind() {
		case reflect.String:
			inp.SetValues(fieldValue.String())
		case reflect.Slice:
			if fieldValue.Type().Elem().Kind() == reflect.String {
				slice := fieldValue.Interface().([]string)
				inp.SetValues(slice...)
			} else {
				// Convert other slice types to string slice
				slice := make([]string, fieldValue.Len())
				for j := 0; j < fieldValue.Len(); j++ {
					slice[j] = fmt.Convert(fieldValue.Index(j).Interface()).String()
				}
				inp.SetValues(slice...)
			}
		default:
			// Convert other types to string using fmt
			inp.SetValues(fmt.Convert(fieldValue.Interface()).String())
		}

		f.Inputs = append(f.Inputs, inp)
	}

	forms = append(forms, f)
	return f, nil
}

// Input returns the input with the given field name, or nil if not found.
func (f *Form) Input(fieldName string) input.Input {
	for _, inp := range f.Inputs {
		if namer, ok := inp.(interface{ Name() string }); ok {
			if namer.Name() == fieldName {
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
		inp.SetOptions(opts...)
	}
	return f
}

// SetValues sets values for the input matching the given field name.
func (f *Form) SetValues(fieldName string, values ...string) *Form {
	inp := f.Input(fieldName)
	if inp != nil {
		inp.SetValues(values...)
	}
	return f
}
