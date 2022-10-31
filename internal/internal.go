package internal

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

const (
	// Tag key on struct fields for VX specific data.
	//
	// Ex:
	// type myStruct struct {
	//	   name: string `vx:"some value in the tag."`
	// }
	VX_TAG_KEY = "vx"
)

func makeType(s string) (reflect.Type, error) {
	var typ reflect.Type

	if s == "any" || s == "interface{}" {
		// Idk how to create a reflect.Type for interface{} except this
		// making a struct with a field of type inteface{} and then use
		// the field's type.
		type a struct {
			A any
		}

		toParse := a{A: "any value here will work ;)"}
		fields, err := ParseStruct(toParse)
		if err != nil {
			return typ, fmt.Errorf("failed to make type for %s", s)
		}

		typ = fields.Fields[0].Type
	} else if s == "bool" {
		typ = reflect.TypeOf(bool(true))
	} else if s == "int" {
		typ = reflect.TypeOf(int(0))
	} else if s == "float64" {
		typ = reflect.TypeOf(float64(0))
	} else if s == "string" {
		typ = reflect.TypeOf(string(""))
	} else if strings.HasPrefix(s, "map") {
		// Something like map[any]any
		leftIdx := strings.Index(s, "[")
		rightIdx := strings.Index(s, "]")

		if leftIdx == -1 || rightIdx == -1 {
			return typ, fmt.Errorf("invalid type in tag '%s'", s)
		}

		keyStr := s[leftIdx+1 : rightIdx]

		keyType, err := makeType(keyStr)
		if err != nil {
			return typ, fmt.Errorf("couldn't make a type for the key '%s' of the map '%s'. %s", keyStr, s, err.Error())
		}

		elemStr := s[rightIdx+1:]

		elemType, err := makeType(elemStr)
		if err != nil {
			return typ, fmt.Errorf("couldn't make a type for the elem '%s' of the map '%s'. %s", elemStr, s, err.Error())
		}

		typ = reflect.MapOf(keyType, elemType)
	} else if strings.HasPrefix(s, "[]") {
		// Something like []any
		sliceTypeStr := s[2:]

		sliceType, err := makeType(sliceTypeStr)
		if err != nil {
			return typ, fmt.Errorf("couldn't make a type for the elem '%s' of the slice '%s'. %s", sliceTypeStr, s, err.Error())
		}

		typ = reflect.SliceOf(sliceType)
	} else if !strings.HasPrefix(s, "[]") && strings.HasPrefix(s, "[") {
		// Something like [10]any
		leftIdx := strings.Index(s, "[")
		rightIdx := strings.Index(s, "]")

		if leftIdx == -1 || rightIdx == -1 {
			return typ, fmt.Errorf("invalid type in tag '%s'", s)
		}

		lenStr := s[leftIdx+1 : rightIdx]

		arrayLen, err := strconv.Atoi(lenStr)
		if err != nil {
			return typ, fmt.Errorf("got invalid length of the array '%s'", lenStr)
		}

		elemStr := s[rightIdx+1:]

		elemType, err := makeType(elemStr)
		if err != nil {
			return typ, fmt.Errorf("couldn't make a type for the elem '%s' of the array '%s'. %s", elemStr, s, err.Error())
		}

		typ = reflect.ArrayOf(arrayLen, elemType)
	} else {
		return typ, fmt.Errorf("cannot make a type for '%s'. it's either invalid or unsupported", s)
	}

	return typ, nil
}

type VxTag struct {
	Type            reflect.Type
	HasExplicitType bool
	Rules           []rule
}

func MakeTag(field VxField) (VxTag, error) {
	tag := VxTag{
		Type:            reflect.TypeOf(nil),
		HasExplicitType: false,
		Rules:           []rule{},
	}

	splits := strings.Split(field.Tag, ",")

	// Looping first time to just get the "type".
	// PERFORMANCE: technicallly the time complexity remains O(n) even if we
	// loop twice over `splits` but maybe consider not looping twice?
	for _, split := range splits {
		if strings.Contains(split, "type=") {
			typeStr := strings.Split(split, "=")[1]

			tagType, err := makeType(typeStr)
			if err != nil {
				err := fmt.Errorf("%s: %s", field.Name, err.Error())
				return tag, err
			}

			tag.Type = tagType
			tag.HasExplicitType = true
		}
	}

	if !tag.HasExplicitType {
		tag.Type = field.Type
	}

	if tag.Type != field.Type && field.Type.Kind() != reflect.Interface {
		err := fmt.Errorf("type mismatch: %s type in struct is '%s' and in tag is '%s'", field.Name, field.Type, tag.Type)
		return tag, err
	}

	// Looping second time to build rules.
	for _, split := range splits {
		if strings.Contains(split, "type=") || strings.Contains(split, "name=") {
			// We have already handled these.
		} else if strings.Contains(split, "minLength=") {
			v := strings.Split(split, "=")[1]

			i, err := strconv.Atoi(v)
			if err != nil {
				return tag, fmt.Errorf("minLength: should be an integer, got %s", v)
			}

			if i <= 0 {
				return tag, fmt.Errorf("minLength: should be greater than 0, got %s", v)
			}

			rule := makeMinLength(i)
			tag.Rules = append(tag.Rules, rule)
		} else if strings.Contains(split, "required") {
			rule := makeRequired()
			tag.Rules = append(tag.Rules, rule)
		} else {
			log.Println("got an invalid value in the tag: ", split)
		}
	}

	return tag, nil
}

