package app

import "example.com/testproject/domain"

// SHOULD REPORT: Direct initialization in a NON-generated file
func UseGeneratedUser() {
	_ = domain.User{ // want "direct construction of sealed struct User is not allowed from outside its package \\(init-scope: same-package\\)"
		ID:   3,
		Name: "another",
		Age:  35,
	}

	_ = &domain.User{ // want "direct construction of sealed struct User is not allowed from outside its package \\(init-scope: same-package\\)"
		ID:   2,
		Name: "user-created",
		Age:  25,
	}
}

// SHOULD REPORT: Direct assignment in a NON-generated file
func DirectAssignment() {
	u := &domain.User{ // want "direct construction of sealed struct User is not allowed from outside its package \\(init-scope: same-package\\)"
		ID:   1,
		Name: "test",
		Age:  40,
	}
	u.ID = 999 // want "direct assignment to field ID of sealed struct User is not allowed outside its receiver methods \\(mutation-scope: receiver\\)"
}
