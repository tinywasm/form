package form_test

import (
	"github.com/tinywasm/form"
	"github.com/tinywasm/form/input"
	"github.com/tinywasm/model"
	"testing"
)

type submitStruct struct {
	model.Fielder
	Nombre string
}

func (s *submitStruct) Schema() []model.Field {
	return []model.Field{
		{Name: "nombre", NotNull: true, Type: input.Text()},
	}
}

func (s *submitStruct) Pointers() []any { return []any{&s.Nombre} }
func (s *submitStruct) Values() []any   { return []any{s.Nombre} }

func runSubmitTests(t *testing.T) {
	t.Run("TestSubmit_CallbackReceivesDataAndDone", func(t *testing.T) {
		s := &submitStruct{Nombre: "Jules"}
		f, _ := form.New("app", s)
		f.SetValues("nombre", "New Name")

		called := false
		var receivedData model.Fielder
		f.OnSubmit(func(data model.Fielder, done func(error)) {
			called = true
			receivedData = data
			done(nil)
		})

		err := f.Submit()
		if err != nil {
			t.Fatalf("Submit failed: %v", err)
		}

		if !called {
			t.Error("OnSubmit callback was not called")
		}

		if receivedData != s {
			t.Errorf("Expected to receive bound struct, got %v", receivedData)
		}

		if s.Nombre != "New Name" {
			t.Errorf("Expected struct field to be synced, got %q", s.Nombre)
		}
	})

	t.Run("TestSubmit_ValidationFailureReturnsError", func(t *testing.T) {
		s := &submitStruct{Nombre: ""} // empty, but NotNull: true
		f, _ := form.New("app", s)
		f.SetValues("nombre", "")

		called := false
		f.OnSubmit(func(data model.Fielder, done func(error)) {
			called = true
			done(nil)
		})

		err := f.Submit()
		if err == nil {
			t.Fatal("Expected validation error, got nil")
		}

		if called {
			t.Error("OnSubmit callback should NOT have been called on validation failure")
		}
	})

	t.Run("TestSubmit_NoResetOnSuccess", func(t *testing.T) {
		s := &submitStruct{Nombre: "Jules"}
		f, _ := form.New("app", s)
		f.NoResetOnSuccess()
		f.SetValues("nombre", "Keep Me")

		f.OnSubmit(func(data model.Fielder, done func(error)) {
			done(nil)
		})

		err := f.Submit()
		if err != nil {
			t.Fatalf("Submit failed: %v", err)
		}

		val := f.Input("nombre").(interface{ GetValue() string }).GetValue()
		if val != "Keep Me" {
			t.Errorf("Expected form to retain value 'Keep Me', got %q", val)
		}
	})

	t.Run("TestSubmit_DefaultResetOnSuccess", func(t *testing.T) {
		s := &submitStruct{Nombre: "Jules"}
		f, _ := form.New("app", s)
		f.SetValues("nombre", "Clear Me")

		f.OnSubmit(func(data model.Fielder, done func(error)) {
			done(nil)
		})

		err := f.Submit()
		if err != nil {
			t.Fatalf("Submit failed: %v", err)
		}

		val := f.Input("nombre").(interface{ GetValue() string }).GetValue()
		if val != "" {
			t.Errorf("Expected empty value after successful submit, got %q", val)
		}
	})

	t.Run("TestSubmit_ResetClearsValues", func(t *testing.T) {
		s := &submitStruct{Nombre: "Jules"}
		f, _ := form.New("app", s)
		f.SetValues("nombre", "New Value")

		f.Reset()

		val := f.Input("nombre").(interface{ GetValue() string }).GetValue()
		if val != "" {
			t.Errorf("Expected empty value after reset, got %q", val)
		}
	})
}
