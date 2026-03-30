package input

import "github.com/tinywasm/fmt"

// checkbox represents a boolean input field.
type checkbox struct{ Base }

// NewCheckbox returns a template instance for use in fmt.Field.Widget (no position).
func NewCheckbox() fmt.Widget {
	return Checkbox("", "")
}

// Checkbox creates a new checkbox input instance.
func Checkbox(parentID, name string) Input {
	c := &checkbox{}
	c.SkipRules = true
	c.InitBase(parentID, name, "checkbox", "check", "boolean", "bool")
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
func (c *checkbox) Clone(parentID, name string) fmt.Widget { return Checkbox(parentID, name) }
