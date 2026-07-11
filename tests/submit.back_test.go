//go:build !wasm

package form_test

import "testing"

func TestSubmit_Back(t *testing.T) {
	runSubmitTests(t)
}
