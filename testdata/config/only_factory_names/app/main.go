package app

import "example.com/testproject/domain"

// SHOULD NOT REPORT: Using constructor function matching factory-names
func WithConstructor() {
	user, _ := domain.NewUser(123, "Alice", 30)
	_ = user

	defaultUser := domain.CreateDefaultUser()
	_ = defaultUser
}

// SHOULD REPORT: Initialization outside factory function
func WithoutConstructor() {
	_ = domain.User{ // want "direct construction of struct User is prohibited, use allowed factory function"
		ID:   123,
		Name: "Bob",
		Age:  25,
	}

	_ = &domain.User{ // want "direct construction of struct User is prohibited, use allowed factory function"
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
	_ = MyStruct{ // want "direct construction of struct MyStruct is prohibited, use allowed factory function"
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

	user.ID = 456      // want "direct assignment to field ID of struct User is prohibited outside allowed scope"
	user.Name = "Dave" // want "direct assignment to field Name of struct User is prohibited outside allowed scope"
	user.Age = 40      // want "direct assignment to field Age of struct User is prohibited outside allowed scope"
}
