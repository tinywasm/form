package input





// tildeSetter is the private contract for widgets that allow toggling accented chars.
// Unexported on purpose — it never appears in user-facing APIs.
type tildeSetter interface{ setTilde(bool) }

// SetTilde toggles tilde acceptance on w and returns w for chaining.
// No-op if w doesn't implement tildeSetter (e.g. checkbox, select).
// Called from ormc-generated code when struct tag `notilde` is present:
//   Widget: input.SetTilde(input.Text(), false)
func SetTilde(w Input, v bool) Input {
	if t, ok := w.(tildeSetter); ok {
		t.setTilde(v)
	}
	return w
}
