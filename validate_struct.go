package form

import (
	"reflect"

	"github.com/tinywasm/fmt"
)

// ValidateData satisfies crudp.DataValidator via Go duck typing.
// It validates the struct in data[0] using this form's configured input rules.
//
// Design:
//   - Reuses the existing form instance — no allocation per call.
//   - Uses pre-computed fieldIndices (set in New) for O(1) field access.
//   - Calls inp.ValidateField(val) which is a pure function — no state mutation.
//   - Thread-safe: reads only immutable form config, never writes.
//
// The action byte follows crudp conventions: 'c'=create, 'r'=read, 'u'=update, 'd'=delete.
// All actions run the same validation; handlers may add action-specific logic on top.
func (f *Form) ValidateData(action byte, data ...any) error {
	if len(data) == 0 {
		return nil
	}
	v := reflect.ValueOf(data[0])
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	for i, inp := range f.Inputs {
		idx := f.fieldIndices[i]
		if idx < 0 {
			continue
		}
		if skipper, ok := inp.(interface{ GetSkipValidation() bool }); ok && skipper.GetSkipValidation() {
			continue
		}
		val := fmt.Convert(v.Field(idx).Interface()).String()
		if err := inp.ValidateField(val); err != nil {
			return err
		}
	}
	return nil
}
