package domain

import "math/rand/v2"

// SHOULD NOT REPORT: Using factory function
func CreateUser(name string, age int) (*User, error) {
	return NewUser(rand.Int(), name, age)
}

// SHOULD NOT REPORT: Function matching "^Create.*" is considered a factory (factory-names)
func CreateDefaultUser() *User {
	return &User{
		ID:   rand.Int(),
		Name: "Default",
		Age:  0,
	}
}

// SHOULD REPORT: Initialization in non-factory function (factory-names)
func BuildUser(name string, age int) *User {
	return &User{ // want "direct construction of struct User is prohibited, use allowed factory function"
		ID:   rand.Int(),
		Name: name,
		Age:  age,
	}
}

// SHOULD REPORT: Direct Assignment in non-receiver function (mutation-scope: receiver)
func UpdateUserWithoutReceiver(u *User, name string, age int) {
	u.Name = name // want "direct assignment to field Name of struct User is prohibited outside allowed scope"
	u.Age = age   // want "direct assignment to field Age of struct User is prohibited outside allowed scope"
}
