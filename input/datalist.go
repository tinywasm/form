package input

import "github.com/tinywasm/fmt"

// datalist represents a datalist input field.
type datalist struct{ Base }

// Datalist creates a new datalist input instance.
func Datalist(parentID, name string) Input {
	dl := &datalist{}
	dl.Base.InitBase(parentID, name, "datalist", "list", "options")
	return dl
}

// ValidateField validates the selected datalist option against Options.Key.
func (dl *datalist) ValidateField(value string) error {
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

// Clone creates a new datalist input with the given parentID and name.
func (dl *datalist) Clone(parentID, name string) Input { return Datalist(parentID, name) }
