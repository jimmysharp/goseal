package app

import "example.com/testproject/domain"

// SHOULD NOT REPORT: Using constructor function
func WithConstructor() {
	user, _ := domain.NewUser(123, "Alice", 30)
	_ = user
}

// SHOULD REPORT: Direct initialization from different package
func WithoutConstructor() {
	user := domain.User{ // want "direct construction of struct User is prohibited, use constructor function"
		ID:   123,
		Name: "Bob",
		Age:  25,
	}
	_ = user
}

// SHOULD REPORT: Direct assignment from different package
func DirectAssignment() {
	user, _ := domain.NewUser(123, "Charlie", 35)
	user.Name = "Dave" // want "direct assignment to field Name of struct User is prohibited, use constructor function"
}
