package input

import "github.com/tinywasm/fmt"

type hour struct{ Base }

// Hour creates a new time input instance.
func Hour(parentID, name string) Input {
	h := &hour{}
	h.Numbers = true
	h.Characters = []rune{':'}
	h.Minimum = 0
	h.Maximum = 5
	h.InitBase(parentID, name, "time", "hour")
	h.SetTitle("formato hora: HH:MM")
	return h
}

// Validate validates HH:MM format rejecting 24:xx.
func (h *hour) Validate(value string) error {
	if value == "" {
		return nil
	}
	if len(value) != 5 {
		return fmt.Err("Hour", "Invalid")
	}
	if value[0] == '2' && value[1] == '4' {
		return fmt.Err("Hour", "Invalid")
	}
	return h.Permitted.Validate(value)
}

// Clone satisfies fmt.Widget — Hour() returns Input which implements Widget.
func (h *hour) Clone(parentID, name string) fmt.Widget { return Hour(parentID, name) }
