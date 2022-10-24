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
		return false, vxErr
	}

	for _, field := range structFields {
		tag, err := internal.MakeTag(field)

		if err != nil {
			vxErr.errors = append(vxErr.errors, err)
			return false, vxErr
		}

		for _, rule := range tag.Rules {
			err := rule.Exec(field.Value)

			if err != nil {
				vxErr.errors = append(vxErr.errors, err)
			}
		}
	}

	return true, vxErr
}
