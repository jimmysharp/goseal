package domain

import (
	"reflect"
	"testing"
)

// SHOULD REPORT: No ignore-files configured, so test files are also checked (ignore-files)
func TestNewUser_Success(t *testing.T) {
	type args struct {
		id   int
		name string
		age  int
	}

	tests := []struct {
		name string
		args args
		want *User
	}{
		{
			name: "normal case",
			args: args{ // want "direct construction of struct args is prohibited, use constructor function"
				id:   1,
				name: "Alice",
				age:  28,
			},
			want: &User{ // want "direct construction of struct User is prohibited, use constructor function"
				ID:   1,
				Name: "Alice",
				Age:  28,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUser(tt.args.id, tt.args.name, tt.args.age)
			if err != nil {
				t.Errorf("NewUser() unexpected error = %v", err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

// SHOULD REPORT: No ignore-files configured, so test files are also checked (ignore-files)
func TestNewUser_Error(t *testing.T) {
	type args struct {
		id   int
		name string
		age  int
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "id is zero",
			args: args{ // want "direct construction of struct args is prohibited, use constructor function"
				id:   0,
				name: "Alice",
				age:  30,
			},
		},
		{
			name: "name is empty",
			args: args{ // want "direct construction of struct args is prohibited, use constructor function"
				id:   1,
				name: "",
				age:  30,
			},
		},
		{
			name: "age is negative",
			args: args{ // want "direct construction of struct args is prohibited, use constructor function"
				id:   1,
				name: "Alice",
				age:  -1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUser(tt.args.id, tt.args.name, tt.args.age)
			if err == nil {
				t.Errorf("NewUser() expected error but got nil")
				return
			}
			if got != nil {
				t.Errorf("NewUser() = %v, want nil", got)
			}
		})
	}
}
