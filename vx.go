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

func ValidateStruct(v any) (VxError, bool) {
	// If ok is false that means the error was caused due to something else other
	// than validating the struct. Something went wrong before the validation.
	ok := true
	vxErr := VxError{
		errors: []error{},
	}

	parsedStruct, err := internal.ParseStruct(v)
	if err != nil {
		vxErr.errors = append(vxErr.errors, err)
		ok = false
		// This is the only case where we early return because if this fails
		// there is literally nothing more we can do after this.
		return vxErr, ok
	}

	// FIXME: Making these maps might not be the best way to do this.
	fieldMap := map[string]internal.StructField{}
	tagMap := map[string]internal.Tag{}

	// NOTE: do not return if `MakeTag` fails for a field. We want to collect
	// all the errors for all fields first and then return them at the end.
	for _, field := range parsedStruct.Fields {
		tag, err := internal.MakeTag(field)
		if err != nil {
			err = fmt.Errorf("%s.%s - %s", parsedStruct.Name, field.Name, err)
			ok = false
			vxErr.errors = append(vxErr.errors, err)
		}
		tagMap[field.Name] = tag
		fieldMap[field.Name] = field
	}

	// We have an internal error, so we are not going to return them without
	// executing rules on the parsedStruct's `Fields`.
	if !ok {
		return vxErr, ok
	}

	for fieldName, tag := range tagMap {
		field := fieldMap[fieldName]
		for _, rule := range tag.Rules {
			err := rule.Exec(field.Value)
			if err != nil {
				err = fmt.Errorf("%s.%s - %s", parsedStruct.Name, field.Name, err)
				vxErr.errors = append(vxErr.errors, err)
			}
		}
	}

	// fmt.Println(v, structFields, vxErr)

	return vxErr, ok
}
