# goseal

`goseal` is a linter that checks direct initialization of Go structs and enforces creation through factory functions. It "seals" your structs, ensuring they can only be instantiated through approved factory functions. Use it to prevent invalid values from leaking into your objects, such as when applying the Value Object pattern.

## Installation

```bash
go install github.com/jimmysharp/goseal/cmd/goseal@latest
```

## Configuration

Create `.goseal.yml` in your project root:

**Note:** While `target-packages` is optional (omitting it will target all packages), it is **strongly recommended** to explicitly specify the packages containing your domain structs. This prevents false positives with structs in third-party or standard library code.

```yaml
# List of regexps for packages containing target structs to protect
# If not specified or empty, all packages are targeted
# Default: []
target-packages:
  - "github\\.com/yourorg/domain/.*"
  - "github\\.com/yourorg/model/.*"

# List of regexps for struct names to exclude from protection
# Even if in target-packages, these structs won't be protected
# Default: []
exclude-structs:
  - "^Config$"
  - "DTO$"

# List of regexps for functions considered as factory functions
# If empty, struct initialization is allowed in any context
# Default: []
factory-names:
  - "^New.*"

# Scope for struct initialization
# - any: Allow initialization from all packages
# - in-target-packages: Allow initialization from packages in target-packages
# - same-package: Allow initialization only within the same package
# Default: same-package
init-scope: same-package

# Scope for field mutation
# - any: Allow field mutation everywhere
# - in-target-packages: Allow field mutation from packages in target-packages
# - receiver: Allow field mutation only in receiver methods
# - same-package: Allow field mutation within the same package
# - never: Never allow field mutation
# Default: receiver
mutation-scope: receiver

# List of regexps for files to ignore
# Default: []
ignore-files:
  - "_test\\.go$"
  - "mock_.*\\.go$"
```

**Note:** Auto-generated files are automatically skipped.

## Usage

### Standalone

```bash
goseal ./...
```

### golangci-lint (custom plugin)

goseal can also be used as a [golangci-lint custom plugin](https://golangci-lint.run/plugins/module-plugins/). When used as a plugin, `.goseal.yml` is not used. Instead, configure settings directly in `.golangci.yml`.

Create `.custom-gcl.yml` to build a custom golangci-lint binary with goseal:

```yaml
version: v2.9.0
plugins:
  - module: 'github.com/jimmysharp/goseal'
    import: 'github.com/jimmysharp/goseal/plugin'
```

Then configure goseal in `.golangci.yml` (the same settings as `.goseal.yml` are available under `settings`):

```yaml
version: "2"
linters:
  enable:
    - goseal
  settings:
    custom:
      goseal:
        type: module
        settings:
          target-packages:
            - "github\\.com/yourorg/domain/.*"
          factory-names:
            - "^New.*"
          init-scope: same-package
          mutation-scope: receiver
```

Build and run:

```bash
golangci-lint custom
./custom-gcl run ./...
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
    // ✅ Allowed inside a factory function
    return &User{
        ID:   id,
        Name: name,
    }
}

// ✅ Allowed: mutation in receiver method (when mutation-scope: receiver)
func (u *User) UpdateName(name string) {
    u.Name = name
}
```

### Example that gets flagged (NG)

```go
package app

import "github.com/yourorg/domain"

func CreateUser() {
    // ❌ Direct use of a struct literal (when init-scope: same-package)
    user := domain.User{
        ID:   "123",
        Name: "Alice",
    }
    
    // ❌ Direct field assignment (when mutation-scope: receiver)
    user.Name = "Bob"
}
```

### Recommended usage (OK)

```go
package app

import "github.com/yourorg/domain"

func CreateUser() {
    // ✅ Use the factory function
    user := domain.NewUser("123", "Alice")
    
    // ✅ Use the method for mutation
    user.UpdateName("Bob")
}
```

## Usecases

### Strict constructor usage with immutability

Force constructor usage in domain packages and prohibit all mutations:

```yaml
target-packages:
  - "github\\.com/yourorg/domain/.*"
factory-names:
  - "^New.*"
init-scope: same-package
mutation-scope: never
```

### Entity Pattern

Allow mutation through methods while preventing direct field access:

```yaml
target-packages:
  - "github\\.com/yourorg/domain/entity/.*"
factory-names:
  - "^New.*"
init-scope: same-package
mutation-scope: receiver
```

### Domain packages with subpackages

Allow cross-package initialization within domain packages (useful when domain is split into subpackages like `domain/user`, `domain/order`, etc.):

```yaml
target-packages:
  - "github\\.com/yourorg/domain/.*"
factory-names:
  - "^New.*"
init-scope: in-target-packages
mutation-scope: never
```

### Trust package boundaries

Allow free operations within domain packages:

```yaml
target-packages:
  - "github\\.com/yourorg/domain/.*"
init-scope: same-package
mutation-scope: same-package
```

### Exclude specific structs

Exclude specific structs (e.g., configuration structs, DTOs) from protection:

```yaml
target-packages:
  - "github\\.com/yourorg/domain/.*"
exclude-structs:
  - "^Config$"
  - "DTO$"
factory-names:
  - "^New.*"
init-scope: same-package
mutation-scope: receiver
```

## License

MIT License – see [LICENSE](LICENSE) for details.
