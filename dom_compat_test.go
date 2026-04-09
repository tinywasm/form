package form

import (
	"testing"

	"github.com/tinywasm/dom"
	"github.com/tinywasm/fmt"
)

type mockFielder struct {
}

func (m *mockFielder) Schema() []fmt.Field {
	return []fmt.Field{}
}

func (m *mockFielder) Pointers() []any {
	return []any{}
}

func (m *mockFielder) Values() []any {
	return []any{}
}

func TestForm_DomComponent(t *testing.T) {
	f, _ := New("parent", &mockFielder{})

	// This assignment will fail to compile if *Form does not implement dom.Component
	var _ dom.Component = f
}
