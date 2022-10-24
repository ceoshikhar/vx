package internal

import (
	"errors"
	"fmt"
)

//
// String specific rules
//

type minLength struct {
	value int
}

func makeMinLength(l int) minLength {
	return minLength{l}
}

func (m minLength) Exec(v any) error {
	s, ok := v.(string)

	if !ok {
		return errors.New("minLength: rule was exec against a value that is not a string")
	}

	if len(s) < m.value {
		return fmt.Errorf("minLength: minimum length allowed was 3 but got %v", len(s))
	}

	return nil
}
