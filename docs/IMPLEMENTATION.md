# Implementation Notes

> See [README.md](../README.md) for the full API reference.
> See [DESIGN.md](DESIGN.md) for architecture.

## File Map

| File | Responsibility |
|------|---------------|
| `form.go` | `Form` struct, `New()`, `Input()`, `SetOptions()`, `SetValues()`, `Namer` |
| `sync.go` | `SyncValues()`, pointer-based field sync |
| `registry.go` | `SetGlobalClass()` |
| `render.go` | `Form.RenderHTML()`, `Form.SetSSR()` |
| `validate.go` | `Form.Validate()` |
| `validate_struct.go` | `Form.ValidateData()` (crudp.DataValidator) |
| `words.go` | Registers form UI words into fmt dictionary |
| `mount.go` | `Form.OnMount()`, `Form.OnUnmount()` (wasm build tag) |
| `input/interface.go` | `Input` interface (embeds `fmt.Widget` + `dom.Component`) |
| `input/base.go` | `Base` struct embedded by all inputs |
| `input/*.go` | 17 concrete input implementations |

## Adding a New Input

1. Create `input/mytype.go` — embed `Base`, configure `fmt.Permitted` rules, implement rendering.
2. Add `NewMyType() fmt.Widget` constructor (template, no position).
3. Add `Clone(parentID, name string) fmt.Widget` method on the concrete type.
4. `Type()` and `Validate()` are inherited from `Base` — no need to implement.
5. Add test cases in `input/inputs_test.go`.
6. Use `input.NewMyType()` in `ormc` schema generation to assign the widget to fields.

## Key Constraints

- Only `github.com/tinywasm/fmt` — no `errors`, `strconv`, `strings`
- No maps in WASM-facing code (increases binary size) — use slices
- No `reflect` at runtime
