package app

import "example.com/testproject/domain"

// NOT SUPPORTED: Zero value declarations are not supported to avoid false positives (see WithZeroValueAndFactoryFunction)
func WithZeroValue() {
	var _ domain.User
	var _ *domain.User
}

// NOT SUPPORTED: Valid pattern - zero value with conditional factory function
// (would be false positive if zero value detection was enabled)
func WithZeroValueAndFactoryFunction(isMale bool) {
	var user *domain.User
	if isMale {
		user, _ = domain.NewUser(1, "Alice", 30)
	} else {
		user, _ = domain.NewUser(2, "Bob", 25)
	}
	_ = user
}

type alias = domain.User

// NOT SUPPORTED: Type alias initialization is not reported
// Type aliases have valid use cases (e.g., to expose internal types for testing)
func WithTypeAlias() {
	_ = alias{
		ID:   123,
		Name: "Eve",
		Age:  28,
	}
}
