package domain

import "math/rand/v2"

type User struct {
	ID   int
	Name string
	Age  int
}

func NewUser(id int, name string, age int) *User {
	return &User{
		ID:   id,
		Name: name,
		Age:  age,
	}
}

func CreateUser(name string, age int) *User {
	return &User{
		ID:   rand.Int(),
		Name: name,
		Age:  age,
	}
}
