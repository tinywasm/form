package input

import "github.com/tinywasm/fmt"

// checkbox represents a boolean input field.
type checkbox struct{ Base }

// Checkbox creates a new checkbox input instance.
func Checkbox(parentID, name string) Input {
	c := &checkbox{}
	c.InitBase(parentID, name, "checkbox", "check", "boolean", "bool")
	return c
}

// ValidateField validates only known boolean string values.
func (c *checkbox) ValidateField(value string) error {
	v := fmt.Convert(value).ToLower().String()
	if v == "" && c.Required {
		return fmt.Err("Field", "Empty", "NotAllowed")
	}
	if v == "" || v == "false" || v == "true" || v == "on" || v == "1" || v == "0" {
		return nil
	}
	return fmt.Err("Format", "Invalid")
}

// Clone creates a new checkbox input with the given parentID and name.
func (c *checkbox) Clone(parentID, name string) Input { return Checkbox(parentID, name) }
