# Struct Tags

Tags are now parsed by `ormc` at code-generation time to populate the `fmt.Field` schema. `tinywasm/form` reads these values from the `fmt.Fielder` interface at runtime.

## Supported Tags (via ormc)

| Tag | Format | Description |
|-----|--------|-------------|
| `form` | `"type"` or `"-"` | Overrides input type or excludes field from form. |
| `options` | `"key1:text1,key2:text2"` | Sets options for select/radio/datalist. |
| `placeholder` | `"string"` | Overrides auto-translated placeholder. |
| `title` | `"string"` | Overrides auto-translated title. |
| `validate` | `"false"` | Skips validation for this field. |

## `form:` Tag Reference

- **Type Override**: `form:"email"` ensures the email input is used.
- **Exclusion**: `form:"-"` skips the field entirely when building the form.

## Runtime Interaction

While tags are parsed during generation, the form still supports programmatic overrides:

```go
f.Input("Field").SetPlaceholder("New Hint")
f.SetOptions("Category", opts...)
```
