package input

import "github.com/tinywasm/fmt"

// hour represents a standard time input.
type hour struct{ Base }

// Hour creates a new time input instance.
func Hour(parentID, name string) Input {
	h := &hour{}
	h.Numbers = true
	h.Characters = []rune{':'}
	h.Minimum = 5
	h.Maximum = 5
	h.InitBase(parentID, name, "time", "hour")
	h.SetTitle("formato hora: HH:MM")
	return h
}

// ValidateField validates HH:MM format rejecting 24:xx.
func (h *hour) ValidateField(value string) error {
	if len(value) >= 2 && value[0] == '2' && value[1] == '4' {
		return fmt.Err("Hour", "Invalid")
	}
	return h.Permitted.Validate(value)
}

// Clone creates a new time input with the given parentID and name.
func (h *hour) Clone(parentID, name string) Input { return Hour(parentID, name) }
