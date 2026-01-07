package input

import (
	"github.com/tinywasm/fmt"
)

// Base contains common logic and fields (State) for all inputs.
// It is intended to be embedded in concrete input structs.
type Base struct {
	id             string
	name           string
	htmlName       string         // The HTML type (e.g., "text", "email")
	aliases        []string       // Field name aliases for matching
	Values         []string       // Multiple values support (for select/checkbox/etc)
	Options        []fmt.KeyValue // Multiple options for select/checkbox/etc
	Placeholder    string
	Title          string
	Required       bool // HTML required attribute
	Disabled       bool // HTML disabled attribute
	Readonly       bool // HTML readonly attribute
	SkipValidation bool // Whether to skip validation for this input
	Attributes     []fmt.KeyValue
}

// InitBase initializes the base fields and constructs the unique ID.
func (b *Base) InitBase(parentID, name, htmlName string, aliases ...string) {
	if parentID != "" {
		b.id = parentID + "." + name
	} else {
		b.id = name
	}
	b.name = name
	b.htmlName = htmlName
	b.aliases = aliases

	// Auto-defaults
	b.Placeholder = "Enter " + name
	b.Title = name + " field"
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

// SetPlaceholder sets the input placeholder.
func (b *Base) SetPlaceholder(ph string) {
	b.Placeholder = ph
}

// GetPlaceholder returns the input placeholder.
func (b *Base) GetPlaceholder() string {
	return b.Placeholder
}

// SetTitle sets the input title (tooltip).
func (b *Base) SetTitle(title string) {
	b.Title = title
}

// GetTitle returns the input title.
func (b *Base) GetTitle() string {
	return b.Title
}

// SetSkipValidation sets whether to skip validation for this input.
func (b *Base) SetSkipValidation(skip bool) {
	b.SkipValidation = skip
}

// GetSkipValidation returns whether to skip validation.
func (b *Base) GetSkipValidation() bool {
	return b.SkipValidation
}

// SetOptions sets multiple options (for select/checkbox/etc).
func (b *Base) SetOptions(opts ...fmt.KeyValue) {
	b.Options = opts
}

// SetAliases sets the field name aliases for matching.
func (b *Base) SetAliases(aliases ...string) {
	b.aliases = aliases
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

// FieldName returns the struct field name (without parent prefix).
func (b *Base) FieldName() string {
	return b.name
}

// HTMLName returns the HTML input type.
func (b *Base) HTMLName() string {
	return b.htmlName
}

// Matches checks if the given field name matches this input's htmlName, name or aliases.
func (b *Base) Matches(fieldName string) bool {
	name := fmt.Convert(fieldName).ToLower().String()
	if b.htmlName == name || b.name == name {
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
	b.Attributes = append(b.Attributes, fmt.KeyValue{Key: key, Value: value})
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
