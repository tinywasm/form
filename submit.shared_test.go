package form

import (
	"testing"
	"github.com/tinywasm/fmt"
	"github.com/tinywasm/form/input"
)

type submitStruct struct {
	fmt.Fielder
	Nombre string `input:"required"`
}

func (s *submitStruct) Schema() []fmt.Field {
	return []fmt.Field{
		{Name: "nombre", NotNull: true, Widget: input.Text()},
	}
}

func (s *submitStruct) Pointers() []any { return []any{&s.Nombre} }
func (s *submitStruct) Values() []any   { return []any{s.Nombre} }

func runSubmitTests(t *testing.T) {
	t.Run("TestSubmit_CallbackReceivesDone", func(t *testing.T) {
		s := &submitStruct{Nombre: "Jules"}
		f, _ := New("app", s)

		called := false
		f.OnSubmit(func(data fmt.Fielder, done func(error)) {
			called = true
			done(nil)
		})

		if f.onSubmit == nil {
			t.Fatal("onSubmit callback not set")
		}

		f.onSubmit(f.data, func(err error) {
			// dummy done
		})

		if !called {
			t.Error("OnSubmit callback was not called")
		}
	})

	t.Run("TestSubmit_ResetClearsValues", func(t *testing.T) {
		s := &submitStruct{Nombre: "Jules"}
		f, _ := New("app", s)
		f.SetValues("nombre", "New Value")

		f.Reset()

		val := f.Input("nombre").(interface{ GetValue() string }).GetValue()
		if val != "" {
			t.Errorf("Expected empty value after reset, got %q", val)
		}
	})
}
