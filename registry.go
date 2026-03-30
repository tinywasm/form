package form

import (
	"github.com/tinywasm/form/input"
)

// Global storage for forms
var forms = make([]*Form, 0)

// Global registry for input types
var registeredInputs = make([]input.Input, 0)

// RegisterInput registers input types for field mapping.
func RegisterInput(inputs ...input.Input) {
	registeredInputs = append(registeredInputs, inputs...)
}

// findInputByType finds a registered input template by its HTMLName.
// Returns nil if no match found.
func findInputByType(htmlType string) input.Input {
	for _, tmpl := range registeredInputs {
		if tmpl.Type() == htmlType {
			return tmpl
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
