package domain

import "math/rand/v2"

// SHOULD NOT REPORT: Using constructor function
func CreateUser(name string, age int) (*User, error) {
	return NewUser(rand.Int(), name, age)
}

// SHOULD REPORT: Direct initialization in same package (allow-same-package defaults to false)
func CreateUserDirect(name string, age int) *User {
	return &User{ // want "direct construction of struct User is prohibited, use constructor function"
		ID:   rand.Int(),
		Name: name,
		Age:  age,
	}
}
