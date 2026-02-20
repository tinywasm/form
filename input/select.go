package input

import "github.com/tinywasm/fmt"

// select_ represents a dropdown selection.
type select_ struct{ Base }

// Select creates a new Select input instance.
func Select(parentID, name string) Input {
	s := &select_{}
	s.Letters = true
	s.Numbers = true
	s.Minimum = 1
	s.InitBase(parentID, name, "select", "role", "tipo")
	return s
}

// RenderHTML renders a select element with options.
func (s *select_) RenderHTML() string {
	out := fmt.GetConv()
	values := s.Base.GetValues()

	out.Write(`<select id="`).Write(s.Base.HandlerName()).Write(`"`)
	out.Write(` name="`).Write(s.Base.FieldName()).Write(`"`)
	if s.Base.Required {
		out.Write(` required`)
	}
	out.Write(`>`)

	for _, opt := range s.Base.GetOptions() {
		selected := false
		for _, v := range values {
			if v == opt.Key {
				selected = true
				break
			}
		}
		out.Write(`<option value="`).Write(opt.Key).Write(`"`)
		if selected {
			out.Write(` selected`)
		}
		out.Write(`>`)
		out.Write(opt.Value)
		out.Write(`</option>`)
	}

	out.Write(`</select>`)
	return out.String()
}

// Clone creates a new Select input.
func (s *select_) Clone(parentID, name string) Input { return Select(parentID, name) }
