package internal

import (
	"reflect"
	"testing"
)

func TestParseStruct(t *testing.T) {
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
			got, err := ParseStruct(test.arg)

			if err != nil {
				t.Errorf(err.Error())
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ParseStruct() = %v, want %v", got, test.want)
			}
		})
	}
}
