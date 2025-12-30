package app

import (
	"example.com/testproject/domain/user"
)

// SHOULD NOT REPORT: Using constructor function
func WithConstructor() {
	u, _ := user.NewUser(123, "Alice", 30)
	_ = u
}

// SHOULD REPORT: Direct initialization outside target packages (init-scope: in-target-packages)
func DirectInitFromApp() {
	_ = user.User{ // want "direct construction of struct User is prohibited outside allowed scope"
		ID:   123,
		Name: "Bob",
		Age:  25,
	}

	_ = &user.User{ // want "direct construction of struct User is prohibited outside allowed scope"
		ID:   123,
		Name: "Bob",
		Age:  25,
	}
}

// SHOULD REPORT: Direct field assignment
func DirectAssignment() {
	u, _ := user.NewUser(123, "Charlie", 35)

	u.ID = 456      // want "direct assignment to field ID of struct User is prohibited outside allowed scope"
	u.Name = "Dave" // want "direct assignment to field Name of struct User is prohibited outside allowed scope"
	u.Age = 40      // want "direct assignment to field Age of struct User is prohibited outside allowed scope"
}
