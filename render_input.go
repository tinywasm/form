package form

import (
	"github.com/tinywasm/dom"
	"github.com/tinywasm/fmt"
	"github.com/tinywasm/form/input"
	"github.com/tinywasm/html"
)

// fieldComponent wraps an input.Input to implement dom.Component.
type fieldComponent struct {
	input.Input
}

func (fc *fieldComponent) String() string {
	return RenderInput(fc.Input)
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

// RenderInput generates the HTML for the input based on its htmlName.
func RenderInput(inp input.Input) string {
	htmlName := inp.HTMLName()
	switch htmlName {
	case "select":
		return renderSelect(inp)
	case "radio":
		return renderRadio(inp)
	case "datalist":
		return renderDatalist(inp)
	}

	out := fmt.GetConv()

	var tag string
	var isInput bool

	if htmlName == "textarea" {
		tag = "textarea"
	} else {
		tag = "input"
		isInput = true
	}

	out.Write("<").Write(tag)
	if isInput {
		out.Write(` type="`).Write(htmlName).Write(`"`)
	}
	out.Write(` id="`).Write(inp.GetID()).Write(`"`)
	out.Write(` name="`).Write(inp.FieldName()).Write(`"`)

	values := inp.GetValues()
	value := ""
	if len(values) > 0 {
		value = values[0]
	}

	if isInput && value != "" {
		out.Write(` value="`).Write(value).Write(`"`)
	}
	if ph := inp.GetPlaceholder(); ph != "" {
		out.Write(` placeholder="`).Write(ph).Write(`"`)
	}
	if title := inp.GetTitle(); title != "" {
		out.Write(` title="`).Write(title).Write(`"`)
	}
	for _, attr := range inp.GetAttributes() {
		if attr.Value != "" {
			out.Write(` `).Write(attr.Key).Write(`="`).Write(attr.Value).Write(`"`)
		}
	}
	if inp.IsRequired() {
		out.Write(` required`)
	}
	if inp.IsDisabled() {
		out.Write(` disabled`)
	}
	if inp.IsReadonly() {
		out.Write(` readonly`)
	}

	if htmlName == "textarea" {
		out.Write(">")
		if value != "" {
			out.Write(value)
		}
		out.Write("</textarea>")
	} else {
		out.Write(">")
	}

	errSpan := html.Span("").ID(inp.ErrorID()).Class("tw-field-error").Attr("aria-live", "polite")
	out.Write(errSpan.String())

	return out.String()
}

// renderSelect generates <select> with <option> elements.
func renderSelect(inp input.Input) string {
	out := fmt.GetConv()
	values := inp.GetValues()
	out.Write(`<select id="`).Write(inp.HandlerName()).Write(`"`)
	out.Write(` name="`).Write(inp.FieldName()).Write(`"`)
	if inp.IsRequired() {
		out.Write(` required`)
	}
	out.Write(`>`)
	for _, opt := range inp.GetOptions() {
		out.Write(`<option value="`).Write(opt.Key).Write(`"`)
		for _, v := range values {
			if v == opt.Key {
				out.Write(` selected`)
				break
			}
		}
		out.Write(`>`).Write(opt.Value).Write(`</option>`)
	}
	out.Write(`</select>`)

	errSpan := html.Span("").ID(inp.ErrorID()).Class("tw-field-error").Attr("aria-live", "polite")
	out.Write(errSpan.String())

	return out.String()
}

// renderRadio generates <label><input type="radio"></label> per option.
func renderRadio(inp input.Input) string {
	out := fmt.GetConv()
	values := inp.GetValues()
	for _, opt := range inp.GetOptions() {
		optID := inp.HandlerName() + "." + opt.Key
		out.Write(`<label>`)
		out.Write(`<input type="radio" id="`).Write(optID).Write(`"`)
		out.Write(` name="`).Write(inp.FieldName()).Write(`"`)
		out.Write(` value="`).Write(opt.Key).Write(`"`)
		for _, v := range values {
			if v == opt.Key {
				out.Write(` checked`)
				break
			}
		}
		out.Write(`>`)
		out.Write(opt.Value)
		out.Write(`</label>`)
	}

	errSpan := html.Span("").ID(inp.ErrorID()).Class("tw-field-error").Attr("aria-live", "polite")
	out.Write(errSpan.String())

	return out.String()
}

// renderDatalist generates <input> linked to a <datalist> element.
func renderDatalist(inp input.Input) string {
	listID := inp.GetID() + "-list"
	inp.AddAttribute("list", listID)
	out := fmt.GetConv()

	// Datalist always uses a text input as the base
	// We need a way to render it as text.
	// Since we can't temporarily change htmlName easily (it's internal to input.Base),
	// we'll manually render the text input part here or use a helper.

	out.Write(`<input type="text" id="`).Write(inp.GetID()).Write(`"`)
	out.Write(` name="`).Write(inp.FieldName()).Write(`"`)

	values := inp.GetValues()
	if len(values) > 0 && values[0] != "" {
		out.Write(` value="`).Write(values[0]).Write(`"`)
	}
	if ph := inp.GetPlaceholder(); ph != "" {
		out.Write(` placeholder="`).Write(ph).Write(`"`)
	}
	if title := inp.GetTitle(); title != "" {
		out.Write(` title="`).Write(title).Write(`"`)
	}
	for _, attr := range inp.GetAttributes() {
		if attr.Value != "" {
			out.Write(` `).Write(attr.Key).Write(`="`).Write(attr.Value).Write(`"`)
		}
	}
	if inp.IsRequired() {
		out.Write(` required`)
	}
	if inp.IsDisabled() {
		out.Write(` disabled`)
	}
	if inp.IsReadonly() {
		out.Write(` readonly`)
	}
	out.Write(">")

	errSpan := html.Span("").ID(inp.ErrorID()).Class("tw-field-error").Attr("aria-live", "polite")
	out.Write(errSpan.String())

	out.Write(`<datalist id="`).Write(listID).Write(`">`)
	for _, opt := range inp.GetOptions() {
		out.Write(`<option value="`).Write(opt.Key).Write(`">`).Write(opt.Value).Write(`</option>`)
	}
	out.Write(`</datalist>`)
	return out.String()
}
