package app

import "example.com/testproject/domain"

func WithConstructor() {
	user, _ := domain.NewUser(123, "Alice", 30)
	_ = user
}

func WithoutConstructor() {
	_ = domain.User{ // want "direct construction of struct User is prohibited, use constructor function"
		ID:   123,
		Name: "Bob",
		Age:  25,
	}
	_ = domain.User{} // want "direct construction of struct User is prohibited, use constructor function"

	_ = &domain.User{ // want "direct construction of struct User is prohibited, use constructor function"
		ID:   123,
		Name: "Bob",
		Age:  25,
	}
	_ = &domain.User{} // want "direct construction of struct User is prohibited, use constructor function"
}

func DirectAssignment() {
	user, _ := domain.NewUser(123, "Charlie", 35)

	user.ID = 456      // want "direct assignment to field ID of struct User is prohibited, use constructor function"
	user.Name = "Dave" // want "direct assignment to field Name of struct User is prohibited, use constructor function"
	user.Age = 40      // want "direct assignment to field Age of struct User is prohibited, use constructor function"
}

type StructInSamePackage struct {
	message string
}

func InSamePackage() {
	_ = StructInSamePackage{
		message: "Hello, World!",
	}
	_ = &StructInSamePackage{
		message: "Hello, World!",
	}
}

func InSlice() {
	_ = []domain.User{
		{ // want "direct construction of struct User is prohibited, use constructor function"
			ID:   1,
			Name: "Alice",
			Age:  30,
		},
		{ // want "direct construction of struct User is prohibited, use constructor function"
			ID:   2,
			Name: "Bob",
			Age:  25,
		},
	}

	_ = []*domain.User{
		{ // want "direct construction of struct User is prohibited, use constructor function"
			ID:   1,
			Name: "Alice",
			Age:  30,
		},
		{ // want "direct construction of struct User is prohibited, use constructor function"
			ID:   2,
			Name: "Bob",
			Age:  25,
		},
	}

	_ = [][]domain.User{
		{
			domain.User{ // want "direct construction of struct User is prohibited, use constructor function"
				ID:   1,
				Name: "Alice",
				Age:  30,
			},
		},
		{
			{ // want "direct construction of struct User is prohibited, use constructor function"
				ID:   2,
				Name: "Bob",
				Age:  25,
			},
		},
	}
}

func InMap() {
	_ = map[int]domain.User{
		1: { // want "direct construction of struct User is prohibited, use constructor function"
			ID:   1,
			Name: "Alice",
			Age:  30,
		},
		2: { // want "direct construction of struct User is prohibited, use constructor function"
			ID:   2,
			Name: "Bob",
			Age:  25,
		},
	}

	_ = map[int]*domain.User{
		1: { // want "direct construction of struct User is prohibited, use constructor function"
			ID:   1,
			Name: "Alice",
			Age:  30,
		},
		2: { // want "direct construction of struct User is prohibited, use constructor function"
			ID:   2,
			Name: "Bob",
			Age:  25,
		},
	}
}

func InArray() {
	users := [2]domain.User{
		{ // want "direct construction of struct User is prohibited, use constructor function"
			ID:   1,
			Name: "Alice",
			Age:  30,
		},
		{ // want "direct construction of struct User is prohibited, use constructor function"
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

func InStructField() {
	_ = Wrapper{
		User: domain.User{ // want "direct construction of struct User is prohibited, use constructor function"
			ID:   1,
			Name: "Alice",
			Age:  30,
		},
	}

	_ = &Wrapper{
		User: domain.User{ // want "direct construction of struct User is prohibited, use constructor function"
			ID:   1,
			Name: "Alice",
			Age:  30,
		},
	}

	_ = PointerWrapper{
		User: &domain.User{ // want "direct construction of struct User is prohibited, use constructor function"
			ID:   1,
			Name: "Alice",
			Age:  30,
		},
	}

	_ = &PointerWrapper{
		User: &domain.User{ // want "direct construction of struct User is prohibited, use constructor function"
			ID:   1,
			Name: "Alice",
			Age:  30,
		},
	}
}
