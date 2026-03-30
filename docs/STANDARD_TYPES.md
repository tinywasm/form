# Standard Input Types

17 built-in types in `input/`. Each provides `NewXxx() fmt.Widget` and `Clone(parentID, name) fmt.Widget`.

| Input | HTML type | Constructor | Key Aliases |
|-------|-----------|-------------|-------------|
| `Address` | `text` | `input.NewAddress()` | `address`, `addr` |
| `Checkbox` | `checkbox` | `input.NewCheckbox()` | `check`, `boolean`, `bool` |
| `Date` | `date` | `input.NewDate()` | `fecha` |
| `Datalist` | `text` | `input.NewDatalist()` | `list`, `options` |
| `Email` | `email` | `input.NewEmail()` | `mail`, `correo` |
| `Filepath` | `text` | `input.NewFilepath()` | `path`, `dir`, `file` |
| `Gender` | `radio` | `input.NewGender()` | `gender`, `sex` |
| `Hour` | `time` | `input.NewHour()` | `hour` |
| `IP` | `text` | `input.NewIP()` | `ip` |
| `Number` | `number` | `input.NewNumber()` | `num`, `amount`, `price`, `age` |
| `Password` | `password` | `input.NewPassword()` | `pass`, `clave`, `pwd` |
| `Phone` | `tel` | `input.NewPhone()` | `phone`, `mobile`, `cell` |
| `Radio` | `radio` | `input.NewRadio()` | -- |
| `Rut` | `text` | `input.NewRut()` | `rut`, `run`, `dni` |
| `Select` | `select` | `input.NewSelect()` | -- |
| `Text` | `text` | `input.NewText()` | `name`, `fullname`, `username` |
| `Textarea` | `textarea` | `input.NewTextarea()` | `description`, `details`, `comments` |

See [input/README.md](../input/README.md) for detailed validation rules and rendering behavior per type.
