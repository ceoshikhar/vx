package vx

import (
	"reflect"
	"testing"
)

func TestParseStructFields(t *testing.T) {
	type user struct {
		name  string `vx:"required"`
		email string `vx:"min=3"`
	}

	tests := []struct {
		name string
		arg  interface{}
		want []StructField
	}{
		{
			name: "user",
			arg:  user{name: "Jon Doe", email: "jon@doe.com"},
			want: []StructField{{Name: "name", Type: "string", Tag: "required", Value: "Jon Doe"}, {Name: "email", Type: "string", Tag: "min=3", Value: "jon@doe.com"}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParseStructFields(test.arg)

			if err != nil {
				t.Errorf(err.Error())
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ParseStructFields() = %v, want %v", got, test.want)
			}
		})
	}
}
