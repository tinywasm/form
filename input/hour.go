package input

import (
	"github.com/tinywasm/fmt"
)

// hour represents a standard time input.
type hour struct {
	Base
	Permitted Permitted
}

// Hour creates a new time input instance.
func Hour(parentID, name string) Input {
	h := &hour{
		Permitted: Permitted{
			Numbers:    true,
			Characters: []rune{':'},
			Minimum:    5,
			Maximum:    5,
		},
	}
	// htmlName: "time", aliases: "hour"
	h.Base.InitBase(parentID, name, "time", "hour")
	h.Base.SetTitle("formato hora: HH:MM")
	return h
}

// HTMLName returns "time".
func (h *hour) HTMLName() string {
	return h.Base.HTMLName()
}

// ValidateField validates the value against Permitted rules and checking for valid hour format.
func (h *hour) ValidateField(value string) error {
	if len(value) >= 2 && value[0] == '2' && value[1] == '4' {
		// Example: "Hora" + "Inválida"
		return fmt.Err("Hour", "Invalid")
	}
	return h.Permitted.Validate(value)
}

// RenderHTML delegates to Base.RenderInput.
func (h *hour) RenderHTML() string {
	return h.Base.RenderInput()
}

// Clone creates a new time input with the given parentID and name.
func (h *hour) Clone(parentID, name string) Input {
	return Hour(parentID, name)
}
