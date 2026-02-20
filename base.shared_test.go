package form_test

import (
	"testing"

	"github.com/tinywasm/fmt"
	"github.com/tinywasm/form"
	"github.com/tinywasm/form/input"
)

func TestForm_NewAndBinding_Shared(t *testing.T) {
	f := createTestForm()

	// Helper to get value via type assertion
	getValue := func(inp interface{}) string {
		if g, ok := inp.(interface{ GetValue() string }); ok {
			return g.GetValue()
		}
		return ""
	}

	// Check if values were bound correctly
	nameInp := f.Input("Name")
	if nameInp == nil {
		t.Fatal("Input 'Name' not found")
	}
	if getValue(nameInp) != "John Doe" {
		t.Errorf("Expected 'John Doe', got '%s'", getValue(nameInp))
	}

	emailInp := f.Input("Email")
	if getValue(emailInp) != "john@example.com" {
		t.Errorf("Expected 'john@example.com', got '%s'", getValue(emailInp))
	}

	genderInp := f.Input("Gender")
	if getValue(genderInp) != "m" {
		t.Errorf("Expected 'm', got '%s'", getValue(genderInp))
	}
}

func TestForm_RenderHTML_SSR_Shared(t *testing.T) {
	f := createTestForm()
	f.SetSSR(true)

	html := f.RenderHTML()

	// Verify SSR specific attributes
	if !fmt.Contains(html, `method="POST"`) {
		t.Error("SSR form should contain method=\"POST\"")
	}

	// Verify inputs are present
	if !fmt.Contains(html, `name="Name"`) {
		t.Error("SSR form should contain input name=\"Name\"")
	}
	if !fmt.Contains(html, `value="John Doe"`) {
		t.Error("SSR form should contain value=\"John Doe\"")
	}
}

func TestForm_AutoDefaults_Shared(t *testing.T) {
	f := createTestForm()

	// Placeholder and title default to the translated field name via fmt.Translate.
	// Without a registered translation for "Name", it passes through as-is.
	nameInp := f.Input("Name")
	if p, ok := nameInp.(interface{ GetPlaceholder() string }); ok {
		if p.GetPlaceholder() != "Name" {
			t.Errorf("Expected 'Name', got '%s'", p.GetPlaceholder())
		}
	}

	if p, ok := nameInp.(interface{ GetTitle() string }); ok {
		if p.GetTitle() != "Name" {
			t.Errorf("Expected 'Name', got '%s'", p.GetTitle())
		}
	}
}

func TestForm_Validate_Shared(t *testing.T) {
	f := createTestForm()

	// Valid form
	if err := f.Validate(); err != nil {
		t.Errorf("Expected valid form, got error: %v", err)
	}

	// Invalid form
	f.SetValues("Name", "") // Too short (min 2)
	if err := f.Validate(); err == nil {
		t.Error("Expected validation error for empty name, got nil")
	}

	// Reset Name to valid value
	f.SetValues("Name", "John Doe")

	// Skip validation test (Address has validate:"false")
	f.SetValues("Address", "") // Address is Text, normally needs min 5 chars
	if err := f.Validate(); err != nil {
		t.Errorf("Expected Address to skip validation, got error: %v", err)
	}
}

func TestForm_CustomInput_Shared(t *testing.T) {
	// 1. Create a custom input
	custom := input.Text("", "custom_field")
	// Add an alias that includes the struct name
	if b, ok := custom.(interface{ SetAliases(...string) }); ok {
		b.SetAliases("customuser.special")
	}
	form.RegisterInput(custom)

	// 2. Struct with a field that should match the custom input
	type CustomUser struct {
		Special string
	}

	f, err := form.New("parent", &CustomUser{})
	if err != nil {
		t.Fatalf("Failed to create form with custom input: %v", err)
	}

	if f.Input("Special") == nil {
		t.Error("Custom input 'Special' not found via StructName.FieldName matching")
	}
}

func TestForm_ValidateData_Shared(t *testing.T) {
	f := createTestForm()

	// Valid struct: all fields filled with valid values — should pass
	valid := &User{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "secret123",
		Gender:   "m",
		Role:     "admin",
		Address:  "123 Main St",
	}
	if err := f.ValidateData('c', valid); err != nil {
		t.Errorf("Expected valid data to pass, got error: %v", err)
	}

	// Invalid email — should fail
	invalid := &User{
		Name:     "John",
		Email:    "not-an-email",
		Password: "secret123",
	}
	if err := f.ValidateData('c', invalid); err == nil {
		t.Error("Expected invalid email to fail ValidateData, got nil")
	}

	// No data — should return nil
	if err := f.ValidateData('c'); err != nil {
		t.Errorf("Expected no-data call to return nil, got: %v", err)
	}
}

// runSharedTests executes all test cases common to both WASM and Standard Lib.
func runSharedTests(t *testing.T) {
	t.Run("NewAndBinding", TestForm_NewAndBinding_Shared)
	t.Run("AutoDefaults", TestForm_AutoDefaults_Shared)
	t.Run("CustomInput", TestForm_CustomInput_Shared)
	t.Run("RenderHTML_SSR", TestForm_RenderHTML_SSR_Shared)
	t.Run("Validate", TestForm_Validate_Shared)
	t.Run("ValidateData", TestForm_ValidateData_Shared)
}
