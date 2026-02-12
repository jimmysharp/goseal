package app

import (
	"example.com/testproject/domain/user"
)

// SHOULD NOT REPORT: Using factory function
func WithFactoryFunction() {
	u, _ := user.NewUser(123, "Alice", 30)
	_ = u
}

// SHOULD REPORT: Direct initialization outside target packages (init-scope: in-target-packages)
func DirectInitFromApp() {
	_ = user.User{ // want "direct construction of sealed struct User is not allowed from outside target packages \\(init-scope: in-target-packages\\)"
		ID:   123,
		Name: "Bob",
		Age:  25,
	}

	_ = &user.User{ // want "direct construction of sealed struct User is not allowed from outside target packages \\(init-scope: in-target-packages\\)"
		ID:   123,
		Name: "Bob",
		Age:  25,
	}
}

// SHOULD REPORT: Direct field assignment
func DirectAssignment() {
	u, _ := user.NewUser(123, "Charlie", 35)

	u.ID = 456      // want "direct assignment to field ID of sealed struct User is not allowed from outside target packages \\(mutation-scope: in-target-packages\\)"
	u.Name = "Dave" // want "direct assignment to field Name of sealed struct User is not allowed from outside target packages \\(mutation-scope: in-target-packages\\)"
	u.Age = 40      // want "direct assignment to field Age of sealed struct User is not allowed from outside target packages \\(mutation-scope: in-target-packages\\)"
}
