package vx

import (
	"fmt"
	"reflect"
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
	fieldMap := map[string]internal.VxField{}
	tagMap := map[string]internal.Tag{}

	// NOTE: do not return if `MakeTag` fails for a field. We want to collect
	// all the errors for all fields first and then return them at the end.
	for _, field := range parsedStruct.Fields {
		tag, err := internal.MakeTag(field)
		if err != nil {
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

	for fieldName, field := range fieldMap {
		switch field.Type.Kind() {
		case reflect.Int:
			{
				field := fieldMap[fieldName]
				v, ok := field.Value.(int)
				if !ok {
					vxErr.errors = append(vxErr.errors, fmt.Errorf("%s should be an int", field.Name))
				}

				field.Value = v
				break
			}
		case reflect.String:
			{
				field := fieldMap[fieldName]
				v, ok := field.Value.(string)
				if !ok {
					vxErr.errors = append(vxErr.errors, fmt.Errorf("%s should be a string", field.Name))
				}

				field.Value = v
				break
			}
		default:
			{
				// Nothing to do here.
			}
		}
	}

	// If we have collected some errors and have reached here that means these
	// errors are due to the fact that some field(s) type casting failed.
	// We don't execute Rules and return here with type casting errors.
	if len(vxErr.errors) > 0 {
		ok = true
		return vxErr, ok
	}

	for fieldName, tag := range tagMap {
		field := fieldMap[fieldName]

		for _, rule := range tag.Rules {
			err := rule.Exec(field)
			if err != nil {
				vxErr.errors = append(vxErr.errors, err)
			}
		}
	}

	// fmt.Println("ParsedStruct:", parsedStruct)
	// fmt.Println("Result:", ok, vxErr)

	return vxErr, ok
}
