package form

import (
	"github.com/tinywasm/dom"
	"github.com/tinywasm/fmt"
	"github.com/tinywasm/form/input"
)

// Form represents a form instance.
type Form struct {
	id           string
	parentID     string // Parent element ID where the form is mounted
	data         fmt.Fielder
	Inputs       []input.Input
	fieldIndices []int                   // Pre-computed struct field index per Input (-1 if not found)
	class        string                  // CSS class(es)
	method       string                  // HTTP method (default POST)
	action       string                  // Form action URL (default: struct name)
	ssrMode      bool                    // Per-form SSR mode (default false)
	submitLabel        string                            // Submit button label (empty = "Submit")
	submitLoadingLabel string                            // Label while submitting (default: label + "...")
	noResetOnSuccess   bool                              // Disable auto-reset after successful submit
	noSubmit           bool                              // True when the form should NOT render a submit button
	onSubmit           func(fmt.Fielder, func(error))    // WASM submit callback
	children           []dom.Component                   // Cached dom components (zero-alloc)
	valueSignals       []*dom.SignalString               // One per input
	errorSignals       []*dom.SignalString               // One per input
	submitting         *dom.SignalBool                   // Global form submitting state
}

// Children returns the form's input fields as dom components (O(1), zero-alloc).
func (f *Form) Children() []dom.Component {
	return f.children
}

// GetID returns the html id that group the form
func (f *Form) GetID() string {
	return f.id
}

// SetID sets the html id that group the form
func (f *Form) SetID(id string) {
	f.id = id
}

// ParentID returns the ID of the parent element.
func (f *Form) ParentID() string {
	return f.parentID
}

// OnSubmit sets the callback for form submission in WASM mode.
func (f *Form) OnSubmit(fn func(fmt.Fielder, func(error))) *Form {
	f.onSubmit = fn
	return f
}

// SubmitLoadingLabel customizes the text on the submit button while submitting.
func (f *Form) SubmitLoadingLabel(text string) *Form {
	f.submitLoadingLabel = text
	return f
}

// NoResetOnSuccess disables the automatic form reset after a successful submit.
func (f *Form) NoResetOnSuccess() *Form {
	f.noResetOnSuccess = true
	return f
}

// SubmitLabel customizes the text on the submit button.
// If never called, the button shows "Submit".
func (f *Form) SubmitLabel(text string) *Form {
	f.submitLabel = text
	return f
}

// HideSubmit disables rendering of the submit button.
// Use this when the form is part of a larger UI that provides its own
// submit control (e.g. an external toolbar). Default is to render one.
func (f *Form) HideSubmit() *Form {
	f.noSubmit = true
	return f
}

// Namer is optionally implemented by Fielder types to provide a custom name.
// If not implemented, the form derives the name from the first Schema field's context.
type Namer interface {
	FormName() string
}

func resolveStructName(data fmt.Fielder) string {
	if n, ok := data.(Namer); ok {
		return n.FormName()
	}
	return "form"
}

// New creates a new Form from a Fielder.
// parentID: ID of the parent DOM element where the form will be mounted.
// Returns an error if any exported field has no matching registered input.
func New(parentID string, data fmt.Fielder) (*Form, error) {
	schema := data.Schema()
	values := fmt.ReadValues(schema, data.Pointers())

	structName := resolveStructName(data)
	formID := parentID + "." + structName

	f := &Form{
		id:           formID,
		parentID:     parentID,
		data:         data,
		Inputs:       make([]input.Input, 0, len(schema)),
		class:        globalClass,
		method:       "POST",
		action:       "/" + structName,
		ssrMode:      false,
		children:     make([]dom.Component, 0, len(schema)),
		valueSignals: make([]*dom.SignalString, 0, len(schema)),
		errorSignals: make([]*dom.SignalString, 0, len(schema)),
		submitting:   dom.NewBool(false),
	}

	for i, field := range schema {
		// Skip auto-increment PKs (not editable)
		if field.IsPK() && field.IsAutoInc() {
			continue
		}

		fieldName := field.Name

		// Resolve input: use Field.Widget directly.
		if field.Widget == nil {
			continue // skip fields with no UI binding
		}

		inp, ok := field.Widget.Clone(formID, fieldName).(input.Input)
		if !ok {
			continue // skip widgets that do not implement input.Input
		}

		// Apply constraint-based defaults from schema
		if field.NotNull {
			if b, ok := inp.(interface{ SetRequired(bool) }); ok {
				b.SetRequired(true)
			}
		}

		// Initial value
		val := fmt.Convert(values[i]).String()

		// Bind current value to input (still needed for SSR/initial state)
		if setter, ok := inp.(interface{ SetValues(...string) }); ok {
			setter.SetValues(val)
		}

		vSig := dom.NewString(val)
		eSig := dom.NewString("")

		f.Inputs = append(f.Inputs, inp)
		f.valueSignals = append(f.valueSignals, vSig)
		f.errorSignals = append(f.errorSignals, eSig)
		f.children = append(f.children, &fieldComponent{inp, vSig, eSig})
		f.fieldIndices = append(f.fieldIndices, i)
	}

	forms = append(forms, f)
	return f, nil
}

// Input returns the input with the given field name, or nil if not found.
func (f *Form) Input(fieldName string) input.Input {
	for _, inp := range f.Inputs {
		if getter, ok := inp.(interface{ FieldName() string }); ok {
			if getter.FieldName() == fieldName {
				return inp
			}
		}
	}
	return nil
}

// SetOptions sets options for the input matching the given field name.
func (f *Form) SetOptions(fieldName string, opts ...fmt.KeyValue) *Form {
	inp := f.Input(fieldName)
	if inp != nil {
		if setter, ok := inp.(interface{ SetOptions(...fmt.KeyValue) }); ok {
			setter.SetOptions(opts...)
		}
	}
	return f
}

// Reset clears all input values and error messages in the DOM and internal state.
func (f *Form) Reset() { f.reset() }

func (f *Form) resolveSubmitLabel() string {
	label := f.submitLabel
	if label == "" {
		label = "Submit"
	}
	return label
}

func (f *Form) reset() {
	for i, inp := range f.Inputs {
		// Reset signals
		f.valueSignals[i].Set("")
		f.errorSignals[i].Set("")

		// Clear internal state (used by SSR/SyncValues if signals not available)
		if setter, ok := inp.(interface{ SetValues(...string) }); ok {
			setter.SetValues("")
		}
	}
}

// SetValues sets values for the input matching the given field name.
func (f *Form) SetValues(fieldName string, values ...string) *Form {
	for i, inp := range f.Inputs {
		if getter, ok := inp.(interface{ FieldName() string }); ok {
			if getter.FieldName() == fieldName {
				val := ""
				if len(values) > 0 {
					val = values[0]
				}
				f.valueSignals[i].Set(val)
				if setter, ok := inp.(interface{ SetValues(...string) }); ok {
					setter.SetValues(values...)
				}
				break
			}
		}
	}
	return f
}
