package form

import (
	"github.com/tinywasm/fmt"
)

// ValidateData validates a Fielder instance using this form's input rules.
// Satisfies the updated crudp.DataValidator interface (with fmt.Fielder).
func (f *Form) ValidateData(action byte, data fmt.Fielder) error {
	values := fmt.ReadValues(data.Schema(), data.Pointers())
	for i, inp := range f.Inputs {
		idx := f.fieldIndices[i]
		if idx < 0 || idx >= len(values) {
			continue
		}
		if skipper, ok := inp.(interface{ GetSkipValidation() bool }); ok && skipper.GetSkipValidation() {
			continue
		}
		val := fmt.Convert(values[idx]).String()
		if err := inp.Validate(val); err != nil {
			return err
		}
	}
	return nil
}
