package input

import "github.com/tinywasm/fmt"

// checkbox represents a boolean input field.
type checkbox struct {
	Base
}

// Checkbox creates a new checkbox input instance.
func Checkbox(parentID, name string) Input {
	c := &checkbox{}
	// htmlName: "checkbox", aliases: "check", "boolean", "bool"
	c.Base.InitBase(parentID, name, "checkbox", "check", "boolean", "bool")
	return c
}

// HTMLName returns "checkbox".
func (c *checkbox) HTMLName() string {
	return c.Base.HTMLName()
}

// ValidateField validates only if it is "true", "false", or "on".
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

// RenderHTML delegates to Base.RenderInput.
func (c *checkbox) RenderHTML() string {
	return c.Base.RenderInput()
}

// Clone creates a new checkbox input with the given parentID and name.
func (c *checkbox) Clone(parentID, name string) Input {
	return Checkbox(parentID, name)
}
