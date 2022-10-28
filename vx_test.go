package vx

import (
	"testing"
)

func TestValidateStruct(t *testing.T) {
	type noTag struct {
		A string
	}

	type emptyTag struct {
		A string `vx:""`
	}

	type ruleMinLength5 struct {
		A string `vx:"minLength=5"`
	}

	type ruleMinLength0 struct {
		A string `vx:"minLength=0"`
	}

	type ruleMinLength3WithAny struct {
		A any `vx:"minLength=3"`
	}

	type ruleMinLength3WithInt struct {
		Lmao int `vx:"minLength=3"`
	}

	type twoStringFields struct {
		A string `vx:"minLength=0"`
		B string `vx:"minLength=ab"`
	}

	type ruleRequired struct {
		A any `vx:"required"`
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
				A: "vx",
			},
			want: want{true, 0},
		},
		{
			name: "emptyTag",
			arg: emptyTag{
				A: "vx",
			},
			want: want{true, 0},
		},
		{
			name: "rule: minLength of 5 is failed",
			arg: ruleMinLength5{
				A: "yolo",
			},
			want: want{true, 1},
		},
		{
			name: "rule: minLength of 5 is passed",
			arg: ruleMinLength5{
				A: "happy",
			},
			want: want{true, 0},
		},
		{
			name: "rule: minLength of 0 should throw an internal error",
			arg: ruleMinLength0{
				A: "happy",
			},
			want: want{false, 1},
		},
		{
			name: "rule: minLength of 3 with field of type any with string value should fail",
			arg: ruleMinLength3WithAny{
				A: "ab",
			},
			want: want{true, 1},
		},
		{
			name: "rule: minLength of 3 with field of type any with int value should fail",
			arg: ruleMinLength3WithAny{
				A: 12,
			},
			want: want{true, 1},
		},
		{
			name: "rule: minLength of 3 with field of type any with string value should pass",
			arg: ruleMinLength3WithAny{
				A: "abcd",
			},
			want: want{true, 0},
		},
		{
			name: "rule: minLength of 3 with field of type any with int value should fail",
			arg: ruleMinLength3WithAny{
				A: 1234,
			},
			want: want{true, 1},
		},
		{
			name: "rule: minLength of 3 with field of type any with bool value should fail",
			arg: ruleMinLength3WithAny{
				A: false,
			},
			want: want{true, 1},
		},
		{
			name: "rule: minLength of 3 with with int value should fail",
			arg: ruleMinLength3WithInt{
				Lmao: 1234,
			},
			want: want{true, 1},
		},
		{
			name: "rule: minlength should return internal error for both fields",
			arg: twoStringFields{
				A: "ab",
				B: "abc",
			},
			want: want{false, 2},
		},
		{
			name: "rule: required should fail because string is empty",
			arg: ruleRequired{
				A: "",
			},
			want: want{true, 1},
		},
		{
			name: "rule: required should fail because field is nil",
			arg: ruleRequired{
				A: nil,
			},
			want: want{true, 1},
		},
		{
			name: "rule: required should pass because 0 is valid",
			arg: ruleRequired{
				A: 0,
			},
			want: want{true, 0},
		},
		{
			name: "rule: required should pass because false is valid",
			arg: ruleRequired{
				A: false,
			},
			want: want{true, 0},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, ok := ValidateStruct(test.arg)

			if ok != test.want.ok {
				t.Error(res.String())
				t.Errorf("expected ok to be %v but got %v, check the errors above", test.want.ok, ok)
			}

			if len(res.Errors) != test.want.count {
				t.Error(res.String())
				t.Errorf("expected to get exactly %v number of validation errors but got %v, check the errors above.", test.want.count, len(res.Errors))
			}
		})
	}
}
