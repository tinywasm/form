package input

import (
	"github.com/tinywasm/fmt"
)

// gender represents a gender selection (radio buttons).
type gender struct {
	Base
	Permitted Permitted
}

// Gender creates a new Gender input instance with default Male/Female options.
func Gender(parentID, name string) Input {
	g := &gender{
		Permitted: Permitted{
			Letters: true,
			Minimum: 1,
			Maximum: 1,
		},
	}
	g.Base.InitBase(parentID+"."+name, name, "radio", "gender", "sexo")
	// Default options
	g.Base.SetOptions(
		fmt.KeyValue{Key: "m", Value: "Male"},
		fmt.KeyValue{Key: "f", Value: "Female"},
	)
	return g
}

// HtmlName returns "radio".
func (g *gender) HtmlName() string {
	return g.Base.GetHtmlName()
}

// ValidateField validates the value.
func (g *gender) ValidateField(value string) error {
	return g.Permitted.Validate(value)
}

// RenderHTML renders radio buttons for each option.
func (g *gender) RenderHTML() string {
	out := fmt.GetConv()
	values := g.Base.GetValues()

	for _, opt := range g.Base.GetOptions() {
		optID := g.Base.ID() + "." + opt.Key
		selected := false
		for _, v := range values {
			if v == opt.Key {
				selected = true
				break
			}
		}

		out.Write(`<label>`)
		out.Write(`<input type="radio" id="`).Write(optID).Write(`"`)
		out.Write(` name="`).Write(g.Base.Name()).Write(`"`)
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

// Clone creates a new Gender input.
func (g *gender) Clone(parentID, name string) Input {
	return Gender(parentID, name)
}
