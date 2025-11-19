package app

import "example.com/testproject/domain"

// Direct initialization in a NON-generated file - SHOULD be detected
func UseGeneratedUser() {
	_ = domain.User{ // want "direct construction of struct User is prohibited, use constructor function"
		ID:   3,
		Name: "another",
		Age:  35,
	}

	_ = &domain.User{ // want "direct construction of struct User is prohibited, use constructor function"
		ID:   2,
		Name: "user-created",
		Age:  25,
	}
}

// Direct assignment in a NON-generated file - SHOULD be detected
func DirectAssignment() {
	u := &domain.User{ // want "direct construction of struct User is prohibited, use constructor function"
		ID:   1,
		Name: "test",
		Age:  40,
	}
	u.ID = 999 // want "direct assignment to field ID of struct User is prohibited, use constructor function"
}
