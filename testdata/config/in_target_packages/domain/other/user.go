package other

import (
	"fmt"

	"example.com/testproject/domain/user"
)

// SHOULD NOT REPORT: Initialization in target packages is allowed (init-scope: in-target-packages)
func NewUser(id int, name string, age int) (*user.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("id must be positive: %d", id)
	}
	if name == "" {
		return nil, fmt.Errorf("name must not be empty")
	}
	if age < 0 {
		return nil, fmt.Errorf("age must be non-negative: %d", age)
	}

	return &user.User{
		ID:   id,
		Name: name,
		Age:  age,
	}, nil
}

// SHOULD NOT REPORT: Non-receiver mutation in target packages is allowed (mutation-scope: in-target-packages)
func Rename(u *user.User, name string) error {
	if name == "" {
		return fmt.Errorf("name must not be empty")
	}
	u.Name = name
	return nil
}
