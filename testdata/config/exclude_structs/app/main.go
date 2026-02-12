package app

import "example.com/testproject/domain"

// SHOULD NOT REPORT: Using factory function
func WithFactoryFunction() {
	user, _ := domain.NewUser(123, "Alice", 30)
	_ = user
}

// SHOULD REPORT: Direct initialization of User without factory function (init-scope: same-package)
func WithoutFactoryFunction() {
	_ = domain.User{ // want "direct construction of sealed struct User is not allowed from outside its package \\(init-scope: same-package\\)"
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

	user.ID = 456      // want "direct assignment to field ID of sealed struct User is not allowed outside its receiver methods \\(mutation-scope: receiver\\)"
	user.Name = "Dave" // want "direct assignment to field Name of sealed struct User is not allowed outside its receiver methods \\(mutation-scope: receiver\\)"
	user.Age = 40      // want "direct assignment to field Age of sealed struct User is not allowed outside its receiver methods \\(mutation-scope: receiver\\)"
}
