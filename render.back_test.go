//go:build !wasm

package form

import "testing"

func TestRender_Back(t *testing.T) {
	runRenderTests(t)
}
