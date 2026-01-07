package input

import "github.com/tinywasm/fmt"

// radio represents a standard radio button input.
type radio struct {
	Base
	Permitted Permitted
}

// Radio creates a new Radio input instance.
func Radio(parentID, name string) Input {
	r := &radio{
		Permitted: Permitted{
			Letters: true,
			Numbers: true,
			Minimum: 1,
		},
	}
	r.Base.InitBase(parentID, name, "radio")
	return r
}

// HTMLName returns "radio".
func (r *radio) HTMLName() string {
	return "radio"
}

// ValidateField validates the value against Permitted rules.
func (r *radio) ValidateField(value string) error {
	return r.Permitted.Validate(value)
}

// RenderHTML renders radio buttons.
func (r *radio) RenderHTML() string {
	out := fmt.GetConv()
	values := r.Base.GetValues()

	for _, opt := range r.Base.GetOptions() {
		optID := r.Base.ID() + "." + opt.Key
		selected := false
		for _, v := range values {
			if v == opt.Key {
				selected = true
				break
			}
		}

		out.Write(`<label>`)
		out.Write(`<input type="radio" id="`).Write(optID).Write(`"`)
		out.Write(` name="`).Write(r.Base.FieldName()).Write(`"`)
		out.Write(` value="`).Write(opt.Key).Write(`"`)
		if selected {
			out.Write(` checked`)
		}
		out.Write(`>`)
		out.Write(opt.Value)
		out.Write(`</label>`)
	}
	return out.String()
}

// Clone creates a new Radio input.
func (r *radio) Clone(parentID, name string) Input {
	return Radio(parentID, name)
}
