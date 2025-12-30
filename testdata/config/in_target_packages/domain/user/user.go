package user

import "fmt"

type User struct {
	ID   int
	Name string
	Age  int
}

// SHOULD NOT REPORT: Function matching "^New.*" is considered a factory (factory-names)
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

// SHOULD NOT REPORT: Assignment in target packages is allowed (mutation-scope: in-target-packages)
func (u *User) UpdateName(name string) error {
	if name == "" {
		return fmt.Errorf("name must not be empty")
	}
	u.Name = name
	return nil
}
