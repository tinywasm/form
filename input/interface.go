package input

import (
	"github.com/tinywasm/dom"
	"github.com/tinywasm/fmt"
)

// Input interface defines the behavior for all form input types.
// It embeds dom.Component to ensure compatibility with the tinywasm/dom ecosystem.
type Input interface {
	fmt.Widget    // Type(), Validate(), Clone(parentID, name) — semantic type contract
	dom.Component // Includes GetID(), SetID(), RenderHTML(), Children()
}
