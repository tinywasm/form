package form_test

import (
	"testing"

	"github.com/tinywasm/dom"
	"github.com/tinywasm/form"
	"github.com/tinywasm/model"
)

type mockFielder struct {
}

func (m *mockFielder) Schema() []model.Field {
	return []model.Field{}
}

func (m *mockFielder) Pointers() []any {
	return []any{}
}

func (m *mockFielder) Values() []any {
	return []any{}
}

func TestForm_DomComponent(t *testing.T) {
	f, _ := form.New("parent", &mockFielder{})

	// This assignment will fail to compile if *Form does not implement dom.Component
	var _ dom.Component = f
}
