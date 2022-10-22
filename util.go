package vx

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// Pretty print non primitive types like struct, map, array, slice.
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

type StructField struct {
	// Name of the field.
	Name string
	// Type of the field.
	Type string
	// The `VX_TAG` tag on the field.
	Tag string
	// Value of this field.
	Value string
}

func ParseStructFields(toParse interface{}) ([]StructField, error) {
	// Could be any underlying type. DO NOT call `.Elem()` on it, it will panic.
	val := reflect.ValueOf(toParse)

	// If its a pointer, resolve its value.
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}

	// Should double check we now have a struct (could still be anything).
	if val.Kind() != reflect.Struct {
		msg := fmt.Sprintf("util.ParseStructFields(): expected struct, received %s", val.Kind().String())
		return nil, errors.New(msg)
	}

	valType := val.Type()
	parsedData := []StructField{}

	for i := 0; i < valType.NumField(); i++ {
		field := valType.Field(i)

		Name := field.Name
		Type := field.Type.String()
		Tag := field.Tag.Get("vx")
		Value := reflect.Indirect(reflect.ValueOf(toParse)).FieldByName(field.Name).String()

		parsedData = append(parsedData, StructField{Name, Type, Tag, Value})
	}

	return parsedData, nil
}
