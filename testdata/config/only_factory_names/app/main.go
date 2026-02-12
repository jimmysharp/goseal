package app

import "example.com/testproject/domain"

// SHOULD NOT REPORT: Using factory function matching factory-names
func WithFactoryFunction() {
	user, _ := domain.NewUser(123, "Alice", 30)
	_ = user

	defaultUser := domain.CreateDefaultUser()
	_ = defaultUser
}

// SHOULD REPORT: Initialization outside factory function
func WithoutFactoryFunction() {
	_ = domain.User{ // want "direct construction of sealed struct User is not allowed outside factory functions \\(factory-names\\)"
		ID:   123,
		Name: "Bob",
		Age:  25,
	}

	_ = &domain.User{ // want "direct construction of sealed struct User is not allowed outside factory functions \\(factory-names\\)"
		ID:   456,
		Name: "Charlie",
		Age:  30,
	}
}

// SHOULD NOT REPORT: Function matching "^New.*" is considered a factory (factory-names)
func NewUser(id int, name string, age int) *domain.User {
	return &domain.User{
		ID:   id,
		Name: name,
		Age:  age,
	}
}

type MyStruct struct {
	Value int
}

// SHOULD REPORT: target-packages is empty (all packages targeted), and not in factory function
func InitLocalStruct() {
	_ = MyStruct{ // want "direct construction of sealed struct MyStruct is not allowed outside factory functions \\(factory-names\\)"
		Value: 100,
	}
}

// SHOULD NOT REPORT: Function matching "^New.*" is considered a factory (factory-names)
func NewMyStruct(value int) *MyStruct {
	return &MyStruct{
		Value: value,
	}
}

// SHOULD REPORT: Direct field assignment outside receiver
func DirectAssignment() {
	user, _ := domain.NewUser(123, "Charlie", 35)

	user.ID = 456      // want "direct assignment to field ID of sealed struct User is not allowed outside its receiver methods \\(mutation-scope: receiver\\)"
	user.Name = "Dave" // want "direct assignment to field Name of sealed struct User is not allowed outside its receiver methods \\(mutation-scope: receiver\\)"
	user.Age = 40      // want "direct assignment to field Age of sealed struct User is not allowed outside its receiver methods \\(mutation-scope: receiver\\)"
}
