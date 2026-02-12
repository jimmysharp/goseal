package domain

import "math/rand/v2"

// SHOULD NOT REPORT: Using factory function
func CreateUser(name string, age int) (*User, error) {
	return NewUser(rand.Int(), name, age)
}

// SHOULD REPORT: Direct initialization in non-factory function (factory-names)
func CreateUserDirect(name string, age int) *User {
	return &User{ // want "direct construction of sealed struct User is not allowed outside factory functions \\(factory-names\\)"
		ID:   rand.Int(),
		Name: name,
		Age:  age,
	}
}

// SHOULD REPORT: Direct Assignment in non-receiver function (mutation-scope: receiver)
func UpdateUserWithoutReceiver(u *User, name string, age int) {
	u.Name = name // want "direct assignment to field Name of sealed struct User is not allowed outside its receiver methods \\(mutation-scope: receiver\\)"
	u.Age = age   // want "direct assignment to field Age of sealed struct User is not allowed outside its receiver methods \\(mutation-scope: receiver\\)"
}
