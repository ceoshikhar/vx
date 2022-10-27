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

type Tag struct {
	Rules []rule
}

func MakeTag(field StructField) (Tag, error) {
	tag := Tag{
		Rules: []rule{},
	}

	splits := strings.Split(field.Tag, ",")

	// Looping second time to build rules.
	for _, split := range splits {
		if strings.Contains(split, "minLength") {
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
		}

		if strings.Contains(split, "required") {
			rule := makeRequired()
			tag.Rules = append(tag.Rules, rule)
		}
	}

	return tag, nil
}

type StructField struct {
	// Name of the field.
	Name string
	// Type of the field.
	Type reflect.Type
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
		Type := field.Type
		Tag := field.Tag.Get(VX_TAG_KEY)
		Value := reflect.Indirect(reflect.ValueOf(toParse)).FieldByName(field.Name).Interface()

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
// Rules allowed for all types.
//

type required struct{}

func makeRequired() required {
	return required{}
}

func (r required) Exec(field StructField) error {
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

func (r minLength) Exec(field StructField) error {
	wrongTypeErr := fmt.Errorf("%s - minLength: rule can be applied to type string or any but got %s", field.Name, TypeOf(field.Value))

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
