package form_test

import (
	"github.com/tinywasm/fmt"
	"github.com/tinywasm/form"
)

// User is a sample struct for testing data binding.
type User struct {
	Name     string // Auto: Text, validates, placeholder="Name" (via fmt.Translate)
	Email    string // Auto: Email, validates
	Password string // Auto: Password, validates
	Gender   string // Auto: Gender (m/f defaults)
	Role     string `options:"admin:Admin,user:User"` // Custom options
	Address  string `validate:"false"`                // No validation
}

func (u *User) Schema() []fmt.Field {
	return []fmt.Field{
		{Name: "Name", Type: fmt.FieldText, NotNull: true},
		{Name: "Email", Type: fmt.FieldText, NotNull: true, Input: "email"},
		{Name: "Password", Type: fmt.FieldText, NotNull: true, Input: "password"},
		{Name: "Gender", Type: fmt.FieldText},
		{Name: "Role", Type: fmt.FieldText},
		{Name: "Address", Type: fmt.FieldText, Input: "-"}, // testing exclusion as per plan? actually Address was "validate:false"
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
