# Implementation Notes

> See [README.md](../README.md) for the full API reference.
> See [DESIGN.md](DESIGN.md) for architecture.
> See [input/README.md](../input/README.md) for all input types.

## File Map

| File | Responsibility |
|------|---------------|
| `form.go` | `Form` struct, `New()`, tag parsing, `SyncValues()`, `Input()`, `SetOptions()`, `SetValues()` |
| `registry.go` | `registeredInputs`, `forms`, `RegisterInput()`, `findInputForField()`, `SetGlobalClass()` |
| `render.go` | `Form.RenderHTML()`, `Form.SetSSR()` |
| `validate.go` | `Form.Validate()` |
| `mount.go` | `Form.OnMount()`, `Form.OnUnmount()` (wasm build tag) |
| `tags.go` | `ParseOptionsTag()`, `GetTagOptions()` helpers |
| `input/base.go` | `Base` struct embedded by all inputs |
| `input/permitted.go` | `Permitted` whitelist validation engine |
| `input/interface.go` | `Input` interface |
| `input/*.go` | Individual input implementations |

## Adding a New Input

1. Create `input/mytype.go` — embed `Base`, add `Permitted`, implement `Input` interface
2. Register in `registry.go` `init()`: `input.MyType("", "")`
3. Add to `input/README.md` table
4. Add test cases in `input/inputs_test.go`

## Key Constraints

- Only `github.com/tinywasm/fmt` — no `errors`, `strconv`, `strings`
- No maps in WASM-facing code (increases binary size) — use slices
- `Permitted.ExtraValidation` for complex rules that can't use whitelist
