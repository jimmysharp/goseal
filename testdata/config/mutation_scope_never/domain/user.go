package domain

import "fmt"

type User struct {
	ID   int
	Name string
	Age  int
}

func NewUser(id int, name string, age int) (*User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("id must be positive: %d", id)
	}
	if name == "" {
		return nil, fmt.Errorf("name must not be empty")
	}
	if age < 0 {
		return nil, fmt.Errorf("age must be non-negative: %d", age)
	}

	return &User{
		ID:   id,
		Name: name,
		Age:  age,
	}, nil
}

// SHOULD REPORT: Mutation is always prohibited (mutation-scope: never)
func (u *User) UpdateName(name string) error {
	if name == "" {
		return fmt.Errorf("name must not be empty")
	}
	u.Name = name // want "direct assignment to field Name of sealed struct User is not allowed anywhere \\(mutation-scope: never\\)"
	return nil
}

// SHOULD REPORT: Mutation is always prohibited (mutation-scope: never)
func UpdateUserAge(u *User, age int) {
	u.Age = age // want "direct assignment to field Age of sealed struct User is not allowed anywhere \\(mutation-scope: never\\)"
}
