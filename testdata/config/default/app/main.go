package app

import "example.com/testproject/domain"

// SHOULD NOT REPORT: Using constructor function
func WithConstructor() {
	user, _ := domain.NewUser(123, "Alice", 30)
	_ = user
}

// SHOULD REPORT: Direct initialization without constructor
func WithoutConstructor() {
	user := domain.User{ // want "direct construction of struct User is prohibited outside allowed scope"
		ID:   123,
		Name: "Bob",
		Age:  25,
	}
	_ = user
}

// SHOULD REPORT: Direct initialization with pointer
func WithoutConstructorByPointer() {
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
