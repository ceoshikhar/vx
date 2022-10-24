package vx

import (
	"strings"
	"testing"
)

func TestSimpleString(t *testing.T) {
	type simpleString struct {
		a string `vx:"type=string,minLength=3"`
	}

	v := simpleString{a: "as"}
	ok, err := ValidateStruct(v)

	if !ok {
		t.Error("Something went wrong while validating the struct:")
		t.Error(err.Error())
	}

	if !strings.Contains(err.Error(), "minLength") {
		t.Error("minLength rule was supposed to fail but found no mention of it in the error")
	}
}