type VxField struct {
	// Name of the field.
	Name string
	// Type of the field.
	Type reflect.Type
	// The `"vx"` tag on the field.
	Tag string
	// Value of this field.
	Value any
	// Type of the value given to this field.
	//
	// VxStruct.ValueType will be different than VxField.Type when:
	// - VxField.Type is reflect.Interface then VxStruct.ValueType will be the
	//   type of the value that is actually passed to this field.
	// - If the struct was created during runtime via something like json.Encode, etc.
	ValueType reflect.Type
}

type VxStruct struct {
	// Name of the struct.
	Name string
	// Fields in the struct.
	Fields []VxField
}

func ParseStruct(toParse interface{}) (VxStruct, error) {
	// Could be any underlying type. DO NOT call `.Elem()` on it, might panic.
	val := reflect.ValueOf(toParse)

	// If its a pointer, resolve its value.
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	// Double check now that we have a struct (could still be anything).
	if val.Kind() != reflect.Struct {
		msg := fmt.Sprintf("expected struct, received %s", val.Kind().String())
		return VxStruct{}, errors.New(msg)
	}

	valType := val.Type()
	fields := []VxField{}

	for i := 0; i < valType.NumField(); i++ {
		field := valType.Field(i)

		Name := field.Name
		Type := field.Type
		Tag := field.Tag.Get(VX_TAG_KEY)

		splits := strings.Split(Tag, ",")
		for _, split := range splits {
			if strings.Contains(split, "name=") {
				v := strings.Split(split, "=")[1]

				if len(v) > 0 {
					// We have a `name` property on the tag, so lets use it.
					Name = v
				}
			}
		}

		Value := reflect.Indirect(reflect.ValueOf(toParse)).FieldByName(field.Name).Interface()
		ValueType := reflect.TypeOf(Value)

		// Checking if we have a Type on a Field which is another custom `type`(Go keyword).
		if strings.Contains(Type.String(), "main") && strings.Contains(ValueType.String(), "main") {
			fmt.Println("Field with a type that is of another type:", Type)

			// Switching over the actual `type`.
			switch Type.Kind() {
			case reflect.Map:
				{
					fmt.Println(Type, Type.Key(), Type.Elem())
				}

			default:
				{
					// Nothing to do here.
				}
			}
		}

		fields = append(fields, VxField{Name, Type, Tag, Value, ValueType})
	}

	return VxStruct{
		Name:   val.Type().Name(),
		Fields: fields,
	}, nil
}

// This interface should be implemented by everything but "type" in the "vx" tag.
type rule interface {
	Exec(field VxField) error
}

//
// Rules allowed for all types.
//

type required struct{}

func makeRequired() required {
	return required{}
}

func (r required) Exec(field VxField) error {
	if field.Value == nil || field.Value == "" {
		return fmt.Errorf("%s is required", field.Name)
	}

	return nil
}

//
// Rules allowed only for string.
//

type minLength struct {
	value int
}

func makeMinLength(l int) minLength {
	return minLength{l}
}

func (r minLength) Exec(field VxField) error {
	if field.Value == nil {
		return nil
	}

	wrongTypeErr := fmt.Errorf("%s - minLength: rule can only be applied to type string but was applied to type %s", field.Name, field.ValueType)

	if field.Type.Kind() != reflect.String && field.Type.Kind() != reflect.Interface {
		return wrongTypeErr
	}

	s, ok := field.Value.(string)
	if !ok {
		return wrongTypeErr
	}

	if len(s) < r.value {
		return fmt.Errorf("%s should have a minimum length of %v but has %v", field.Name, r.value, len(s))
	}

	return nil
}
