package form_test

import (
	"testing"

	"github.com/tinywasm/fmt"
	"github.com/tinywasm/form"
	"github.com/tinywasm/form/input"
)

type testUser struct {
	id    string
	name  string
	email string
}

func (u *testUser) Schema() []fmt.Field {
	return []fmt.Field{
		{Name: "id", Type: fmt.FieldText, PK: true, Widget: input.Text("", "")},
		{Name: "name", Type: fmt.FieldText, NotNull: true, Widget: input.Text("", "")},
		{Name: "email", Type: fmt.FieldText, NotNull: true, Widget: input.Email("", "")},
	}
}
func (u *testUser) Values() []any    { return []any{u.id, u.name, u.email} }
func (u *testUser) Pointers() []any  { return []any{&u.id, &u.name, &u.email} }
func (u *testUser) FormName() string { return "user" }

func TestNewWithFielder(t *testing.T) {
	u := &testUser{id: "1", name: "John", email: "john@example.com"}
	f, err := form.New("parent", u)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	if f.GetID() != "parent.user" {
		t.Errorf("Expected form ID 'parent.user', got '%s'", f.GetID())
	}

	if len(f.Inputs) != 3 {
		t.Errorf("Expected 3 inputs, got %d", len(f.Inputs))
	}
}

type autoUser struct {
	id int64
}

func (u *autoUser) Schema() []fmt.Field {
	return []fmt.Field{
		{Name: "id", Type: fmt.FieldInt, PK: true, AutoInc: true},
	}
}
func (u *autoUser) Values() []any    { return []any{u.id} }
func (u *autoUser) Pointers() []any  { return []any{&u.id} }
func (u *autoUser) FormName() string { return "auto" }

func TestNewAutoIncPKExcludedReal(t *testing.T) {
	u := &autoUser{id: 1}
	f, err := form.New("parent", u)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}
	if len(f.Inputs) != 0 {
		t.Errorf("Expected 0 inputs (PK+AutoInc should be skipped), got %d", len(f.Inputs))
	}
}

func TestSyncValuesText(t *testing.T) {
	u := &testUser{id: "1", name: "Old"}
	f, err := form.New("p", u)
	if err != nil {
		t.Fatal(err)
	}

	f.SetValues("name", "New")
	err = f.SyncValues(u)
	if err != nil {
		t.Fatal(err)
	}

	if u.name != "New" {
		t.Errorf("Expected 'New', got '%s'", u.name)
	}
}

func TestValidateDataValid(t *testing.T) {
	u := &testUser{id: "100", name: "John Doe", email: "john@example.com"}
	f, err := form.New("p", u)
	if err != nil {
		t.Fatal(err)
	}

	err = f.ValidateData('u', u)
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
}
