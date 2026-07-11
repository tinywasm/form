package form

import (
	"github.com/tinywasm/dom"
	"github.com/tinywasm/form/input"
)

// fieldComponent wraps an input.Input to implement dom.Component.
type fieldComponent struct {
	input.Input
	value *dom.SignalString
	err   *dom.SignalString
}

func (fc *fieldComponent) String() string {
	return fc.Render().String()
}

func (fc *fieldComponent) GetID() string {
	return fc.Input.GetID()
}

func (fc *fieldComponent) SetID(id string) {
	fc.Input.SetID(id)
}

func (fc *fieldComponent) Children() []dom.Component {
	return nil
}

// Renderer is an optional capability for custom inputs that own their markup.
// The form still owns the field wrapper (div.tw-field), the error span, and
// validation: the widget must call onInput with the new value on user input —
// the form updates the value signal and runs live validation. The value
// signal carries the initial value and programmatic updates (SetValues).
type Renderer interface {
	RenderInput(value *dom.SignalString, onInput func(string)) *dom.Element
}

func (fc *fieldComponent) validate(val string) {
	if err := fc.Input.Validate(val); err != nil {
		fc.err.Set(err.Error())
	} else {
		fc.err.Set("")
	}
}

func (fc *fieldComponent) Render() *dom.Element {
	container := dom.NewElement("div").Class("tw-field")

	if r, ok := fc.Input.(Renderer); ok {
		container.Child(r.RenderInput(fc.value, func(v string) {
			fc.value.Set(v)
			fc.validate(v)
		}))
	} else {
		htmlName := fc.Input.HTMLName()
		switch htmlName {
		case "radio":
			fc.renderRadio(container)
		case "select":
			fc.renderSelect(container)
		case "datalist":
			fc.renderDatalist(container)
		default:
			fc.renderInput(container)
		}
	}

	errSpan := dom.NewElement("span").
		ID(fc.Input.ErrorID()).
		Class("tw-field-error").
		Attr("aria-live", "polite").
		BindText(fc.err).
		BindClassFunc("tw-field-error--visible", func() bool {
			return fc.err.Get() != ""
		})

	container.Child(errSpan)
	return container
}

func (fc *fieldComponent) renderInput(container *dom.Element) {
	tag := "input"
	htmlName := fc.Input.HTMLName()
	if htmlName == "textarea" {
		tag = "textarea"
	}

	el := dom.NewElement(tag).
		ID(fc.Input.GetID()).
		Attr("name", fc.Input.FieldName())

	if tag == "input" {
		el.Attr("type", htmlName)
	}

	// Initial value for SSR
	val := fc.value.Get()
	if val != "" {
		if htmlName == "textarea" {
			el.Text(val)
		} else {
			el.Attr("value", val)
		}
	}

	// Two-way binding
	el.Bind(fc.value)
	el.On("input", func(e dom.Event) {
		val := e.TargetValue()
		fc.value.Set(val)
		fc.validate(val)
	})

	applyCommonAttrs(el, fc.Input)
	container.Child(el)
}

func (fc *fieldComponent) renderSelect(container *dom.Element) {
	el := dom.NewElement("select").
		ID(fc.Input.HandlerName()).
		Attr("name", fc.Input.FieldName())

	if fc.Input.IsRequired() {
		el.Attr("required", "")
	}

	val := fc.value.Get()

	// Two-way binding for select
	el.Bind(fc.value)
	el.On("change", func(e dom.Event) {
		val := e.TargetValue()
		fc.value.Set(val)
		fc.validate(val)
	})

	for _, opt := range fc.Input.GetOptions() {
		option := dom.NewElement("option").Attr("value", opt.Key).Text(opt.Value)
		if val != "" && opt.Key == val {
			option.Attr("selected", "")
		}
		el.Child(option)
	}
	container.Child(el)
}

func (fc *fieldComponent) renderRadio(container *dom.Element) {
	group := dom.NewElement("div").Class("tw-radio-group")
	val := fc.value.Get()
	for _, opt := range fc.Input.GetOptions() {
		optID := fc.Input.HandlerName() + "." + opt.Key
		label := dom.NewElement("label")

		radio := dom.NewElement("input").
			Attr("type", "radio").
			ID(optID).
			Attr("name", fc.Input.FieldName()).
			Attr("value", opt.Key)

		if val != "" && opt.Key == val {
			radio.Attr("checked", "")
		}

		// Reactive checked state
		radio.BindAttrBoolFunc("checked", func() bool {
			return fc.value.Get() == opt.Key
		})

		radio.On("change", func(e dom.Event) {
			if e.TargetChecked() {
				fc.value.Set(opt.Key)
				fc.validate(opt.Key)
			}
		})

		label.Child(radio)
		label.Child(dom.NewElement("span").Text(opt.Value))
		group.Child(label)
	}
	container.Child(group)
}

func (fc *fieldComponent) renderDatalist(container *dom.Element) {
	listID := fc.Input.GetID() + "-list"

	el := dom.NewElement("input").
		Attr("type", "text").
		ID(fc.Input.GetID()).
		Attr("name", fc.Input.FieldName()).
		Attr("list", listID)

	// Two-way binding
	el.Bind(fc.value)
	el.On("input", func(e dom.Event) {
		val := e.TargetValue()
		fc.value.Set(val)
		fc.validate(val)
	})

	applyCommonAttrs(el, fc.Input)
	container.Child(el)

	datalist := dom.NewElement("datalist").ID(listID)
	for _, opt := range fc.Input.GetOptions() {
		datalist.Child(dom.NewElement("option").Attr("value", opt.Key).Text(opt.Value))
	}
	container.Child(datalist)
}

func applyCommonAttrs(el *dom.Element, inp input.Input) {
	if ph := inp.GetPlaceholder(); ph != "" {
		el.Attr("placeholder", ph)
	}
	if title := inp.GetTitle(); title != "" {
		el.Attr("title", title)
	}
	for _, attr := range inp.GetAttributes() {
		if attr.Value != "" {
			el.Attr(attr.Key, attr.Value)
		}
	}
	if inp.IsRequired() {
		el.Attr("required", "")
	}
	if inp.IsDisabled() {
		el.Attr("disabled", "")
	}
	if inp.IsReadonly() {
		el.Attr("readonly", "")
	}
}

// RenderInput is kept for backward compatibility and as a standalone helper
func RenderInput(inp input.Input) *dom.Element {
	fc := &fieldComponent{
		Input: inp,
		value: dom.NewString(""),
		err:   dom.NewString(""),
	}
	// Note: this Render() returns a div.tw-field containing the input + error span.
	return fc.Render()
}
