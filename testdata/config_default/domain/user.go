package domain

import "fmt"

type User struct {
	ID   int
	Name string
	Age  int
}

// SHOULD NOT REPORT: "^New.*" is default constructor pattern (struct-packages)
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
