package vx

import (
	"strings"
	"testing"
)

func TestValidateStruct(t *testing.T) {
	type simpleString struct {
		a string `vx:"minLength=3"`
	}

	v := simpleString{a: "1"}
	ok, err := ValidateStruct(v)

	if !ok {
		t.Error("Something went wrong while validating the struct:")
		t.Error(err.Error())
	}

	if len(err.errors) != 1 {
		t.Error(err.Error())
		t.Errorf("expected to get exactly 1 validation error but got %v, check the errors above.", len(err.errors))
	}

	if !strings.Contains(err.Error(), "minLength") {
		t.Error("minLength rule was supposed to fail but found no mention of it in the error")
	}
}
