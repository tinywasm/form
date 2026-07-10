package input
import "github.com/tinywasm/model"


import (
	"github.com/tinywasm/fmt"
)

// Input interface defines the behavior for all form input types.
type Input interface {
	model.Kind
	Clone(parentID, name string) Input
	FieldName() string
	SetRequired(bool)
	AddAttribute(key, value string)

	// Metadata getters for rendering
	GetID() string
	SetID(string)
	GetValues() []string
	GetOptions() []fmt.KeyValue
	GetPlaceholder() string
	GetTitle() string
	IsRequired() bool
	IsDisabled() bool
	IsReadonly() bool
	GetAttributes() []fmt.KeyValue
	ErrorID() string
	HTMLName() string
	HandlerName() string
}
