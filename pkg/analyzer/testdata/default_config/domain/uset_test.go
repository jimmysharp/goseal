package domain

import (
	"reflect"
	"testing"
)

func TestNewUser(t *testing.T) {
	type args struct {
		id   int
		name string
		age  int
	}

	tests := []struct {
		name string
		args args
		want User
	}{
		{
			name: "normal case",
			args: args{ // want "direct construction of struct args is prohibited, use constructor function"
				id:   1,
				name: "Alice",
				age:  28,
			},
			want: User{ // want "direct construction of struct User is prohibited, use constructor function"
				ID:   1,
				Name: "Alice",
				Age:  28,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUser(tt.args.id, tt.args.name, tt.args.age)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
