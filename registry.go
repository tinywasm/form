package form

import (
	"github.com/tinywasm/fmt"
	"github.com/tinywasm/form/input"
)

// Global storage for forms
var forms = make([]*Form, 0)

// Global registry for input types
var registeredInputs = make([]input.Input, 0)

func init() {
	RegisterInput(
		input.Text("", ""),
		input.Email("", ""),
		input.Password("", ""),
		input.Gender("", ""),
		input.Select("", ""),
		input.Radio("", ""),
		input.Address("", ""),
	)
}

// RegisterInput registers input types for field mapping.
func RegisterInput(inputs ...input.Input) {
	registeredInputs = append(registeredInputs, inputs...)
}

// findInputForField searches for a registered input that matches the field name.
func findInputForField(fieldName, structName string) input.Input {
	name := fmt.Convert(fieldName).ToLower().String()
	fullName := fmt.Convert(structName + "." + fieldName).ToLower().String()

	for _, inp := range registeredInputs {
		if matcher, ok := inp.(interface{ Matches(string) bool }); ok {
			if matcher.Matches(name) || matcher.Matches(fullName) {
				return inp
			}
		}
	}
	return nil
}

// Global class configuration
var globalClass string

func SetGlobalClass(classes ...string) {
	for _, c := range classes {
		if globalClass != "" {
			globalClass += " "
		}
		globalClass += c
	}
}
