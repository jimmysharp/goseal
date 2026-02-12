package app

import "example.com/testproject/domain"

// WithFactoryFunction uses a factory function (should NOT be reported).
func WithFactoryFunction() {
	user, _ := domain.NewUser(123, "Alice", 30)
	_ = user
}

// WithoutFactoryFunction uses direct struct literal (should be reported).
func WithoutFactoryFunction() {
	_ = domain.User{
		ID:   123,
		Name: "Bob",
		Age:  25,
	}
}

// AssignmentInReceiver uses a receiver method (should NOT be reported).
func AssignmentInReceiver() {
	user, _ := domain.NewUser(123, "Charlie", 35)
	user.UpdateName("Dave")
}

// DirectAssignment assigns to a field directly (should be reported).
func DirectAssignment() {
	user, _ := domain.NewUser(123, "Charlie", 35)
	user.Name = "Dave"
}
