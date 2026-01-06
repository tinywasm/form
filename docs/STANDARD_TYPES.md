# Standard Input Types & Validation

`tinywasm/form` includes "Contextual Intelligence". It automatically maps field names and data types to specific HTML5 input types with predefined validation rules.

This mapping allows you to define a `Name` field in your struct, and the library automatically enforces:
- HTML: `<input type="text">`
- Validation: Allow letters/spaces, Min 3 chars, Max 50 chars.

## Default Registry

The following types are **pre-registered** internally using `form.RegisterType`. They serve as the baseline "smart" behavior.

| Field Name Pattern | HTML Type | Validation / Constraints |
| :--- | :--- | :--- |
| `Name`, `Text`, `*Name` | `text` | **Permitted**: Letters, Spaces, `.` `,` `()`<br>**Min**: 2, **Max**: 100 |
| `Email`, `Mail` | `email` | **Permitted**: Letters, Numbers, `@._`<br>**Extra**: Email Structure Check |
| `Password`, `Pass` | `password` | **Permitted**: All generic chars<br>**Min**: 5, **Max**: 50 |
| `Phone`, `Mobile` | `tel` | **Permitted**: Numbers only<br>**Min**: 7, **Max**: 11 |
| `Age`, `Year` | `number` | **Permitted**: Numbers only |
| `Date`, `BirthDate` | `date` | Format `YYYY-MM-DD` |
| `RUT` (ChileID) | `text` | **Permitted**: Numbers, `k`, `K`, `-`<br>**Extra**: Modulo 11 Algorithm |
| `IP` | `text` | **Permitted**: Numbers, `.` `:`<br>**Extra**: IPv4/IPv6 Logic |

## Code Transfer Plan (Reference -> New)
The following logic from `Archive/mono` maps directly here:
*   `Email`: Uses logic from `field.go:98`
*   `RUT`: Uses logic from `field.go:193` & `validation.go`
*   `IP`: Uses logic from `field.go:140`

## Extensibility (`RegisterType`)

You can register new types using Go closures for full power (no JSON/Serialization limitations).

```go
form.RegisterType("MyCustomType", form.TypeConfig{
    InputHTML: "text",
    Permitted: form.Permitted{
        Letters: true,
        ExtraValidation: func(val string) error {
            // My custom logic
            return nil
        },
    },
})
```

## Validation Logic (`Permitted`)

Internal validation uses a "Permissible" whitelist approach to ensure security and data integrity.

### `Permitted` Structure
Every Type has a `Permitted` configuration:
*   **Letters**: `bool` (A-Z, a-z)
*   **Numbers**: `bool` (0-9)
*   **Characters**: `[]rune` (Allowed special chars e.g., `@`, `.`, `-`)
*   **Min/Max**: `int` (Length constraints)
*   **StartWith**: `*Permitted` (Rules for the first character)
*   **ExtraValidation**: `func(val string) error` (Custom logic)
