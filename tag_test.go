package vx

import (
	"testing"
)

// Tests all the possible kind of values for `Type` field of `Tag`.
func TestAllVxTypePossible(t *testing.T) {

	tests := []struct {
		name string
		arg  string
		want VxType
	}{
		{
			name: "TYPE_BYTE",
			arg:  "type=byte",
			want: TYPE_BYTE,
		},
		{
			name: "TYPE_FLOAT",
			arg:  "type=float",
			want: TYPE_FLOAT,
		},
		{
			name: "TYPE_INT",
			arg:  "type=int,required",
			want: TYPE_INT,
		},
		{
			name: "TYPE_RUNE",
			arg:  "type=rune",
			want: TYPE_RUNE,
		},
		{
			name: "TYPE_STRING",
			arg:  "required,type=string,minLength=3",
			want: TYPE_STRING,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := MakeTagFromString(test.arg)

			if got.Type != test.want {
				t.Errorf("got %v, want %v", got.Type, test.want)
			}
		})
	}

}
