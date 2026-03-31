package input

import "github.com/tinywasm/fmt"

type datalist struct{ Base }

// Datalist creates a new datalist input instance.
func Datalist() Input {
	dl := &datalist{}
	dl.Base.InitBase("", "", "datalist")
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
func (dl *datalist) Clone(parentID, name string) fmt.Widget {
	c := *dl
	c.InitBase(parentID, name, "datalist")
	return &c
}
