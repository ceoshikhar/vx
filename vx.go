package vx

import (
	"fmt"
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

func ValidateStruct(v any) (ok bool, err error) {
	structFields, err := internal.ParseStruct(v)

	if err != nil {
		return false, err
	}

	vxErr := VxError{
		errors: []error{},
	}

	for _, field := range structFields {
		tag, err := internal.MakeTag(field.Tag)

		if tag.Type != field.Type {
			err = fmt.Errorf("type mismatch - field '%s' type in struct is '%s' and type in tag is '%s'", field.Name, field.Type, tag.Type)
			return false, err
		}

		if err != nil {
			return false, err
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
