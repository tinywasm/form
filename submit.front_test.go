//go:build wasm

package form

import "testing"

func TestSubmit_Front(t *testing.T) {
	runSubmitTests(t)
}
