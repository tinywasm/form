package input

import "github.com/tinywasm/fmt"

// datalist represents a datalist input field.
type datalist struct {
	Base
}

// Datalist creates a new datalist input instance.
func Datalist(parentID, name string) Input {
	dl := &datalist{}
	// htmlName: "text" (with a <datalist> under it), aliases: "list", "options"
	dl.Base.InitBase(parentID, name, "text", "list", "options")
	return dl
}

// HTMLName returns "text".
func (dl *datalist) HTMLName() string {
	return dl.Base.HTMLName() // Real input is text, accompanied by datalist
}

// ValidateField validates the selected datalist option.
func (dl *datalist) ValidateField(value string) error {
	v := fmt.Convert(value).ToLower().String()
	if v == "" && dl.Required {
		return fmt.Err("Field", "Empty", "NotAllowed")
	}

	if len(dl.Options) == 0 {
		return nil
	}

	for _, opt := range dl.Options {
		// Datalist match against Key or Value depending on your business logic. Here we will match Key.
		// Using fmt.Convert for lowercase string comparison
		if fmt.Convert(opt.Key).ToLower().String() == v {
			return nil
		}
	}

	if v != "" {
		return fmt.Err("Value", value, "NotAllowed", "in", dl.name)
	}

	return nil
}

// RenderHTML overrides Base.RenderInput to include a <datalist> element attached to the <input>.
func (dl *datalist) RenderHTML() string {
	out := fmt.GetConv()

	// Ensure input links to the datalist via the "list" attribute
	dl.AddAttribute("list", dl.id+"-list")

	// Core <input>
	out.Write(dl.Base.RenderInput())

	// Datalist <datalist> block
	out.Write(`<datalist id="`).Write(dl.id).Write(`-list">`)
	for _, opt := range dl.Options {
		out.Write(`<option value="`).Write(opt.Key).Write(`">`).Write(opt.Value).Write(`</option>`)
	}
	out.Write(`</datalist>`)

	return out.String()
}

// Clone creates a new datalist input with the given parentID and name.
func (dl *datalist) Clone(parentID, name string) Input {
	return Datalist(parentID, name)
}
