package domain

import "math/rand/v2"

// SHOULD NOT REPORT: Using constructor function
func CreateUser(name string, age int) (*User, error) {
	return NewUser(rand.Int(), name, age)
}

// SHOULD NOT REPORT: Same package initialization is allowed (allow-same-package: true)
func CreateUserDirect(name string, age int) *User {
	return &User{
		ID:   rand.Int(),
		Name: name,
		Age:  age,
	}
}
