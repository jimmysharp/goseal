# conseal

`conseal` is a linter that checks direct instantiation of Go structs and enforces creation through constructor functions. Use it to prevent invalid values from leaking into your objects, such as when applying the Value Object pattern.

## Installation

```bash
go install github.com/jimmysharp/conseal/cmd/conseal@latest
```

## Configuration

Create `.conseal.yml` in your project root:

```yaml
# List of rgexps for packages containing target structs
packages:
  - "github\\.com/yourorg/domain/.*"
  - "github\\.com/yourorg/model/.*"

# List of regexps for functions considered as constructors
constructors:
  - "^New.*"

# List of regexps for files to ignore
ignore-files:
  - "_test\\.go$"
  - "mock_.*\\.go$"

# Whether to allow direct struct construction within the same package
allow-same-package: false
```

## Usage

```bash
conseal ./...
```

## Examples

### Example domain object

```go
package domain

type User struct {
    ID   string
    Name string
}

func NewUser(id, name string) *User {
    // ✅ Allowed inside a constructor
    return &User{
        ID:   id,
        Name: name,
    }
}
```

### Example that gets flagged (NG)

```go
package app

import "github.com/yourorg/domain"

func CreateUser() {
    // ❌ Direct use of a struct literal
    user := domain.User{
        ID:   "123",
        Name: "Alice",
    }
    
    // ❌ Direct field assignment
    user.Name = "Bob"
}
```

### Recommended usage (OK)

```go
package app

import "github.com/yourorg/domain"

func CreateUser() {
    // ✅ Use the constructor
    user := domain.NewUser("123", "Alice")
    _ = user
}
```

## License

MIT License – see [LICENSE](LICENSE) for details.
