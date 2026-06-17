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
	Values         []string       // Multiple values support (for select/checkbox/etc)
	Options        []fmt.KeyValue // Multiple options for select/checkbox/etc
	Placeholder    string
	Title          string
	Required       bool // HTML required attribute
	Disabled       bool // HTML disabled attribute
	Readonly       bool // HTML readonly attribute
	SkipValidation bool // Whether to skip validation for this input
	Attributes     []fmt.KeyValue
	fmt.Permitted  // anonymous embed: promotes Letters, Numbers, Validate(), etc.
}

// InitBase initializes the base fields and constructs the unique ID.
func (b *Base) InitBase(parentID, name, htmlName string) {
	if parentID != "" {
		b.id = parentID + "." + name
	} else {
		b.id = name
	}
	b.name = name
	b.htmlName = htmlName

	// Only apply defaults if not already set (preserves values during Clone)
	if b.Placeholder == "" {
		b.Placeholder = fmt.Translate(name).String()
	}
	if b.Title == "" {
		b.Title = fmt.Translate(name).String()
	}
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

// GetOptions returns all options.
func (b *Base) GetOptions() []fmt.KeyValue {
	return b.Options
}

// GetSelectedValue returns the first value in Values, or empty if none.
func (b *Base) GetSelectedValue() string {
	return b.GetValue()
}

// GetID returns the component's unique identifier.
func (b *Base) GetID() string {
	return b.id
}

// ErrorID returns the ID of the associated error span for this input.
func (b *Base) ErrorID() string { return b.id + ".error" }

// SetID sets the component's unique identifier.
func (b *Base) SetID(id string) {
	b.id = id
}

// Type satisfies fmt.Widget.Type(). Returns the semantic input type name.
func (b *Base) Type() string { return b.htmlName }

// Validate satisfies fmt.Widget.Validate().
// Concrete structs embedding Base can override this to provide specialized validation.
func (b *Base) Validate(value string) error {
	if value == "" && b.Required {
		return fmt.Err("field", b.name, "is required")
	}
	return b.Permitted.Validate(b.name, value)
}

// SetRequired sets the required attribute.
func (b *Base) SetRequired(req bool) {
	b.Required = req
}

// IsRequired returns true if the input is required.
func (b *Base) IsRequired() bool {
	return b.Required
}

// IsDisabled returns true if the input is disabled.
func (b *Base) IsDisabled() bool {
	return b.Disabled
}

// IsReadonly returns true if the input is readonly.
func (b *Base) IsReadonly() bool {
	return b.Readonly
}

// GetAttributes returns custom attributes.
func (b *Base) GetAttributes() []fmt.KeyValue {
	return b.Attributes
}

// HandlerName returns the component's unique identifier.
// Deprecated: use GetID instead.
func (b *Base) HandlerName() string {
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

// AddAttribute adds a custom attribute to the input.
func (b *Base) AddAttribute(key, value string) {
	b.Attributes = append(b.Attributes, fmt.KeyValue{Key: key, Value: value})
}
