package input

import (
	"github.com/tinywasm/dom"
	"github.com/tinywasm/fmt"
)

// Input interface defines the behavior for all form input types.
// It embeds dom.Component to ensure compatibility with the tinywasm/dom ecosystem.
type Input interface {
	dom.Component // Includes ID() and RenderHTML()

	HtmlName() string                  // Standard HTML5 type (e.g., "text", "email")
	ValidateField(value string) error  // Self-contained validation logic
	Clone(parentID, name string) Input // Creates a new instance with given parentID and name

	// State getters/setters
	GetValue() string
	GetValues() []string
	SetValues(...string)
	GetOptions() []fmt.KeyValue
	SetOptions(...fmt.KeyValue)
	GetSelectedValue() string
}
