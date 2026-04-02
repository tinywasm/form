# Struct Tags

Tags are parsed by `ormc` (the ORM code generator) at compile-time to populate the `fmt.Field` schema. `tinywasm/form` reads these values from the `fmt.Fielder` interface at runtime.

The unified tag for all form-related metadata is `input:`.

## The `input:` Tag

The `input:` tag combines widget type, validation rules, and other metadata into a single comma-separated string.

### 1. Widget Types
Overrides the default widget (which is `text` for strings, `number` for ints).

| Tag | Resulting Widget |
|-----|------------------|
| `input:"email"` | `input.Email()` |
| `input:"password"` | `input.Password()` |
| `input:"textarea"` | `input.Textarea()` |
| `input:"checkbox"` | `input.Checkbox()` |
| `input:"select"` | `input.Select()` |
| `input:"radio"` | `input.Radio()` |
| `input:"datalist"` | `input.Datalist()` |
| `input:"phone"` | `input.Phone()` |
| `input:"date"` | `input.Date()` |
| `input:"-"` | Skips the field in form rendering |

### 2. Validation Rules
Mapped directly to `fmt.Permitted` and `NotNull` properties in the schema.

| Rule | Effect |
|------|--------|
| `required` | Sets `NotNull: true` |
| `min=5` | Sets `Minimum: 5` (length or value) |
| `max=100` | Sets `Maximum: 100` |
| `letters` | Requires letters only |
| `numbers` | Requires numbers only |
| `spaces` | Allows spaces |

### 3. Metadata
| Prefix | Description |
|--------|-------------|
| `title="My Label"` | Overrides the auto-translated label |
| `placeholder="Search..."` | Overrides the auto-translated placeholder |
| `options="A:Alpha,B:Beta"` | Sets static options for select/radio |

## Example

```go
type Contact struct {
    Name     string `input:"required,min=2,max=50,title=\"Full Name\""`
    Email    string `input:"email,required"`
    Bio      string `input:"textarea,max=500,placeholder=\"Tell us about yourself\""`
    Category string `input:"select,options=\"1:Worker,2:Manager\",required"`
}
```

## Runtime Interaction

While tags are parsed during generation, the form still supports programmatic overrides at runtime:

```go
f.Input("Name").SetTitle("Legal Name")
f.SetOptions("Category", dynamicOptions...)
```
