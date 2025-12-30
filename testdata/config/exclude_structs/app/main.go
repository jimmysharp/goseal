package app

import "example.com/testproject/domain"

// SHOULD NOT REPORT: Using constructor function
func WithConstructor() {
	user, _ := domain.NewUser(123, "Alice", 30)
	_ = user
}

// SHOULD REPORT: Direct initialization of User without constructor (init-scope: same-package)
func WithoutConstructor() {
	_ = domain.User{ // want "direct construction of struct User is prohibited outside allowed scope"
		ID:   123,
		Name: "Bob",
		Age:  25,
	}
}

// SHOULD NOT REPORT: UserDTO is excluded from protection (matches ".*DTO$")
func WithDTO() {
	_ = domain.UserDTO{
		ID:   1,
		Name: "Alice",
		Age:  30,
	}

	dto := &domain.UserDTO{
		ID:   2,
		Name: "Bob",
		Age:  25,
	}
	dto.ID = 999
	dto.Name = "Updated"
}

// SHOULD REPORT: Direct field assignment to User (mutation-scope: receiver)
func DirectAssignment() {
	user, _ := domain.NewUser(123, "Charlie", 35)

	user.ID = 456      // want "direct assignment to field ID of struct User is prohibited outside allowed scope"
	user.Name = "Dave" // want "direct assignment to field Name of struct User is prohibited outside allowed scope"
	user.Age = 40      // want "direct assignment to field Age of struct User is prohibited outside allowed scope"
}
