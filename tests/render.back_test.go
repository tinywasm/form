//go:build !wasm

package form_test

import "testing"

func TestRender_Back(t *testing.T) {
	runRenderTests(t)
}
