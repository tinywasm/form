# Standard Input Types

19 built-in types in `input/`. Each provides `Xxx() model.Kind` and `Clone(parentID, name) model.Kind`.

| Input | HTML type | Constructor | Key Aliases |
|-------|-----------|-------------|-------------|
| `Address` | `text` | `input.Address()` | `address`, `addr` |
| `Checkbox` | `checkbox` | `input.Checkbox()` | `check`, `boolean`, `bool` |
| `Date` | `date` | `input.Date()` | `fecha` |
| `Datalist` | `text` | `input.Datalist()` | `list`, `options` |
| `Decimal` | `number` | `input.Decimal()` | `price`, `amount`, `percentage` |
| `Email` | `email` | `input.Email()` | `mail`, `correo` |
| `Filepath` | `text` | `input.Filepath()` | `path`, `dir`, `file` |
| `Gender` | `radio` | `input.Gender()` | `gender`, `sex` |
| `Hour` | `time` | `input.Hour()` | `hour` |
| `IP` | `text` | `input.IP()` | `ip` |
| `Number` | `number` | `input.Number()` | `num`, `age`, `count`, `quantity` |
| `Password` | `password` | `input.Password()` | `pass`, `clave`, `pwd` |
| `Phone` | `tel` | `input.Phone()` | `phone`, `mobile`, `cell` |
| `Radio` | `radio` | `input.Radio()` | -- |
| `Rut` | `text` | `input.Rut()` | `rut`, `run`, `dni` |
| `Search` | `search` | `input.Search()` | `search`, `find`, `q` |
| `Select` | `select` | `input.Select()` | -- |
| `Text` | `text` | `input.Text()` | `name`, `fullname`, `username` |
| `Textarea` | `textarea` | `input.Textarea()` | `description`, `details`, `comments` |

See [input/README.md](../input/README.md) for detailed validation rules and rendering behavior per type.
