package form

import "github.com/tinywasm/model"

import (
	"github.com/tinywasm/dom"
	"github.com/tinywasm/fmt"
	"github.com/tinywasm/form/input"
)

// Form represents a form instance.
type Form struct {
	id                 string
	parentID           string // Parent element ID where the form is mounted
	data               model.Fielder
	Inputs             []input.Input
	fieldIndices       []int                            // Pre-computed struct field index per Input (-1 if not found)
	class              string                           // CSS class(es)
	method             string                           // HTTP method (default POST)
	action             string                           // Form action URL (default: struct name)
	ssrMode            bool                             // Per-form SSR mode (default false)
	submitLabel        string                           // Submit button label (empty = "Submit")
	submitLoadingLabel string                           // Label while submitting (default: label + "...")
	noResetOnSuccess   bool                             // Disable auto-reset after successful submit
	noSubmit           bool                             // True when the form should NOT render a submit button
	onSubmit           func(model.Fielder, func(error)) // WASM submit callback
	onFieldChange      func()                           // fires when a field is committed (blur/change) — auto-save hook
	children           []dom.Component                  // Cached dom components (zero-alloc)
	valueSignals       []*dom.SignalString              // One per input
	errorSignals       []*dom.SignalString              // One per input
	submitting         *dom.SignalBool                  // Global form submitting state
	locked             *dom.SignalBool                  // Whole-form read-only gate (see SetLocked)
	focused            string                           // id Focus() last targeted (see FocusedFieldID)
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
func (f *Form) OnSubmit(fn func(model.Fielder, func(error))) *Form {
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

// SetLocked gates every field to read-only/disabled (whole-form, not per-field)
// without discarding their values — used by a host UI to show an existing
// record before an explicit "edit" action unlocks it. Reactive: takes effect
// immediately on an already-rendered form.
func (f *Form) SetLocked(v bool) *Form {
	f.locked.Set(v)
	return f
}

// Focus moves keyboard focus to the form's first field — a host UI calls this
// when entering an editable state (e.g. crudview's "+" / ⋮ Editar) so the user
// can start typing immediately instead of having to click into the form. A
// no-op if the form has no fields. Imperative, not reactive: the form's DOM
// already exists by the time a host unlocks it (this never runs on first
// mount), so a direct dom.Get+Focus is enough — no binding needed.
func (f *Form) Focus() *Form {
	if len(f.Inputs) == 0 {
		return f
	}
	f.focused = f.Inputs[0].GetID()
	if ref, ok := dom.Get(f.focused); ok {
		ref.Focus()
	}
	return f
}

// FocusedFieldID returns the id Focus() last targeted (empty if never called,
// or the form has no fields). Real focus movement is a WASM-only DOM side
// effect (a no-op in the backend/SSR stub); this makes the INTENT observable
// in any build, e.g. for the view/conformance "New/Edit focuses the first
// field" clause to assert against without a live DOM.
func (f *Form) FocusedFieldID() string { return f.focused }

// OnFieldChange registers a callback fired every time a field is committed by the
// user: blur for text/textarea/datalist, change for select/radio. This is the
// hook a host uses for auto-save (no explicit Save button) — the callback runs
// AFTER the field's own value/validate update, so the form's data is current.
func (f *Form) OnFieldChange(fn func()) *Form {
	f.onFieldChange = fn
	return f
}

// SetClass appends CSS classes to this form (on top of any global classes
// set via SetGlobalClass). Chainable.
func (f *Form) SetClass(classes ...string) *Form {
	for _, c := range classes {
		if f.class != "" {
			f.class += " "
		}
		f.class += c
	}
	return f
}

// Namer is optionally implemented by Fielder types to provide a custom name.
// If not implemented, the form derives the name from the first Schema field's context.
type Namer interface {
	FormName() string
}

func resolveStructName(data model.Fielder) string {
	if n, ok := data.(Namer); ok {
		return n.FormName()
	}
	return "form"
}

// New creates a new Form from a Fielder.
// parentID: ID of the parent DOM element where the form will be mounted.
// Returns an error if any exported field has no matching registered input.
func New(parentID string, data model.Fielder) (*Form, error) {
	schema := data.Schema()
	values := model.ReadValues(schema, data.Pointers())

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
		locked:       dom.NewBool(false),
	}

	for i, field := range schema {
		// Skip auto-increment PKs (not editable)
		if field.IsPK() && field.IsAutoInc() {
			continue
		}

		fieldName := field.Name

		// Resolve input: use Field.Type asserting to input.Input.
		inpKind, ok := field.Type.(input.Input)
		if !ok {
			continue // skip fields with no UI binding
		}

		inp := inpKind.Clone(formID, fieldName)

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
		// A closure, not f.onFieldChange by value: OnFieldChange is meant to be
		// called AFTER New() returns (chainable, like HideSubmit) — capturing the
		// field directly here would freeze it at nil since registration happens
		// later. The closure re-reads f.onFieldChange at commit time instead.
		f.children = append(f.children, &fieldComponent{inp, vSig, eSig, f.locked, func() {
			if f.onFieldChange != nil {
				f.onFieldChange()
			}
		}})
		f.fieldIndices = append(f.fieldIndices, i)
	}

	if len(f.Inputs) == 0 {
		return nil, fmt.Errf("form.New: %s has no renderable field — every Field.Type is a "+
			"plain model.Kind, not a form input.Input. Declare the widget in the model "+
			"Definition (input.Text(), input.Number(), …) instead of model.Text()/model.Int()",
			structName)
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

// Submit runs the full submit pipeline programmatically: syncs input values
// into the bound struct, validates, and (if valid) fires the OnSubmit
// callback. Returns the first validation error, or nil if the submission
// was dispatched. The async result of the submission itself is delivered
// through the OnSubmit callback's done function.
func (f *Form) Submit() error {
	// Sync all values from signals to struct
	f.SyncValues(f.data)

	// Validate all (final check)
	if err := f.Validate(); err != nil {
		return err
	}

	if f.onSubmit != nil {
		f.submitting.Set(true)
		f.onSubmit(f.data, func(err error) {
			f.submitting.Set(false)
			if err == nil && !f.noResetOnSuccess {
				f.reset()
			}
		})
	}
	return nil
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
