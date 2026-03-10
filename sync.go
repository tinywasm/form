package form

import (
	"github.com/tinywasm/fmt"
)

// SyncValues copies all input values back into the bound struct
// via the Fielder's Pointers() method.
func (f *Form) SyncValues(data fmt.Fielder) error {
	pointers := data.Pointers()
	schema := data.Schema()

	for i, inp := range f.Inputs {
		idx := f.fieldIndices[i]
		if idx < 0 || idx >= len(pointers) {
			continue
		}

		var values []string
		if getter, ok := inp.(interface{ GetValues() []string }); ok {
			values = getter.GetValues()
		}

		ptr := pointers[idx]
		field := schema[idx]

		if len(values) == 0 {
			// Zero the field
			zeroField(ptr, field.Type)
			continue
		}

		writeField(ptr, field.Type, values)
	}
	return nil
}

// zeroField sets a field to its zero value via its pointer.
func zeroField(ptr any, ft fmt.FieldType) {
	switch ft {
	case fmt.FieldText:
		if p, ok := ptr.(*string); ok {
			*p = ""
		}
	case fmt.FieldInt:
		if p, ok := ptr.(*int64); ok {
			*p = 0
		}
	case fmt.FieldFloat:
		if p, ok := ptr.(*float64); ok {
			*p = 0
		}
	case fmt.FieldBool:
		if p, ok := ptr.(*bool); ok {
			*p = false
		}
	}
}

// writeField writes string values into a field via its pointer.
func writeField(ptr any, ft fmt.FieldType, values []string) {
	switch ft {
	case fmt.FieldText:
		if p, ok := ptr.(*string); ok {
			*p = values[0]
		}
	case fmt.FieldInt:
		if p, ok := ptr.(*int64); ok {
			val, _ := fmt.Convert(values[0]).Int64()
			*p = val
		}
	case fmt.FieldFloat:
		if p, ok := ptr.(*float64); ok {
			val, _ := fmt.Convert(values[0]).Float64()
			*p = val
		}
	case fmt.FieldBool:
		if p, ok := ptr.(*bool); ok {
			val, _ := fmt.Convert(values[0]).Bool()
			*p = val
		}
	}
}
