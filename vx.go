package vx

import (
	"fmt"
	"reflect"
	"strings"
	"vx/internal"
)

type VxResult struct {
	Errors []error
}

func (v VxResult) Error() string {
	var sb strings.Builder

	for _, err := range v.Errors {
		sb.WriteString("\n")
		sb.WriteString(err.Error())
	}

	return sb.String()
}

func ValidateStruct(v any) (VxResult, bool) {
	// If ok is false that means the error was caused due to something else other
	// than validating the struct. Something went wrong before the validation.
	ok := true
	vxErr := VxResult{
		Errors: []error{},
	}

	parsedStruct, err := internal.ParseStruct(v)
	if err != nil {
		vxErr.Errors = append(vxErr.Errors, err)
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
			vxErr.Errors = append(vxErr.Errors, err)
		}

		tagMap[field.Name] = tag
		fieldMap[field.Name] = field
	}

	// We have an internal error, so we are not going to return them without
	// executing rules on the parsedStruct's `Fields`.
	if !ok {
		return vxErr, ok
	}

	for _, field := range parsedStruct.Fields {
		if field.Type != field.ValueType && field.Type.Kind() != reflect.Interface {
			vxErr.Errors = append(vxErr.Errors, fmt.Errorf("%s should be of type %s", field.Name, field.Type))
		}

		// This is a bit tricky. Although we have a "required" Rule to check and warn
		// the user if a required value is not present or "empty" but if the developer
		// forgets to add the "required" Rule and then reads this field, that would
		// lead to a panic during runtime, which is NOT SO GOOD !!
		//
		// So in order to prevent the runtime panics for Vx users, we will set the
		// field.Value to be a default of the type that was being expected. If the
		// expected type is any, we will default it to empty string.
		if field.Value == nil {
			fmt.Printf("YIKES! This is bad, %s is nil which can be a nasty runtime error.", field.Name)
		}

	}

	// If we have collected some errors and have reached here that means these
	// errors are due to the fact that some field(s) type casting failed.
	// We don't execute Rules and return here with type casting errors.
	if len(vxErr.Errors) > 0 {
		ok = true
		return vxErr, ok
	}

	for fieldName, tag := range tagMap {
		field := fieldMap[fieldName]

		for _, rule := range tag.Rules {
			err := rule.Exec(field)
			if err != nil {
				vxErr.Errors = append(vxErr.Errors, err)
			}
		}
	}

	fmt.Println("\nParsedStruct:", parsedStruct)
	fmt.Println("\nResult:", ok, vxErr)

	return vxErr, ok
}
