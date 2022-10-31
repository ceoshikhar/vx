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

func (v VxResult) String() string {
	var sb strings.Builder

	for _, err := range v.Errors {
		sb.WriteString("\n")
		sb.WriteString(err.Error())
	}

	return sb.String()
}

func (v VxResult) StringArray() []string {
	var errArray []string = []string{}

	for _, err := range v.Errors {
		errArray = append(errArray, err.Error())
	}

	return errArray
}

// Parses and validates all the field's values of the given struct `v` against
// the rules mentioned in the "vx" tag.
//
// Returns (res VxResult, ok bool) where `res` is the result with `res.Errors`
// containing all the errors that happened during the entire functional call
// and `ok` represents whether the `res.Errors` were generated due to something
// else other than the actual validation. Something went wrong before the
// validation, like the parsing of the struct or making `Tag`.
func ValidateStruct(v any) (res VxResult, ok bool) {
	ok = true
	res = VxResult{
		Errors: []error{},
	}

	parsedStruct, err := internal.ParseStruct(v)
	if err != nil {
		res.Errors = append(res.Errors, err)
		ok = false
		// This is the only case where we early return because if this fails
		// there is literally nothing more we can do after this.
		return res, ok
	}

	// FIXME: Making these maps might not be the best way to do this.
	fieldMap := map[string]internal.VxField{}
	tagMap := map[string]internal.VxTag{}

	// NOTE: do not return if `MakeTag` fails for a field. We want to collect
	// all the errors for all fields first and then return them at the end.
	for _, field := range parsedStruct.Fields {
		tag, err := internal.MakeTag(field)
		if err != nil {
			ok = false
			res.Errors = append(res.Errors, err)
		}

		tagMap[field.Name] = tag
		fieldMap[field.Name] = field
	}

	// We have an internal error, so we are not going to return them without
	// executing rules on the parsedStruct's `Fields`.
	if !ok {
		return res, ok
	}

	for _, field := range parsedStruct.Fields {
		// NOTE: `field.ValueType.Kind()` panics when `field.Value` is `nil` !!!
		if field.Value != nil {
			tag := tagMap[field.Name]

			if tag.Type != field.ValueType && field.Type.Kind() == reflect.Interface && tag.HasExplicitType && tag.Type.Kind() != reflect.Interface {
				switch tag.Type.Kind() {
				case reflect.Slice:
					var actualElemType reflect.Type = nil
					hasElems := false

					mySlice, ok := field.Value.([]any)
					if ok {
						for _, elem := range mySlice {
							hasElems = true

							if actualElemType == nil || reflect.TypeOf(elem) != tag.Type.Elem() {
								actualElemType = reflect.TypeOf(elem)
							}
						}
					}

					{
						if tag.Type.Elem().Kind() != field.ValueType.Elem().Kind() {
							if tag.Type.Elem().Kind() != reflect.Interface && field.ValueType.Elem().Kind() != reflect.Interface || (hasElems && tag.Type.Elem() != actualElemType) {
								elemType := field.ValueType.Elem()

								if hasElems {
									elemType = actualElemType
								}

								err = fmt.Errorf("%s should be an array of elem of type %s but got %s", field.Name, tag.Type.Elem(), elemType)
								res.Errors = append(res.Errors, err)
							}
						}
					}
				case reflect.Array:
					{
						var actualElemType reflect.Type = nil
						hasElems := false

						mySlice, ok := field.Value.([]any)
						if ok {
							for _, elem := range mySlice {
								hasElems = true

								if actualElemType == nil || reflect.TypeOf(elem) != tag.Type.Elem() {
									actualElemType = reflect.TypeOf(elem)
								}
							}
						}

						if tag.Type.Elem().Kind() != field.ValueType.Elem().Kind() {
							if tag.Type.Elem().Kind() != reflect.Interface && field.ValueType.Elem().Kind() != reflect.Interface || (hasElems && tag.Type.Elem() != actualElemType) {
								elemType := field.ValueType.Elem()

								if hasElems {
									elemType = actualElemType
								}

								err = fmt.Errorf("%s should be an array of elem of type %s but got %s", field.Name, tag.Type.Elem(), elemType)
								res.Errors = append(res.Errors, err)
							}
						}

						if tag.Type.Len() != field.ValueType.Len() {
							err = fmt.Errorf("%s should be an array of length %d but got %d", field.Name, tag.Type.Len(), field.ValueType.Len())
							res.Errors = append(res.Errors, err)
						}
					}
				case reflect.Map:
					{
						var actualKeyType, actualElemType reflect.Type = nil, nil
						hasElems := false

						myMap, ok := field.Value.(map[string]any)
						if ok {
							for key, elem := range myMap {
								hasElems = true

								// If `key` if of type `any` in `Field.Value` then `ValueType.Key()`
								// can be of more than 1 type. We want `actualKeyType` to change
								// only when it's nil (first time) and when the reflect.TypeOf(key)
								// doens't match to the key type from the tag.
								//
								// One thing to note is that the error message will show `actualKeyType`
								// which will be the last wrong type key/elem that we found.
								//
								// Same goes for `actualElemType`.
								if actualKeyType == nil || reflect.TypeOf(key) != tag.Type.Key() {
									actualKeyType = reflect.TypeOf(key)
								}

								if actualElemType == nil || reflect.TypeOf(elem) != tag.Type.Elem() {
									actualElemType = reflect.TypeOf(elem)
								}
							}
						}

						if tag.Type.Key().Kind() != field.ValueType.Key().Kind() {
							if tag.Type.Key().Kind() != reflect.Interface && field.ValueType.Key().Kind() != reflect.Interface || (hasElems && tag.Type.Key() != actualKeyType) {
								keyType := field.ValueType.Key()

								if hasElems {
									keyType = actualKeyType
								}

								err = fmt.Errorf("%s should be a map with key of type %s and elem of type %s but got map with key of type %s", field.Name, tag.Type.Key(), tag.Type.Elem(), keyType)
								res.Errors = append(res.Errors, err)
							}
						}

						if tag.Type.Elem().Kind() != field.ValueType.Elem().Kind() {
							if tag.Type.Elem().Kind() != reflect.Interface && field.ValueType.Elem().Kind() != reflect.Interface || (hasElems && tag.Type.Elem() != actualElemType) {
								elemType := field.ValueType.Elem()

								if hasElems {
									elemType = actualElemType
								}

								err = fmt.Errorf("%s should be a map with key of type %s and elem of type %s but got map with elem of type %s", field.Name, tag.Type.Key(), tag.Type.Elem(), elemType)
								res.Errors = append(res.Errors, err)
							}
						}
					}
				default:
					err = fmt.Errorf("%s should be of type %s but got %s", field.Name, tag.Type, field.ValueType)
					res.Errors = append(res.Errors, err)
				}

			}
		}

		// This is a bit tricky. Although we have a "required" Rule to check and warn
		// the user if a required value is not present or "empty" but if the developer
		// forgets to add the "required" Rule and then reads this field, that would
		// lead to a panic during runtime, which is NOT SO GOOD !!
		//
		// So in order to prevent the runtime panics for Vx users, we should set the
		// field.Value to be a default of the type that was being expected. If the
		// expected type is any, we will default it to empty string.
		if field.Value == nil {
			fmt.Printf("YIKES! This is bad, %s is nil which can be a nasty runtime error.", field.Name)
		}

	}

	for fieldName, tag := range tagMap {
		field := fieldMap[fieldName]

		for _, rule := range tag.Rules {
			err := rule.Exec(field)
			if err != nil {
				res.Errors = append(res.Errors, err)
			}
		}
	}

	// fmt.Println("\nParsedStruct:", parsedStruct)
	// fmt.Println("\nResult:", ok, res)

	return res, ok
}
