package form

import (
	"reflect"

	"github.com/tinywasm/fmt"
	"github.com/tinywasm/form/input"
)

// Form represents a form instance.
type Form struct {
	id           string
	parentID     string // Parent element ID where the form is mounted
	Value        any
	Inputs       []input.Input
	fieldIndices []int           // Pre-computed struct field index per Input (-1 if not found)
	class        string          // CSS class(es)
	method       string          // HTTP method (default POST)
	action       string          // Form action URL (default: struct name)
	ssrMode      bool            // Per-form SSR mode (default false)
	onSubmit     func(any) error // WASM submit callback
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
func (f *Form) OnSubmit(fn func(any) error) *Form {
	f.onSubmit = fn
	return f
}

// SyncValues synchronizes all input values back to the source struct.
func (f *Form) SyncValues() error {
	v := reflect.ValueOf(f.Value)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	for _, inp := range f.Inputs {
		fieldName := inp.FieldName()
		field := v.FieldByName(fieldName)
		if !field.IsValid() || !field.CanSet() {
			continue
		}

		var values []string
		if getter, ok := inp.(interface{ GetValues() []string }); ok {
			values = getter.GetValues()
		}
		if len(values) == 0 {
			field.Set(reflect.Zero(field.Type()))
			continue
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(values[0])
		case reflect.Slice:
			if field.Type().Elem().Kind() == reflect.String {
				field.Set(reflect.ValueOf(values))
			}
		}
	}
	return nil
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
		id:       formID,
		parentID: parentID,
		Value:    structPtr,
		Inputs:   make([]input.Input, 0),
		class:    globalClass,
		method:   "POST",
		action:   "/" + structName,
		ssrMode:  false,
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldName := field.Name

		if !fmt.HasUpperPrefix(fieldName) {
			continue
		}

		template := findInputForField(fieldName, structName)
		if template == nil {
			return nil, fmt.Err("field", fieldName, "no matching input registered")
		}
		inp := template.Clone(formID, fieldName)

		// Get struct tags
		tag := string(field.Tag)
		conv := fmt.Convert(tag)

		// 1. Check for validate:"false"
		if valTag, _ := conv.TagValue("validate"); valTag == "false" {
			if b, ok := inp.(interface{ SetSkipValidation(bool) }); ok {
				b.SetSkipValidation(true)
			}
		}

		// 2. Custom Placeholder
		if ph, _ := conv.TagValue("placeholder"); ph != "" {
			if b, ok := inp.(interface{ SetPlaceholder(string) }); ok {
				b.SetPlaceholder(ph)
			}
		}

		// 3. Custom Title
		if title, _ := conv.TagValue("title"); title != "" {
			if b, ok := inp.(interface{ SetTitle(string) }); ok {
				b.SetTitle(title)
			}
		}

		// 4. Parse options tag: `options:"key1:text1,key2:text2"`
		if opts := conv.TagPairs("options"); len(opts) > 0 {
			if setter, ok := inp.(interface{ SetOptions(...fmt.KeyValue) }); ok {
				setter.SetOptions(opts...)
			}
		}

		// Bind struct field value to input
		fieldValue := v.Field(i)
		if setter, ok := inp.(interface{ SetValues(...string) }); ok {
			switch fieldValue.Kind() {
			case reflect.String:
				setter.SetValues(fieldValue.String())
			case reflect.Slice:
				if fieldValue.Type().Elem().Kind() == reflect.String {
					slice := fieldValue.Interface().([]string)
					setter.SetValues(slice...)
				} else {
					slice := make([]string, fieldValue.Len())
					for j := 0; j < fieldValue.Len(); j++ {
						slice[j] = fmt.Convert(fieldValue.Index(j).Interface()).String()
					}
					setter.SetValues(slice...)
				}
			default:
				setter.SetValues(fmt.Convert(fieldValue.Interface()).String())
			}
		}

		f.Inputs = append(f.Inputs, inp)
	}

	// Pre-compute struct field indices for each input.
	// Stored as []int so ValidateData can use v.Field(idx) — O(1) access, no per-call FieldByName search.
	f.fieldIndices = make([]int, len(f.Inputs))
	for i, inp := range f.Inputs {
		if sf, ok := t.FieldByName(inp.FieldName()); ok {
			f.fieldIndices[i] = sf.Index[0]
		} else {
			f.fieldIndices[i] = -1
		}
	}

	forms = append(forms, f)
	return f, nil
}

// Input returns the input with the given field name, or nil if not found.
func (f *Form) Input(fieldName string) input.Input {
	for _, inp := range f.Inputs {
		if inp.FieldName() == fieldName {
			return inp
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
