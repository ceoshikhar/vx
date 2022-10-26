package vx

import (
	"strings"
	"vx/internal"
)

type VxError struct {
	errors []error
}

func (v VxError) Error() string {
	var sb strings.Builder

	for _, err := range v.errors {
		sb.WriteString("\n")
		sb.WriteString(err.Error())
	}

	return sb.String()
}

func ValidateStruct(v any) (bool, VxError) {
	vxErr := VxError{
		errors: []error{},
	}

	structFields, err := internal.ParseStruct(v)
	if err != nil {
		vxErr.errors = append(vxErr.errors, err)
		// This is the only case where we early return because if this fails
		// there is literally nothing more we can do after this.
		return false, vxErr
	}

	// NOTE: do not return if `MakeTag` fails for a field. We want to collect
	// all the errors for all fields first and then return them at the end.
	for _, field := range structFields {
		tag, err := internal.MakeTag(field)
		if err != nil {
			vxErr.errors = append(vxErr.errors, err)
		}

		for _, rule := range tag.Rules {
			err := rule.Exec(field.Value)
			if err != nil {
				vxErr.errors = append(vxErr.errors, err)
			}
		}
	}

	// fmt.Println(v, structFields, vxErr)

	return true, vxErr
}
