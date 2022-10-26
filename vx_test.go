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

	type ruleMinLength3WithInt struct {
		a int `vx:"minLength=3"`
	}

	type twoStringFields struct {
		a string `vx:"minLength=0"`
		b string `vx:"minLength=ab"`
	}

	type structAnyTagInt struct {
		age interface{} `vx:"type=int"`
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
			name: "rule: minLength of 3 with field of type any with string value should fail",
			arg: ruleMinLength3WithAny{
				a: "ab",
			},
			want: want{true, 1},
		},
		{
			name: "rule: minLength of 3 with field of type any with int value should fail",
			arg: ruleMinLength3WithAny{
				a: 12,
			},
			want: want{true, 1},
		},
		{
			name: "rule: minLength of 3 with field of type any with string value should pass",
			arg: ruleMinLength3WithAny{
				a: "abcd",
			},
			want: want{true, 0},
		},
		{
			name: "rule: minLength of 3 with field of type any with int value should pass",
			arg: ruleMinLength3WithAny{
				a: 1234,
			},
			want: want{true, 0},
		},
		{
			name: "rule: minLength of 3 with field of type any with bool value should pass",
			arg: ruleMinLength3WithAny{
				a: false,
			},
			want: want{true, 0},
		},
		{
			name: "rule: minLength of 3 with with int value should fail",
			arg: ruleMinLength3WithInt{
				a: 1234,
			},
			want: want{false, 1},
		},
		{
			name: "rule: minlength should return internal error for both fields",
			arg: twoStringFields{
				a: "ab",
				b: "abc",
			},
			want: want{false, 2},
		},
		{
			name: "type expected to be int error",
			arg: structAnyTagInt{
				age: "abc",
			},
			want: want{true, 1},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, ok := ValidateStruct(test.arg)

			if ok != test.want.ok {
				t.Error(res.Error())
				t.Errorf("expected ok to be %v but got %v, check the errors above", test.want.ok, ok)
			}

			if len(res.errors) != test.want.count {
				t.Error(res.Error())
				t.Errorf("expected to get exactly %v number of validation errors but got %v, check the errors above.", test.want.count, len(res.errors))
			}

		})
	}
}
