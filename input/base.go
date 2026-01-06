package input

import (
	"github.com/tinywasm/fmt"
)

// Attribute key-value pair for HTML attributes.
type Attribute struct {
	Key   string
	Value string
}

// Base contains common logic and fields (State) for all inputs.
// It is intended to be embedded in concrete input structs.
type Base struct {
	id          string
	name        string
	Value       string
	Placeholder string
	Title       string
	// Attributes slice instead of map for optimization
	Attributes []Attribute
}

// InitBase initializes the base fields.
func (b *Base) InitBase(id, name string) {
	b.id = id
	b.name = name
	// No map make needed
}

// ID returns the component's unique identifier.
func (b *Base) ID() string {
	return b.id
}

// AddAttribute adds a custom attribute to the input.
func (b *Base) AddAttribute(key, value string) {
	b.Attributes = append(b.Attributes, Attribute{Key: key, Value: value})
}

// RenderInput generates the standard HTML tag for the input.
// htmlType: "text", "email", "textarea", "select", etc.
func (b *Base) RenderInput(htmlType string) string {
	out := fmt.GetConv()

	var tag string
	var closeTag string
	isInput := false

	switch htmlType {
	case "textarea":
		tag = "textarea"
		closeTag = "</textarea>"
	case "select":
		tag = "select"
		closeTag = "></select>" // Standard generic close, assumes options added separately or empty
	default:
		tag = "input"
		isInput = true
		closeTag = ">"
	}

	out.Write("<").Write(tag)

	if isInput {
		out.Write(` type="`).Write(htmlType).Write(`"`)
	}

	out.Write(` id="`).Write(b.id).Write(`"`)
	out.Write(` name="`).Write(b.name).Write(`"`)

	// Standard attributes (Value handled differently for textarea)
	if isInput && b.Value != "" {
		out.Write(` value="`).Write(b.Value).Write(`"`)
	}

	if b.Placeholder != "" {
		out.Write(` placeholder="`).Write(b.Placeholder).Write(`"`)
	}
	if b.Title != "" {
		out.Write(` title="`).Write(b.Title).Write(`"`)
	}

	// Render generic attributes from slice
	for _, attr := range b.Attributes {
		if attr.Value != "" {
			out.Write(` `).Write(attr.Key).Write(`="`).Write(attr.Value).Write(`"`)
		}
	}

	// Close tag content
	if htmlType == "textarea" {
		out.Write(">")
		if b.Value != "" {
			out.Write(b.Value)
		}
		out.Write(closeTag)
	} else {
		out.Write(closeTag)
	}

	return out.String()
}
