package domain

func BuildUserInSamePackage() *User {
	return &User{
		ID:   999,
		Name: "SamePackage",
		Age:  40,
	}
}

func UpdateUserInSamePackage(u *User) {
	u.Name = "Updated"
}
