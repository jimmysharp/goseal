package domain

import "math/rand/v2"

// SHOULD NOT REPORT: Initialization in same package is allowed (init-scope: same-package)
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

// SHOULD NOT REPORT: Assignment same packages is allowed (mutation-scope: same-package)
func UpdateUserWithoutReceiver(u *User, name string, age int) *User {
	u.Name = name
	u.Age = age

	return u
}
