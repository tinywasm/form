//go:build !wasm

package form

// IMPORTANT: this build tag is correct. The CSS is emitted during SSR and the browser
// receives it in the <style> of the initial HTML. WASM code does not need CSS.

import "github.com/tinywasm/css"

// RenderCSS returns the form's CSS contribution (additive — see tinywasm/css contract).
// Call from the project's ssr.go aggregate so assetmin picks it up.
func RenderCSS() *css.Stylesheet {
	return css.NewStylesheet(
		css.Rule(".tw-field-error",
			css.Display(css.Block),
			css.FontSize(css.TextSm),
			css.Color(css.ColorError),
			css.MinHeight(css.Em(1.2)),
		),
		css.Rule(".tw-field-error--visible",
			css.FontWeight(css.FontWeightMedium),
		),
	)
}
