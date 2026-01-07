package input

import (
	"github.com/tinywasm/dom"
)

// Input interface defines the behavior for all form input types.
// It embeds dom.Component to ensure compatibility with the tinywasm/dom ecosystem.
type Input interface {
	dom.Component // Includes ID() and RenderHTML()

	HTMLName() string                  // Standard HTML5 type (e.g., "text", "email")
	FieldName() string                 // Struct field name (without parent prefix)
	ValidateField(value string) error  // Self-contained validation logic
	Clone(parentID, name string) Input // Creates a new instance with given parentID and name
}
