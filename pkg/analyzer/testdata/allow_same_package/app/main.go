package app

import "example.com/testproject/domain"

func WithConstructor() {
	user := domain.NewUser(123, "Alice", 30)
	_ = user
}

func WithoutConstructor() {
	user := domain.User{ // want "direct construction of struct User is prohibited, use constructor function"
		ID:   123,
		Name: "Bob",
		Age:  25,
	}
	_ = user
}

func DirectAssignment() {
	user := domain.NewUser(123, "Charlie", 35)
	user.Name = "Dave" // want "direct assignment to field Name of struct User is prohibited, use constructor function"
}
