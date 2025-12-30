package app

import "example.com/testproject/domain"

// SHOULD NOT REPORT: Using factory function
func WithFactoryFunction() {
	user, _ := domain.NewUser(123, "Alice", 30)
	_ = user
}

// SHOULD REPORT: Direct initialization without factory function
func WithoutFactoryFunction() {
	_ = domain.User{ // want "direct construction of struct User is prohibited outside allowed scope"
		ID:   123,
		Name: "Bob",
		Age:  25,
	}
	_ = domain.User{} // want "direct construction of struct User is prohibited outside allowed scope"

	_ = &domain.User{ // want "direct construction of struct User is prohibited outside allowed scope"
		ID:   123,
		Name: "Bob",
		Age:  25,
	}
	_ = &domain.User{} // want "direct construction of struct User is prohibited outside allowed scope"
}

// SHOULD NOT REPORT: Assignment through method
func AssignmentInReceiver() {
	user, _ := domain.NewUser(123, "Charlie", 35)

	user.UpdateName("Dave")
}

// SHOULD REPORT: Direct field assignment
func DirectAssignment() {
	user, _ := domain.NewUser(123, "Charlie", 35)

	user.ID = 456      // want "direct assignment to field ID of struct User is prohibited outside allowed scope"
	user.Name = "Dave" // want "direct assignment to field Name of struct User is prohibited outside allowed scope"
	user.Age = 40      // want "direct assignment to field Age of struct User is prohibited outside allowed scope"
}

type StructInSamePackage struct {
	message string
}

// SHOULD NOT REPORT: Only "domain" package is targeted (struct-packages)
func InSamePackage() {
	_ = StructInSamePackage{
		message: "Hello, World!",
	}
	_ = &StructInSamePackage{
		message: "Hello, World!",
	}
}

// SHOULD REPORT: Direct initialization in slice
func InSlice() {
	_ = []domain.User{
		{ // want "direct construction of struct User is prohibited outside allowed scope"
			ID:   1,
			Name: "Alice",
			Age:  30,
		},
		{ // want "direct construction of struct User is prohibited outside allowed scope"
			ID:   2,
			Name: "Bob",
			Age:  25,
		},
	}

	_ = []*domain.User{
		{ // want "direct construction of struct User is prohibited outside allowed scope"
			ID:   1,
			Name: "Alice",
			Age:  30,
		},
		{ // want "direct construction of struct User is prohibited outside allowed scope"
			ID:   2,
			Name: "Bob",
			Age:  25,
		},
	}

	_ = [][]domain.User{
		{
			domain.User{ // want "direct construction of struct User is prohibited outside allowed scope"
				ID:   1,
				Name: "Alice",
				Age:  30,
			},
		},
		{
			{ // want "direct construction of struct User is prohibited outside allowed scope"
				ID:   2,
				Name: "Bob",
				Age:  25,
			},
		},
	}
}

// SHOULD REPORT: Direct initialization in map
func InMap() {
	_ = map[int]domain.User{
		1: { // want "direct construction of struct User is prohibited outside allowed scope"
			ID:   1,
			Name: "Alice",
			Age:  30,
		},
		2: { // want "direct construction of struct User is prohibited outside allowed scope"
			ID:   2,
			Name: "Bob",
			Age:  25,
		},
	}

	_ = map[int]*domain.User{
		1: { // want "direct construction of struct User is prohibited outside allowed scope"
			ID:   1,
			Name: "Alice",
			Age:  30,
		},
		2: { // want "direct construction of struct User is prohibited outside allowed scope"
			ID:   2,
			Name: "Bob",
			Age:  25,
		},
	}
}

// SHOULD REPORT: Direct initialization in array
func InArray() {
	users := [2]domain.User{
		{ // want "direct construction of struct User is prohibited outside allowed scope"
			ID:   1,
			Name: "Alice",
			Age:  30,
		},
		{ // want "direct construction of struct User is prohibited outside allowed scope"
			ID:   2,
			Name: "Bob",
			Age:  25,
		},
	}
	_ = users
}

type Wrapper struct {
	User domain.User
}

type PointerWrapper struct {
	User *domain.User
}

// SHOULD REPORT: Direct initialization in struct field
func InStructField() {
	_ = Wrapper{
		User: domain.User{ // want "direct construction of struct User is prohibited outside allowed scope"
			ID:   1,
			Name: "Alice",
			Age:  30,
		},
	}

	_ = &Wrapper{
		User: domain.User{ // want "direct construction of struct User is prohibited outside allowed scope"
			ID:   1,
			Name: "Alice",
			Age:  30,
		},
	}

	_ = PointerWrapper{
		User: &domain.User{ // want "direct construction of struct User is prohibited outside allowed scope"
			ID:   1,
			Name: "Alice",
			Age:  30,
		},
	}

	_ = &PointerWrapper{
		User: &domain.User{ // want "direct construction of struct User is prohibited outside allowed scope"
			ID:   1,
			Name: "Alice",
			Age:  30,
		},
	}
}
