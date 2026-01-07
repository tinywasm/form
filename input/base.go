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
	htmlName    string         // The HTML type (e.g., "text", "email")
	aliases     []string       // Field name aliases for matching
	Values      []string       // Multiple values support (for select/checkbox/etc)
	Options     []fmt.KeyValue // Multiple options for select/checkbox/etc
	Placeholder string
	Title       string
	Required    bool // HTML required attribute
	Disabled    bool // HTML disabled attribute
	Readonly    bool // HTML readonly attribute
	Attributes  []Attribute
}

// InitBase initializes the base fields.
func (b *Base) InitBase(id, name, htmlName string, aliases ...string) {
	b.id = id
	b.name = name
	b.htmlName = htmlName
	b.aliases = aliases
}

// SetValues sets the input values.
func (b *Base) SetValues(v ...string) {
	b.Values = v
}

// GetValue returns the first value (for simple inputs).
func (b *Base) GetValue() string {
	if len(b.Values) > 0 {
		return b.Values[0]
	}
	return ""
}

// GetValues returns all input values.
func (b *Base) GetValues() []string {
	return b.Values
}

// SetOptions sets multiple options (for select/checkbox/etc).
func (b *Base) SetOptions(opts ...fmt.KeyValue) {
	b.Options = opts
}

// GetOptions returns all options.
func (b *Base) GetOptions() []fmt.KeyValue {
	return b.Options
}

// GetSelectedValue returns the first value in Values, or empty if none.
func (b *Base) GetSelectedValue() string {
	return b.GetValue()
}

// ID returns the component's unique identifier.
func (b *Base) ID() string {
	return b.id
}

// Name returns the field name (without parent prefix).
func (b *Base) Name() string {
	return b.name
}

// HtmlName returns the HTML input type.
func (b *Base) GetHtmlName() string {
	return b.htmlName
}

// Matches checks if the given field name matches this input's htmlName or aliases.
func (b *Base) Matches(fieldName string) bool {
	name := fmt.Convert(fieldName).ToLower().String()
	if b.htmlName == name {
		return true
	}
	for _, alias := range b.aliases {
		if alias == name {
			return true
		}
	}
	return false
}

// AddAttribute adds a custom attribute to the input.
func (b *Base) AddAttribute(key, value string) {
	b.Attributes = append(b.Attributes, Attribute{Key: key, Value: value})
}

// RenderInput generates the standard HTML tag for the input.
func (b *Base) RenderInput() string {
	out := fmt.GetConv()

	var tag string
	var closeTag string
	isInput := false

	switch b.htmlName {
	case "textarea":
		tag = "textarea"
		closeTag = "</textarea>"
	case "select":
		tag = "select"
		closeTag = "></select>"
	default:
		tag = "input"
		isInput = true
		closeTag = ">"
	}

	out.Write("<").Write(tag)

	if isInput {
		out.Write(` type="`).Write(b.htmlName).Write(`"`)
	}

	out.Write(` id="`).Write(b.id).Write(`"`)
	out.Write(` name="`).Write(b.name).Write(`"`)

	if isInput && b.GetValue() != "" {
		out.Write(` value="`).Write(b.GetValue()).Write(`"`)
	}

	if b.Placeholder != "" {
		out.Write(` placeholder="`).Write(b.Placeholder).Write(`"`)
	}
	if b.Title != "" {
		out.Write(` title="`).Write(b.Title).Write(`"`)
	}

	for _, attr := range b.Attributes {
		if attr.Value != "" {
			out.Write(` `).Write(attr.Key).Write(`="`).Write(attr.Value).Write(`"`)
		}
	}

	// Boolean attributes
	if b.Required {
		out.Write(` required`)
	}
	if b.Disabled {
		out.Write(` disabled`)
	}
	if b.Readonly {
		out.Write(` readonly`)
	}

	if b.htmlName == "textarea" {
		out.Write(">")
		if b.GetValue() != "" {
			out.Write(b.GetValue())
		}
		out.Write(closeTag)
	} else {
		out.Write(closeTag)
	}

	return out.String()
}
