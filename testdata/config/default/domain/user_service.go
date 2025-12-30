package domain

import "math/rand/v2"

// SHOULD NOT REPORT: Using factory function
func CreateUser(name string, age int) (*User, error) {
	return NewUser(rand.Int(), name, age)
}

// SHOULD NOT REPORT: Initialization in same package is allowed (init-scope: same-package)
func CreateUserDirect(name string, age int) *User {
	return &User{
		ID:   rand.Int(),
		Name: name,
		Age:  age,
	}
}

// SHOULD REPORT: Direct Assignment in non-receiver function (mutation-scope: receiver)
func UpdateUserWithoutReceiver(u *User, name string, age int) *User {
	u.Name = name // want "direct assignment to field Name of struct User is prohibited outside allowed scope"
	u.Age = age   // want "direct assignment to field Age of struct User is prohibited outside allowed scope"

	return u
}
