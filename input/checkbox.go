package input

import "github.com/tinywasm/fmt"

// checkbox represents a boolean input field.
type checkbox struct{ Base }

// Checkbox creates a new checkbox input instance.
func Checkbox() Input {
	c := &checkbox{}
	c.InitBase("", "", "checkbox")
	return c
}

// Validate validates only known boolean string values.
func (c *checkbox) Validate(value string) error {
	v := fmt.Convert(value).ToLower().String()
	if v == "" && c.Required {
		return fmt.Err("Field", "Empty", "NotAllowed")
	}
	if v == "" || v == "false" || v == "true" || v == "on" || v == "1" || v == "0" {
		return nil
	}
	return fmt.Err("Format", "Invalid")
}

// Clone satisfies fmt.Widget — Checkbox() returns Input which implements Widget.
func (c *checkbox) Clone(parentID, name string) fmt.Widget {
	c2 := *c
	c2.InitBase(parentID, name, "checkbox")
	return &c2
}
