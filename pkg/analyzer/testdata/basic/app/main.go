package app

import "example.com/testproject/domain"

type Wrapper struct {
	User domain.User
}

type PointerWrapper struct {
	User *domain.User
}

func WithConstructor() {
	user := domain.NewUser(123, "Alice", 30)
	_ = user
}

func WithoutConstructor() {
	user := domain.User{ // want "direct construction of struct User is prohibited, use constructor function"
		ID:   123,
		Name: "Bob",
		Age:  25,
	}
	_ = user
}

func WithoutConstructorByPointer() {
	user := &domain.User{ // want "direct construction of struct User is prohibited, use constructor function"
		ID:   123,
		Name: "Bob",
		Age:  25,
	}
	_ = user
}

func DirectAssignment() {
	user := domain.NewUser(123, "Charlie", 35)

	user.Name = "Dave" // want "direct assignment to field Name of struct User is prohibited, use constructor function"
}

func InSlice() {
	users := []domain.User{
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

func InSlicePointer() {
	users := []*domain.User{
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

func InMap() {
	userMap := map[int]domain.User{
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
	_ = userMap
}

func InMapPointer() {
	userMap := map[int]*domain.User{
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
	_ = userMap
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

func InNestedSlice() {
	nestedUsers := [][]domain.User{
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
	_ = nestedUsers
}

func InStructField() {
	wrapper := Wrapper{
		User: domain.User{ // want "direct construction of struct User is prohibited, use constructor function"
			ID:   1,
			Name: "Alice",
			Age:  30,
		},
	}
	_ = wrapper
}

func InStructFieldPointer() {
	pointerWrapper := PointerWrapper{
		User: &domain.User{ // want "direct construction of struct User is prohibited, use constructor function"
			ID:   1,
			Name: "Alice",
			Age:  30,
		},
	}
	_ = pointerWrapper
}
