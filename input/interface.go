package input

import (
	"github.com/tinywasm/dom"
)

// Input interface defines the behavior for all form input types.
// It embeds dom.Component to ensure compatibility with the tinywasm/dom ecosystem.
// RenderCSS and RenderJS are optional (implemented via type assertion if needed).
type Input interface {
	dom.Component // Includes ID() and RenderHTML()

	HtmlName() string                 // Standard HTML5 type (e.g., "text", "email")
	ValidateField(value string) error // Self-contained validation logic
}
