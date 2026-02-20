package form_test

import (
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
