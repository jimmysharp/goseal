package app

import "example.com/testproject/domain"

// SHOULD REPORT: Direct initialization in a NON-generated file
func UseGeneratedUser() {
	_ = domain.User{ // want "direct construction of struct User is prohibited outside allowed scope"
		ID:   3,
		Name: "another",
		Age:  35,
	}

	_ = &domain.User{ // want "direct construction of struct User is prohibited outside allowed scope"
		ID:   2,
		Name: "user-created",
		Age:  25,
	}
}

// SHOULD REPORT: Direct assignment in a NON-generated file
func DirectAssignment() {
	u := &domain.User{ // want "direct construction of struct User is prohibited outside allowed scope"
		ID:   1,
		Name: "test",
		Age:  40,
	}
	u.ID = 999 // want "direct assignment to field ID of struct User is prohibited outside allowed scope"
}
