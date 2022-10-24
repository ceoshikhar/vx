package vx

import (
	"testing"
)

func TestSimpleString(t *testing.T) {
	type simpleString struct {
		a string `vx:"type=string,minLength=3"`
	}

	v := simpleString{a: "as"}
	ok, errors := ValidateStruct(v)

	if !ok {
		for _, err := range errors {
			t.Log(err.Error())
		}
		t.Error("Something went wrong while validating the struct, check the errors above.")
	}

	if len(errors) != 1 {
		for _, err := range errors {
			t.Log(err.Error())
		}
		t.Errorf("Expected to get exactly 1 validation error but got %v, check the errors above.", len(errors))
	}
}
