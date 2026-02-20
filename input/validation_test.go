package input

import "testing"

// Test_Validation runs all field validation cases.
// Add new cases here when adding inputs or new validation rules.
// Cases are sorted alphabetically by input type.
func Test_Validation(t *testing.T) {
	cases := []tc{
		// ── Address ──────────────────────────────────────────────────────────
		{"Address", "valid full", "Av. Libertad 123, Dpto #4", "", nil},
		{"Address", "too short", "Ab", "chars", nil},
		{"Address", "empty", "", "chars", nil},
		{"Address", "at sign", "user@home", "not allowed", nil},

		// ── Checkbox ─────────────────────────────────────────────────────────
		{"Checkbox", "true", "true", "", nil},
		{"Checkbox", "false", "false", "", nil},
		{"Checkbox", "on", "on", "", nil},
		{"Checkbox", "1", "1", "", nil},
		{"Checkbox", "0", "0", "", nil},
		{"Checkbox", "empty not required", "", "", nil},
		{"Checkbox", "invalid word", "yes", "invalid", nil},
		{"Checkbox", "invalid number", "2", "invalid", nil},

		// ── Datalist ─────────────────────────────────────────────────────────
		{"Datalist", "valid key 1", "1", "", opts12},
		{"Datalist", "valid key 2", "2", "", opts12},
		{"Datalist", "invalid key", "0", "notallowed", opts12},
		{"Datalist", "empty not required", "", "", opts12},
		{"Datalist", "no options always pass", "any", "", nil},

		// ── Date ─────────────────────────────────────────────────────────────
		{"Date", "valid date", "2002-12-03", "", nil},
		{"Date", "leap year 2020", "2020-02-29", "", nil},
		{"Date", "not leap year 2023", "2023-02-29", "invalid", nil},
		{"Date", "june has 30", "2023-06-31", "invalid", nil},
		{"Date", "month 13", "2023-13-01", "invalid", nil},
		{"Date", "slash format", "21/12/1998", "not allowed", nil},
		{"Date", "empty", "", "chars", nil},
		{"Date", "year too short", "200-01-01", "chars", nil},

		// ── Email ────────────────────────────────────────────────────────────
		{"Email", "valid", "user@example.com", "", nil},
		{"Email", "dots and dashes", "my.name-test@sub.domain.org", "", nil},
		{"Email", "empty", "", "chars", nil},
		{"Email", "too short", "a@b", "chars", nil},
		{"Email", "space", "user @mail.com", "not allowed", nil},

		// ── Filepath ─────────────────────────────────────────────────────────
		{"Filepath", "windows path", ".\\files\\1234\\", "", nil},
		{"Filepath", "unix path", "./files/1234/", "", nil},
		{"Filepath", "single char", "5", "", nil},
		{"Filepath", "leading backslash", "\\files\\", "\\", nil},
		{"Filepath", "space in path", ".\\path with space\\", "space", nil},
		{"Filepath", "empty", "", "chars", nil},

		// ── Gender ───────────────────────────────────────────────────────────
		{"Gender", "male", "m", "", nil},
		{"Gender", "female", "f", "", nil},
		{"Gender", "empty", "", "chars", nil},

		// ── Hour ─────────────────────────────────────────────────────────────
		{"Hour", "valid 12:30", "12:30", "", nil},
		{"Hour", "valid 00:00", "00:00", "", nil},
		{"Hour", "valid 23:59", "23:59", "", nil},
		{"Hour", "24:00 invalid", "24:00", "invalid", nil},
		{"Hour", "empty", "", "chars", nil},
		{"Hour", "no colon 4 chars", "1230", "chars", nil},

		// ── IP ───────────────────────────────────────────────────────────────
		{"IP", "valid ipv4", "192.168.1.1", "", nil},
		{"IP", "valid ipv6", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", "", nil},
		{"IP", "all zeros", "0.0.0.0", "invalid", nil},
		{"IP", "mixed dot and colon", "192.168:1.1", "invalid", nil},
		{"IP", "empty", "", "chars", nil},

		// ── Number ───────────────────────────────────────────────────────────
		{"Number", "valid 100", "100", "", nil},
		{"Number", "valid 0", "0", "", nil},
		{"Number", "negative", "-100", "not allowed", nil},
		{"Number", "letters", "abc", "not allowed", nil},
		{"Number", "empty", "", "chars", nil},

		// ── Password ─────────────────────────────────────────────────────────
		{"Password", "valid mixed", "c0ntra3!", "", nil},
		{"Password", "long", "MyP@ssw0rd123456", "", nil},
		{"Password", "empty", "", "chars", nil},
		{"Password", "too short", "a", "chars", nil},

		// ── Phone ────────────────────────────────────────────────────────────
		{"Phone", "valid intl", "+56 9 1234 5678", "", nil},
		{"Phone", "valid local", "912345678", "", nil},
		{"Phone", "too long", "+001-800-555-0123-99", "chars", nil},
		{"Phone", "too short", "123", "chars", nil},
		{"Phone", "empty", "", "chars", nil},
		{"Phone", "letter after min length", "+56-abc-123-456", "not allowed", nil},

		// ── Radio ────────────────────────────────────────────────────────────
		{"Radio", "valid key m", "m", "", nil},
		{"Radio", "valid key f", "f", "", nil},
		{"Radio", "empty", "", "chars", nil},

		// ── Rut ──────────────────────────────────────────────────────────────
		{"Rut", "valid 7863697-1", "7863697-1", "", nil},
		{"Rut", "valid K uppercase", "20373221-K", "", nil},
		{"Rut", "valid k lowercase", "20373221-k", "", nil},
		{"Rut", "no hyphen alpha", "15890022k", "hyphen", nil},
		{"Rut", "no hyphen digits", "177344788", "hyphen", nil},
		{"Rut", "wrong check digit", "7863697-2", "invalid", nil},
		{"Rut", "leading zero", "01234567-1", "invalid", nil},
		{"Rut", "empty", "", "chars", nil},

		// ── Select ───────────────────────────────────────────────────────────
		{"Select", "valid", "admin", "", nil},
		{"Select", "empty", "", "chars", nil},

		// ── Text ─────────────────────────────────────────────────────────────
		{"Text", "valid with tilde", "Juan Pérez", "", nil},
		{"Text", "valid with dots", "Dr. Smith (Jr.)", "", nil},
		{"Text", "too short", "A", "chars", nil},
		{"Text", "empty", "", "chars", nil},
		{"Text", "at sign", "user@domain", "not allowed", nil},

		// ── Textarea ─────────────────────────────────────────────────────────
		{"Textarea", "valid long", "IRRITACION EN PIEL DE ROSTRO. ALERGIAS NO.", "", nil},
		{"Textarea", "valid with newline", "Line one.\nLine two.", "", nil},
		{"Textarea", "empty", "", "chars", nil},
		{"Textarea", "too short", "Hi", "chars", nil},
	}

	for _, c := range cases {
		c := c
		t.Run(c.t+"/"+c.name, func(t *testing.T) {
			inp := buildInput(t, c.t, c.opts)
			checkErr(t, inp.ValidateField(c.val), c.err)
		})
	}
}
