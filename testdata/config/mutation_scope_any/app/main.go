package app

import "example.com/testproject/domain"

// SHOULD NOT REPORT: Using factory function
func WithFactoryFunction() {
	user, _ := domain.NewUser(123, "Alice", 30)
	_ = user
}

// SHOULD REPORT: Direct initialization without factory (init-scope: same-package)
func WithoutFactoryFunction() {
	_ = domain.User{ // want "direct construction of sealed struct User is not allowed from outside its package \\(init-scope: same-package\\)"
		ID:   123,
		Name: "Bob",
		Age:  25,
	}
}

// SHOULD NOT REPORT: Mutation is always allowed (mutation-scope: any)
func DirectAssignment() {
	user, _ := domain.NewUser(123, "Charlie", 35)

	user.ID = 456
	user.Name = "Dave"
	user.Age = 40
}

// SHOULD NOT REPORT: Mutation is always allowed (mutation-scope: any)
func AssignmentThroughFunction() {
	user, _ := domain.NewUser(123, "Eve", 28)

	domain.UpdateUserAge(user, 29)
}
