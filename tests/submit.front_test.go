//go:build wasm

package form_test

import (
	"testing"
)

func TestSubmit_Front(t *testing.T) {
	runSubmitTests(t)
}
