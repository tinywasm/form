package input

import "testing"

// Test_Validation runs all field validation cases.
// Add new cases here when adding inputs or new validation rules.
// Cases are sorted alphabetically by input type.
func Test_Validation(t *testing.T) {
	cases := []tc{
		// ── Address ──────────────────────────────────────────────────────────
		{"Address", "valid full", "Av. Libertad 123, Dpto #4", "", nil, false},
		{"Address", "too short", "Ab", "chars", nil, false},
		{"Address", "empty", "", "chars", nil, false},
		{"Address", "at sign", "user@home", "not allowed", nil, false},

		// ── Checkbox ─────────────────────────────────────────────────────────
		{"Checkbox", "true", "true", "", nil, false},
		{"Checkbox", "false", "false", "", nil, false},
		{"Checkbox", "on", "on", "", nil, false},
		{"Checkbox", "1", "1", "", nil, false},
		{"Checkbox", "0", "0", "", nil, false},
		{"Checkbox", "empty not required", "", "", nil, false},
		{"Checkbox", "empty required", "", "empty", nil, true},
		{"Checkbox", "invalid word", "yes", "invalid", nil, false},
		{"Checkbox", "invalid number", "2", "invalid", nil, false},

		// ── Datalist ─────────────────────────────────────────────────────────
		{"Datalist", "valid key 1", "1", "", opts12, false},
		{"Datalist", "valid key 2", "2", "", opts12, false},
		{"Datalist", "invalid key", "0", "notallowed", opts12, false},
		{"Datalist", "empty not required", "", "", opts12, false},
		{"Datalist", "no options always pass", "any", "", nil, false},

		// ── Date ─────────────────────────────────────────────────────────────
		{"Date", "valid date", "2002-12-03", "", nil, false},
		{"Date", "leap year 2020", "2020-02-29", "", nil, false},
		{"Date", "not leap year 2023", "2023-02-29", "invalid", nil, false},
		{"Date", "june has 30", "2023-06-31", "invalid", nil, false},
		{"Date", "month 13", "2023-13-01", "invalid", nil, false},
		{"Date", "slash format", "21/12/1998", "not allowed", nil, false},
		{"Date", "empty", "", "chars", nil, false},
		{"Date", "year too short", "200-01-01", "chars", nil, false},

		// ── Email ────────────────────────────────────────────────────────────
		{"Email", "valid", "user@example.com", "", nil, false},
		{"Email", "dots and dashes", "my.name-test@sub.domain.org", "", nil, false},
		{"Email", "empty", "", "chars", nil, false},
		{"Email", "too short", "a@b", "chars", nil, false},
		{"Email", "space", "user @mail.com", "not allowed", nil, false},

		// ── Filepath ─────────────────────────────────────────────────────────
		{"Filepath", "windows path", ".\\files\\1234\\", "", nil, false},
		{"Filepath", "unix path", "./files/1234/", "", nil, false},
		{"Filepath", "single char", "5", "", nil, false},
		{"Filepath", "leading backslash", "\\files\\", "\\", nil, false},
		{"Filepath", "space in path", ".\\path with space\\", "space", nil, false},
		{"Filepath", "empty", "", "chars", nil, false},

		// ── Gender ───────────────────────────────────────────────────────────
		{"Gender", "male", "m", "", nil, false},
		{"Gender", "female", "f", "", nil, false},
		{"Gender", "empty", "", "chars", nil, false},

		// ── Hour ─────────────────────────────────────────────────────────────
		{"Hour", "valid 12:30", "12:30", "", nil, false},
		{"Hour", "valid 00:00", "00:00", "", nil, false},
		{"Hour", "valid 23:59", "23:59", "", nil, false},
		{"Hour", "24:00 invalid", "24:00", "invalid", nil, false},
		{"Hour", "empty", "", "", nil, false},
		{"Hour", "no colon 4 chars", "1230", "invalid", nil, false},

		// ── IP ───────────────────────────────────────────────────────────────
		{"IP", "valid ipv4", "192.168.1.1", "", nil, false},
		{"IP", "valid ipv6", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", "", nil, false},
		{"IP", "all zeros", "0.0.0.0", "invalid", nil, false},
		{"IP", "mixed dot and colon", "192.168:1.1", "invalid", nil, false},
		{"IP", "empty", "", "chars", nil, false},

		// ── Number ───────────────────────────────────────────────────────────
		{"Number", "valid 100", "100", "", nil, false},
		{"Number", "valid 0", "0", "", nil, false},
		{"Number", "negative", "-100", "not allowed", nil, false},
		{"Number", "letters", "abc", "not allowed", nil, false},
		{"Number", "empty", "", "chars", nil, false},

		// ── Password ─────────────────────────────────────────────────────────
		{"Password", "valid mixed", "c0ntra3!", "", nil, false},
		{"Password", "long", "MyP@ssw0rd123456", "", nil, false},
		{"Password", "empty", "", "chars", nil, false},
		{"Password", "too short", "a", "chars", nil, false},

		// ── Phone ────────────────────────────────────────────────────────────
		{"Phone", "valid intl", "+56 9 1234 5678", "", nil, false},
		{"Phone", "valid local", "912345678", "", nil, false},
		{"Phone", "too long", "+001-800-555-0123-99", "chars", nil, false},
		{"Phone", "too short", "123", "chars", nil, false},
		{"Phone", "empty", "", "chars", nil, false},
		{"Phone", "letter after min length", "+56-abc-123-456", "not allowed", nil, false},

		// ── Radio ────────────────────────────────────────────────────────────
		{"Radio", "valid key m", "m", "", nil, false},
		{"Radio", "valid key f", "f", "", nil, false},
		{"Radio", "empty", "", "chars", nil, false},

		// ── Rut ──────────────────────────────────────────────────────────────
		{"Rut", "valid 7863697-1", "7863697-1", "", nil, false},
		{"Rut", "valid K uppercase", "20373221-K", "", nil, false},
		{"Rut", "valid k lowercase", "20373221-k", "", nil, false},
		{"Rut", "no hyphen alpha", "15890022k", "hyphen", nil, false},
		{"Rut", "no hyphen digits", "177344788", "hyphen", nil, false},
		{"Rut", "wrong check digit", "7863697-2", "invalid", nil, false},
		{"Rut", "leading zero", "01234567-1", "invalid", nil, false},
		{"Rut", "empty", "", "chars", nil, false},

		// ── Search ───────────────────────────────────────────────────────────
		{"Search", "valid search", "golang tinywasm", "", nil, false},
		{"Search", "empty search", "", "", nil, false},

		// ── Select ───────────────────────────────────────────────────────────
		{"Select", "valid", "admin", "", nil, false},
		{"Select", "empty", "", "chars", nil, false},

		// ── Text ─────────────────────────────────────────────────────────────
		{"Text", "valid with tilde", "Juan Pérez", "", nil, false},
		{"Text", "valid with dots", "Dr. Smith (Jr.)", "", nil, false},
		{"Text", "too short", "A", "chars", nil, false},
		{"Text", "empty", "", "chars", nil, false},
		{"Text", "at sign", "user@domain", "not allowed", nil, false},

		// ── Textarea ─────────────────────────────────────────────────────────
		{"Textarea", "valid long", "IRRITACION EN PIEL DE ROSTRO. ALERGIAS NO.", "", nil, false},
		{"Textarea", "valid with newline", "Line one.\nLine two.", "", nil, false},
		{"Textarea", "empty", "", "chars", nil, false},
		{"Textarea", "too short", "Hi", "chars", nil, false},
	}

	for _, c := range cases {
		c := c
		t.Run(c.t+"/"+c.name, func(t *testing.T) {
			inp := buildInput(t, c.t, c.opts)
			if c.req {
				if s, ok := inp.(interface{ SetRequired(bool) }); ok {
					s.SetRequired(true)
				}
			}
			checkErr(t, inp.Validate(c.val), c.err)
		})
	}
}
