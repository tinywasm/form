package form_test

import (
	"github.com/tinywasm/fmt"
	"github.com/tinywasm/form"
	"github.com/tinywasm/form/input"
)

// User is a sample struct for testing data binding.
type User struct {
	Name     string
	Email    string
	Password string
	Gender   string
	Role     string
	Address  string
}

func (u *User) Schema() []fmt.Field {
	return []fmt.Field{
		{Name: "Name", Type: fmt.FieldText, NotNull: true, Widget: input.NewText()},
		{Name: "Email", Type: fmt.FieldText, NotNull: true, Widget: input.NewEmail()},
		{Name: "Password", Type: fmt.FieldText, NotNull: true, Widget: input.NewPassword()},
		{Name: "Gender", Type: fmt.FieldText, Widget: input.NewGender()},
		{Name: "Role", Type: fmt.FieldText, Widget: input.NewSelect()},
		{Name: "Address", Type: fmt.FieldText, Widget: input.NewAddress()},
	}
}

func (u *User) Values() []any {
	return []any{u.Name, u.Email, u.Password, u.Gender, u.Role, u.Address}
}

func (u *User) Pointers() []any {
	return []any{&u.Name, &u.Email, &u.Password, &u.Gender, &u.Role, &u.Address}
}

func (u *User) FormName() string {
	return "user"
}

// createTestForm is a helper to create a form for testing.
func createTestForm() *form.Form {
	u := &User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "secretpassword",
		Gender:   "m",
		Role:     "admin",
		Address:  "123 Main St",
	}
	f, err := form.New("test-parent", u)
	if err != nil {
		panic(err)
	}
	return f
}
