package input

import "github.com/tinywasm/fmt"

// datalist represents a datalist input field.
// NewDatalist returns a template instance for use in fmt.Field.Widget (no position).
func NewDatalist() fmt.Widget { return Datalist("", "") }

type datalist struct{ Base }

// Datalist creates a new datalist input instance.
func Datalist(parentID, name string) Input {
	dl := &datalist{}
	dl.SkipRules = true
	dl.Base.InitBase(parentID, name, "datalist", "list", "options")
	return dl
}

// Validate validates the selected datalist option against Options.Key.
func (dl *datalist) Validate(value string) error {
	if value == "" || len(dl.Options) == 0 {
		return nil
	}
	for _, opt := range dl.Options {
		if opt.Key == value {
			return nil
		}
	}
	return fmt.Err("Value", value, "NotAllowed", "in", dl.name)
}

// Clone satisfies fmt.Widget — Datalist() returns Input which implements Widget.
func (dl *datalist) Clone(parentID, name string) fmt.Widget { return Datalist(parentID, name) }
