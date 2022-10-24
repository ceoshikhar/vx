package internal

import (
	"errors"
	"fmt"
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

type VxType string

const (
	TYPE_EMPTY       VxType = ""               // No "type" was explicity declared in the vx tag.
	TYPE_UNKNOWN     VxType = "vx_unknown"     // We do not understand the underlying type.
	TYPE_UNSUPPORTED VxType = "vx_unsupported" // We understand the underlying type but don't support it yet.
	TYPE_ANY         VxType = "vx_any"         // interface{}
	TYPE_STRING      VxType = "vx_string"
)

func MakeVxType(s string) VxType {
	switch s {
	case "string":
		return TYPE_STRING
	case "interface {}":
		return TYPE_ANY
	case "bool", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "float32", "float64", "complex64", "complex128", "float", "int", "uint", "uintptr", "byte", "rune":
		{
			fmt.Println("MakeVxType: got unsupported type:", s)
			return TYPE_UNSUPPORTED
		}
	default:
		{
			fmt.Println("couldn't figure out VxType from the string:", s)
			return TYPE_UNKNOWN
		}
	}
}

type Tag struct {
	Type  VxType
	Rules []rule
}

func MakeTag(f StructField) (Tag, error) {
	tag := Tag{
		Type:  TYPE_EMPTY,
		Rules: []rule{},
	}

	splits := strings.Split(f.Tag, ",")

	// Looping first time to just get the "type".
	// @PERFORMANCE: technicallly the time complexity remains O(n) even if we
	// loop twice over `splits` but maybe consider not looping twice?
	for _, split := range splits {
		if strings.Contains(split, "type") {
			tag.Type = MakeVxType(strings.Split(split, "=")[1])
		}

	}

	if tag.Type == TYPE_EMPTY {
		tag.Type = f.Type
	}

	if tag.Type != f.Type {
		err := fmt.Errorf("type mismatch - field '%s' type in struct is '%s' and type in tag is '%s'", f.Name, f.Type, tag.Type)
		return tag, err
	}

	// Looping second time to build rules.
	for _, split := range splits {
		if strings.Contains(split, "minLength") {
			if tag.Type != TYPE_ANY && tag.Type != TYPE_STRING {
				return tag, fmt.Errorf("minLength: rule is applicable only to value of TYPE_STRING, but got type %s", tag.Type)
			}

			v := strings.Split(split, "=")[1]

			i, err := strconv.Atoi(v)
			if err != nil {
				return tag, fmt.Errorf("MakeTag: value provided to rule minLength should be an integer, got %s", v)
			}

			if i <= 0 {
				return tag, fmt.Errorf("MakeTag: value provided to rule minLength should be greater than 0, got %s", v)
			}

			rule := makeMinLength(i)
			tag.Rules = append(tag.Rules, rule)
		}
	}

	return tag, nil
}

type StructField struct {
	// Name of the field.
	Name string
	// Type of the field.
	Type VxType
	// The `VX_TAG` tag on the field.
	Tag string
	// Value of this field.
	Value string
}

func ParseStruct(toParse interface{}) ([]StructField, error) {
	// Could be any underlying type. DO NOT call `.Elem()` on it, might panic.
	val := reflect.ValueOf(toParse)

	// If its a pointer, resolve its value.
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	// Double check now that we have a struct (could still be anything).
	if val.Kind() != reflect.Struct {
		msg := fmt.Sprintf("util.ParseStruct(): expected struct, received %s", val.Kind().String())
		return nil, errors.New(msg)
	}

	valType := val.Type()
	fields := []StructField{}

	for i := 0; i < valType.NumField(); i++ {
		field := valType.Field(i)

		Name := field.Name
		Type := MakeVxType(field.Type.String())
		Tag := field.Tag.Get(VX_TAG_KEY)
		Value := fmt.Sprintf("%v", reflect.Indirect(reflect.ValueOf(toParse)).FieldByName(field.Name))

		fields = append(fields, StructField{Name, Type, Tag, Value})
	}

	return fields, nil
}

// This interface should be implemented by everything but "type" in the "vx" tag.
type rule interface {
	Exec(v any) error
}

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
