package app

import "example.com/testproject/domain"

// SHOULD NOT REPORT: Using factory function
func WithFactoryFunction() {
	user, _ := domain.NewUser(123, "Alice", 30)
	_ = user
}

// SHOULD REPORT: Direct initialization without factory function
func WithoutFactoryFunction() {
	user := domain.User{ // want "direct construction of struct User is prohibited outside allowed scope"
		ID:   123,
		Name: "Bob",
		Age:  25,
	}
	_ = user
}

// SHOULD REPORT: Direct initialization with pointer
func WithoutFactoryFunctionByPointer() {
	user := &domain.User{ // want "direct construction of struct User is prohibited outside allowed scope"
		ID:   123,
		Name: "Bob",
		Age:  25,
	}
	_ = user
}

// SHOULD NOT REPORT: Assignment through method
func AssignmentInReceiver() {
	user, _ := domain.NewUser(123, "Charlie", 35)

	user.UpdateName("Dave")
}

// SHOULD REPORT: Direct field assignment
func DirectAssignment() {
	user, _ := domain.NewUser(123, "Charlie", 35)

	user.Name = "Dave" // want "direct assignment to field Name of struct User is prohibited outside allowed scope"
}
