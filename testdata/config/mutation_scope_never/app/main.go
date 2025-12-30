package app

import "example.com/testproject/domain"

// SHOULD NOT REPORT: Using factory function
func WithFactoryFunction() {
	user, _ := domain.NewUser(123, "Alice", 30)
	_ = user
}

// SHOULD REPORT: Direct initialization without factory function (init-scope: same-package)
func WithoutFactoryFunction() {
	_ = domain.User{ // want "direct construction of struct User is prohibited outside allowed scope"
		ID:   123,
		Name: "Bob",
		Age:  25,
	}
}

// SHOULD REPORT: mutation-scope is "never", so all field assignments are prohibited
func DirectAssignment() {
	user, _ := domain.NewUser(123, "Charlie", 35)

	user.ID = 456      // want "direct assignment to field ID of struct User is prohibited outside allowed scope"
	user.Name = "Dave" // want "direct assignment to field Name of struct User is prohibited outside allowed scope"
	user.Age = 40      // want "direct assignment to field Age of struct User is prohibited outside allowed scope"
}
