//go:build !wasm

package form

import "testing"

func TestSubmit_Back(t *testing.T) {
	runSubmitTests(t)
}
