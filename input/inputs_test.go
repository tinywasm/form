package input

import (
	"testing"

	"github.com/tinywasm/fmt"
	_ "github.com/tinywasm/fmt/dictionary"
)

type testStruct struct {
	fieldType string // Name of the input struct instantiation
	name      string // Test case description
	inputData string // Input data for the test case
	expected  string // Expected substring in error message (case-insensitive)
	opts      []fmt.KeyValue
}

func Test_MigratedInputs(t *testing.T) {
	tests := []testStruct{
		// DataList tests
		{"Datalist", "Credencial válida (1)", "1", "", []fmt.KeyValue{{Key: "1", Value: "Admin"}, {Key: "3", Value: "Editor"}}},
		{"Datalist", "Credencial válida (3)", "3", "", []fmt.KeyValue{{Key: "1", Value: "Admin"}, {Key: "3", Value: "Editor"}}},
		{"Datalist", "Valor 0 no permitido", "0", "datalist_field", []fmt.KeyValue{{Key: "1", Value: "Admin"}, {Key: "3", Value: "Editor"}}},

		// Date tests
		{"Date", "Formato correcto", "2002-12-03", "", nil},
		{"Date", "Año bisiesto", "2020-02-29", "", nil},
		{"Date", "Año no bisiesto", "2023-02-29", "invalid", nil},
		{"Date", "Junio 31", "2023-06-31", "invalid", nil},
		{"Date", "Formato incorrecto", "21/12/1998", "not allowed", nil},
		{"Date", "Sin datos", "", "chars", nil},

		// FilePath tests
		{"Filepath", "Ruta correcta", ".\\files\\1234\\", "", nil},
		{"Filepath", "Ruta relativa con slash", "./files/1234/", "", nil},
		{"Filepath", "Número como ruta", "5", "", nil},
		{"Filepath", "Ruta sin punto inicial", "\\files\\1234\\", "\\", nil},
		{"Filepath", "Espacios en blanco", ".\\path with white space\\", "space", nil},

		// IP tests
		{"IP", "IPv4 válida", "192.168.1.1", "", nil},
		{"IP", "IPv6 válida", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", "", nil},
		{"IP", "Dirección 0.0.0.0", "0.0.0.0", "invalid", nil},
		{"IP", "Sin datos", "", "chars", nil},

		// Number tests
		{"Number", "Número correcto", "100", "", nil},
		{"Number", "Número negativo", "-100", "not allowed", nil},
		{"Number", "Texto en número", "lOO", "not allowed", nil},

		// Password tests
		{"Password", "Números y letras", "c0ntra3", "", nil},
		{"Password", "Sin datos", "", "chars", nil},
		{"Password", "Menos de 5 chars", "1", "chars", nil},

		// Rut tests
		{"Rut", "Sin guión", "15890022k", "hyphen", nil},
		{"Rut", "Sin guión (num)", "177344788", "hyphen", nil},
		{"Rut", "Correcto", "7863697-1", "", nil},
		{"Rut", "K mayúscula", "20373221-K", "", nil},
		{"Rut", "k minúscula", "20373221-k", "", nil},

		// Textarea tests
		{"Textarea", "Texto largo válido", "IRRITACION EN PIEL DE ROSTRO. ALERGIAS NO.", "", nil},
		{"Textarea", "Sin datos", "", "chars", nil},
		{"Textarea", "Solo espacio", " ", "chars", nil},

		// Checkbox tests
		{"Checkbox", "True value", "true", "", nil},
		{"Checkbox", "False value", "false", "", nil},
		{"Checkbox", "Valor vacío (no required)", "", "", nil},
		{"Checkbox", "Valor inválido", "hola", "invalid", nil},

		// Hour tests
		{"Hour", "Formato hh:mm", "12:30", "", nil},
		{"Hour", "Las 24 no existe", "24:00", "invalid", nil},
	}

	for _, tt := range tests {
		t.Run(tt.fieldType+"-"+tt.name, func(t *testing.T) {
			var inp Input
			id := "test_id"
			name := "test_name"

			switch tt.fieldType {
			case "Datalist":
				dl := Datalist(id, "datalist_field")
				if setter, ok := dl.(interface{ SetOptions(...fmt.KeyValue) }); ok && len(tt.opts) > 0 {
					setter.SetOptions(tt.opts...)
				}
				inp = dl
			case "Date":
				inp = Date(id, name)
			case "Filepath":
				inp = Filepath(id, name)
			case "IP":
				inp = IP(id, name)
			case "Number":
				inp = Number(id, name)
			case "Password":
				inp = Password(id, name)
			case "Rut":
				inp = Rut(id, name)
			case "Textarea":
				inp = Textarea(id, name)
			case "Checkbox":
				inp = Checkbox(id, name)
			case "Hour":
				inp = Hour(id, name)
			default:
				t.Fatalf("Unknown input type: %s", tt.fieldType)
			}

			err := inp.ValidateField(tt.inputData)
			var got string
			if err != nil {
				got = err.Error()
			}

			// We use lowecase Contains or direct match since exact translations might differ
			if tt.expected != "" {
				gotLower := fmt.Convert(got).ToLower().String()
				expLower := fmt.Convert(tt.expected).ToLower().String()
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.expected)
				} else if !fmt.Contains(gotLower, expLower) {
					t.Errorf("expected error containing %q, got %q", tt.expected, got)
				}
			} else if err != nil {
				t.Errorf("expected no error, got %q", got)
			}
		})
	}
}
