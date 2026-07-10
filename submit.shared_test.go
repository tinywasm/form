package form

import "github.com/tinywasm/model"

import (
	"github.com/tinywasm/form/input"
	"testing"
)

type submitStruct struct {
	model.Fielder
	Nombre string `input:"required"`
}

func (s *submitStruct) Schema() []model.Field {
	return []model.Field{
		{Name: "nombre", NotNull: true, Type: input.Text()},
	}
}

func (s *submitStruct) Pointers() []any { return []any{&s.Nombre} }
func (s *submitStruct) Values() []any   { return []any{s.Nombre} }

func runSubmitTests(t *testing.T) {
	t.Run("TestSubmit_CallbackReceivesDone", func(t *testing.T) {
		s := &submitStruct{Nombre: "Jules"}
		f, _ := New("app", s)

		called := false
		f.OnSubmit(func(data model.Fielder, done func(error)) {
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

	t.Run("TestSubmit_NoResetOnSuccess", func(t *testing.T) {
		s := &submitStruct{Nombre: "Jules"}
		f, _ := New("app", s)
		f.NoResetOnSuccess()
		f.SetValues("nombre", "Valor Original")

		doneCalled := false
		f.onSubmit = func(data model.Fielder, done func(error)) {
			done(nil)
		}
		f.onSubmit(f.data, func(err error) {
			doneCalled = true
			if err == nil && !f.noResetOnSuccess {
				f.reset()
			}
		})

		if !doneCalled {
			t.Error("done callback was not called")
		}

		val := f.Input("nombre").(interface{ GetValue() string }).GetValue()
		if val == "" {
			t.Error("Expected form to retain values when NoResetOnSuccess is set")
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
