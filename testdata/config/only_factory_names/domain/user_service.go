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
