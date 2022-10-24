package internal

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
			want: TYPE_UNSUPPORTED,
		},
		{
			name: "TYPE_FLOAT",
			arg:  "type=float",
			want: TYPE_UNSUPPORTED,
		},
		{
			name: "TYPE_INT",
			arg:  "type=int,required",
			want: TYPE_UNSUPPORTED,
		},
		{
			name: "TYPE_RUNE",
			arg:  "type=rune",
			want: TYPE_UNSUPPORTED,
		},
		{
			name: "TYPE_STRING",
			arg:  "required,type=string,minLength=3",
			want: TYPE_STRING,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := MakeTag(test.arg)

			if err != nil {
				t.Errorf(err.Error())
			}

			if got.Type != test.want {
				t.Errorf("got %v, want %v", got.Type, test.want)
			}
		})
	}

}
