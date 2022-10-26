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
	TYPE_ANY         VxType = "interface {}"   // interface{}
	TYPE_INT         VxType = "int"
	TYPE_STRING      VxType = "string"
)

func MakeVxType(s string) VxType {
	switch s {
	case "int":
		return TYPE_INT
	case "string":
		return TYPE_STRING
	case "interface {}":
		return TYPE_ANY
	case "bool", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "float32", "float64", "complex64", "complex128", "float", "uint", "uintptr", "byte", "rune":
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

func MakeTag(field StructField) (Tag, error) {
	tag := Tag{
		Type:  TYPE_EMPTY,
		Rules: []rule{},
	}

	splits := strings.Split(field.Tag, ",")

	// Looping first time to just get the "type".
	// PERFORMANCE: technicallly the time complexity remains O(n) even if we
	// loop twice over `splits` but maybe consider not looping twice?
	for _, split := range splits {
		if strings.Contains(split, "type") {
			tag.Type = MakeVxType(strings.Split(split, "=")[1])
		}

	}

	// No explicit `type` was provided in the tag.
	if tag.Type == TYPE_EMPTY {
		tag.Type = field.Type
	}

	if tag.Type != field.Type && field.Type != TYPE_ANY {
		err := fmt.Errorf("type mismatch - field '%s' type in struct is '%s' and type in tag is '%s'", field.Name, field.Type, tag.Type)
		return tag, err
	}

	// Looping second time to build rules.
	for _, split := range splits {
		if strings.Contains(split, "minLength") {
			if tag.Type != TYPE_ANY && tag.Type != TYPE_STRING {
				return tag, fmt.Errorf("minLength: rule is applicable only to value of type %s and %s, but got type %s", TYPE_ANY, TYPE_STRING, tag.Type)
			}

			v := strings.Split(split, "=")[1]

			i, err := strconv.Atoi(v)
			if err != nil {
				return tag, fmt.Errorf("minLength: value provided to rule should be an integer, got %s", v)
			}

			if i <= 0 {
				return tag, fmt.Errorf("minLength: value provided to rule should be greater than 0, got %s", v)
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
	Value any
}

type VxStruct struct {
	// Name of the struct.
	Name string
	// Fields in the struct.
	Fields []StructField
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
	fields := []StructField{}

	for i := 0; i < valType.NumField(); i++ {
		field := valType.Field(i)

		Name := field.Name
		Type := MakeVxType(field.Type.String())
		Tag := field.Tag.Get(VX_TAG_KEY)
		Value := fmt.Sprintf("%v", reflect.Indirect(reflect.ValueOf(toParse)).FieldByName(field.Name))

		fields = append(fields, StructField{Name, Type, Tag, Value})
	}

	return VxStruct{
		Name:   val.Type().Name(),
		Fields: fields,
	}, nil
}

// This interface should be implemented by everything but "type" in the "vx" tag.
type rule interface {
	Exec(field StructField) error
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

func (m minLength) Exec(field StructField) error {
	s, ok := field.Value.(string)

	if !ok {
		return errors.New("minLength: rule was exec against a value that is not a string")
	}

	if len(s) < m.value {
		return fmt.Errorf("minimum length allowed is 3 but got %v", len(s))
	}

	return nil
}
