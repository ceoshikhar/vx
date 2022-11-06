package vx

import (
	"testing"
)

type want struct {
	ok    bool
	count int
}

type validateStructTest struct {
	name string
	arg  any
	want want
}

func runValidateStructTests(tests []validateStructTest, t *testing.T) {
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

func TestTypeInTag(t *testing.T) {
	type anyToBool struct {
		A any `vx:"type=bool"`
	}

	type anyToInt struct {
		A any `vx:"type=int"`
	}

	type anyToFloat64 struct {
		A any `vx:"type=float64"`
	}

	type anyToString struct {
		A any `vx:"type=string"`
	}

	type anyToSlice struct {
		A any `vx:"type=[]string"`
	}

	type anyToArray struct {
		A any `vx:"type=[10]string"`
	}

	type anyToMap struct {
		A any `vx:"type=map[string]string"`
	}

	type anyToAny struct {
		A any `vx:"type=any"`
	}

	type anyToStringSlice struct {
		A any `vx:"type=[]string"`
	}

	type anyToStringToStringMap struct {
		A any `vx:"type=map[string]string"`
	}

	type anyToStringArray struct {
		A any `vx:"type=[10]string"`
	}

	type simple struct {
		A any `vx:"name=simpleA, type=string"`
	}

	type complex struct {
		A simple
	}

	tests := []validateStructTest{
		{
			name: "anyToBool with bool value should not give an error",
			arg: anyToBool{
				A: false,
			},
			want: want{true, 0},
		},
		{
			name: "anyToInt with int value should not give an error",
			arg: anyToInt{
				A: 0,
			},
			want: want{true, 0},
		},
		{
			name: "anyToFloat64 with float value should not give an error",
			arg: anyToFloat64{
				A: 0.0,
			},
			want: want{true, 0},
		},
		{
			name: "anyToString with string value should not give an error",
			arg: anyToString{
				A: "",
			},
			want: want{true, 0},
		},
		{
			name: "anyToSlice with []string value should not give an error",
			arg: anyToSlice{
				A: []string{},
			},
			want: want{true, 0},
		},
		{
			name: "anyToArray with [10]string{} value should not give an error",
			arg: anyToArray{
				A: [10]string{},
			},
			want: want{true, 0},
		},
		{
			name: "anyToMap with map[string]string{} value should not give an error",
			arg: anyToMap{
				A: map[string]string{},
			},
			want: want{true, 0},
		},
		{
			name: "anyToAny with []string{} value should not give an error",
			arg: anyToAny{
				A: []string{},
			},
			want: want{true, 0},
		},
		{
			name: "anyToStringSlice with []string{} value should not give an error",
			arg: anyToAny{
				A: []string{},
			},
			want: want{true, 0},
		},
		{
			name: "anyToStringSlice with []string value should not give an error",
			arg: anyToStringSlice{
				A: []string{},
			},
			want: want{true, 0},
		},
		{
			name: "anyToStringToStringMap with map[string]string value should not give an error",
			arg: anyToStringToStringMap{
				A: map[string]string{},
			},
			want: want{true, 0},
		},
		{
			name: "anyToBool with string value should give an error",
			arg: anyToBool{
				A: "",
			},
			want: want{true, 1},
		},
		{
			name: "anyToStringSlice with []any value should not give an error",
			arg: anyToStringSlice{
				A: []any{},
			},
			want: want{true, 0},
		},
		{
			name: "anyToStringArray with [10]string value should not give an error",
			arg: anyToStringArray{
				A: [10]any{},
			},
			want: want{true, 0},
		},
		{
			name: "anyToStringToStringMap with map[any]any value should not give an error",
			arg: anyToStringToStringMap{
				A: map[any]any{},
			},
			want: want{true, 0},
		},
		{
			name: "anyToStringToStringMap with map[string]any value should not give an error",
			arg: anyToStringToStringMap{
				A: map[string]any{},
			},
			want: want{true, 0},
		},
		{
			name: "anyToStringToStringMap with map[any]string value should not give an error",
			arg: anyToStringToStringMap{
				A: map[any]string{},
			},
			want: want{true, 0},
		},
		{
			name: "complex with string value should not give an error",
			arg: complex{
				A: simple{A: "abc"},
			},
			want: want{true, 0},
		},
		{
			name: "complex with int value should give an error",
			arg: complex{
				A: simple{A: 123},
			},
			want: want{true, 1},
		},
	}

	runValidateStructTests(tests, t)
}

func TestRequired(t *testing.T) {
	type aAny struct {
		A any `vx:"required"`
	}

	tests := []validateStructTest{
		{
			name: "aAny with bool value should not give an error",
			arg: aAny{
				A: false,
			},
			want: want{true, 0},
		},
		{
			name: "aAny with int value should not give an error",
			arg: aAny{
				A: 0,
			},
			want: want{true, 0},
		},
		{
			name: "aAny with float64 value should not give an error",
			arg: aAny{
				A: 0.0,
			},
			want: want{true, 0},
		},
		{
			name: "aAny with non-empty string value should not give an error",
			arg: aAny{
				A: "can't be empty",
			},
			want: want{true, 0},
		},
		{
			name: "aAny with slice value should not give an error",
			arg: aAny{
				A: []any{},
			},
			want: want{true, 0},
		},
		{
			name: "aAny with map value should not give an error",
			arg: aAny{
				A: map[any]any{},
			},
			want: want{true, 0},
		},
		{
			name: "aAny with nil value should get an error",
			arg:  aAny{},
			want: want{true, 1},
		},
	}

	runValidateStructTests(tests, t)
}

// func TestValidateStruct(t *testing.T) {
// 	type noTag struct {
// 		A string
// 	}

// 	type emptyTag struct {
// 		A string `vx:""`
// 	}

// 	type ruleMinLength5 struct {
// 		A string `vx:"minLength=5"`
// 	}

// 	type ruleMinLength0 struct {
// 		A string `vx:"minLength=0"`
// 	}

// 	type ruleMinLength3WithAny struct {
// 		A any `vx:"minLength=3"`
// 	}

// 	type ruleMinLength3WithInt struct {
// 		Lmao int `vx:"minLength=3"`
// 	}

// 	type twoStringFields struct {
// 		A string `vx:"minLength=0"`
// 		B string `vx:"minLength=ab"`
// 	}

// 	type ruleRequired struct {
// 		A any `vx:"required"`
// 	}

// 	type want struct {
// 		ok    bool
// 		count int
// 	}

// 	tests := []struct {
// 		name string
// 		arg  interface{}
// 		want want
// 	}{
// 		{
// 			name: "noTag",
// 			arg: noTag{
// 				A: "vx",
// 			},
// 			want: want{true, 0},
// 		},
// 		{
// 			name: "emptyTag",
// 			arg: emptyTag{
// 				A: "vx",
// 			},
// 			want: want{true, 0},
// 		},
// 		{
// 			name: "rule: minLength of 5 is failed",
// 			arg: ruleMinLength5{
// 				A: "yolo",
// 			},
// 			want: want{true, 1},
// 		},
// 		{
// 			name: "rule: minLength of 5 is passed",
// 			arg: ruleMinLength5{
// 				A: "happy",
// 			},
// 			want: want{true, 0},
// 		},
// 		{
// 			name: "rule: minLength of 0 should throw an internal error",
// 			arg: ruleMinLength0{
// 				A: "happy",
// 			},
// 			want: want{false, 1},
// 		},
// 		{
// 			name: "rule: minLength of 3 with field of type any with string value should fail",
// 			arg: ruleMinLength3WithAny{
// 				A: "ab",
// 			},
// 			want: want{true, 1},
// 		},
// 		{
// 			name: "rule: minLength of 3 with field of type any with int value should fail",
// 			arg: ruleMinLength3WithAny{
// 				A: 12,
// 			},
// 			want: want{true, 1},
// 		},
// 		{
// 			name: "rule: minLength of 3 with field of type any with string value should pass",
// 			arg: ruleMinLength3WithAny{
// 				A: "abcd",
// 			},
// 			want: want{true, 0},
// 		},
// 		{
// 			name: "rule: minLength of 3 with field of type any with int value should fail",
// 			arg: ruleMinLength3WithAny{
// 				A: 1234,
// 			},
// 			want: want{true, 1},
// 		},
// 		{
// 			name: "rule: minLength of 3 with field of type any with bool value should fail",
// 			arg: ruleMinLength3WithAny{
// 				A: false,
// 			},
// 			want: want{true, 1},
// 		},
// 		{
// 			name: "rule: minLength of 3 with with int value should fail",
// 			arg: ruleMinLength3WithInt{
// 				Lmao: 1234,
// 			},
// 			want: want{true, 1},
// 		},
// 		{
// 			name: "rule: minlength should return internal error for both fields",
// 			arg: twoStringFields{
// 				A: "ab",
// 				B: "abc",
// 			},
// 			want: want{false, 2},
// 		},
// 		{
// 			name: "rule: required should fail because string is empty",
// 			arg: ruleRequired{
// 				A: "",
// 			},
// 			want: want{true, 1},
// 		},
// 		{
// 			name: "rule: required should fail because field is nil",
// 			arg: ruleRequired{
// 				A: nil,
// 			},
// 			want: want{true, 1},
// 		},
// 		{
// 			name: "rule: required should pass because 0 is valid",
// 			arg: ruleRequired{
// 				A: 0,
// 			},
// 			want: want{true, 0},
// 		},
// 		{
// 			name: "rule: required should pass because false is valid",
// 			arg: ruleRequired{
// 				A: false,
// 			},
// 			want: want{true, 0},
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			res, ok := ValidateStruct(test.arg)

// 			if ok != test.want.ok {
// 				t.Error(res.String())
// 				t.Errorf("expected ok to be %v but got %v, check the errors above", test.want.ok, ok)
// 			}

// 			if len(res.Errors) != test.want.count {
// 				t.Error(res.String())
// 				t.Errorf("expected to get exactly %v number of validation errors but got %v, check the errors above.", test.want.count, len(res.Errors))
// 			}
// 		})
// 	}
// }
