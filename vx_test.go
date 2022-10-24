package vx

import (
	"testing"
)

func TestValidateStruct(t *testing.T) {
	type noTag struct {
		a string
	}

	type emptyTag struct {
		a string `vx:""`
	}

	type ruleMinLength5 struct {
		a string `vx:"minLength=5"`
	}

	type ruleMinLength0 struct {
		a string `vx:"minLength=0"`
	}

	type ruleMinLength3WithAny struct {
		a interface{} `vx:"minLength=3"`
	}

	type want struct {
		ok    bool
		count int
	}

	tests := []struct {
		name string
		arg  interface{}
		want want
	}{
		{
			name: "noTag",
			arg: noTag{
				a: "vx",
			},
			want: want{true, 0},
		},
		{
			name: "emptyTag",
			arg: emptyTag{
				a: "vx",
			},
			want: want{true, 0},
		},
		{
			name: "rule: minLength of 5 is failed",
			arg: ruleMinLength5{
				a: "yolo",
			},
			want: want{true, 1},
		},
		{
			name: "rule: minLength of 5 is passed",
			arg: ruleMinLength5{
				a: "happy",
			},
			want: want{true, 0},
		},
		{
			name: "rule: minLength of 0 should throw an internal error",
			arg: ruleMinLength0{
				a: "happy",
			},
			want: want{false, 1},
		},
		{
			name: "rule: minLength of 3 with field of type any",
			arg: ruleMinLength3WithAny{
				a: "ab",
			},
			want: want{true, 1},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ok, err := ValidateStruct(test.arg)

			if ok != test.want.ok {
				t.Error(err.Error())
				t.Errorf("expected ok to be %v but got %v, check the errors above", ok, test.want.ok)
			}

			if len(err.errors) != test.want.count {
				t.Error(err.Error())
				t.Errorf("expected to get exactly %v number of validation errors but got %v, check the errors above.", test.want.count, len(err.errors))
			}

		})
	}
}
