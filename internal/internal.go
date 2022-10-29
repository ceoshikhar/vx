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

var typeToKind = map[string]reflect.Kind{
	"bool":           reflect.Bool,
	"int":            reflect.Int,
	"int8":           reflect.Int8,
	"int16":          reflect.Int16,
	"int32":          reflect.Int32,
	"int64":          reflect.Int64,
	"uint":           reflect.Uint,
	"uint8":          reflect.Uint8,
	"uint16":         reflect.Uint16,
	"uint32":         reflect.Uint32,
	"uint64":         reflect.Uint64,
	"uintptr":        reflect.Uintptr,
	"float32":        reflect.Float32,
	"float64":        reflect.Float64,
	"complex64":      reflect.Complex64,
	"complex128":     reflect.Complex128,
	"array":          reflect.Array,
	"chan":           reflect.Chan,
	"func":           reflect.Func,
	"interface":      reflect.Interface,
	"map":            reflect.Map,
	"ptr":            reflect.Ptr,
	"slice":          reflect.Slice,
	"string":         reflect.String,
	"struct":         reflect.Struct,
	"unsafe.Pointer": reflect.UnsafePointer,
}

func stringToKind(s string) reflect.Kind {
	kind, ok := typeToKind[s]

	if !ok {
		kind = reflect.Invalid
	}

	return kind
}

type VxTag struct {
	Kind  reflect.Kind
	Rules []rule
}

func MakeTag(field VxField) (VxTag, error) {
	tag := VxTag{
		Kind:  reflect.Invalid,
		Rules: []rule{},
	}

	splits := strings.Split(field.Tag, ",")

	// Looping first time to just get the "type".
	// PERFORMANCE: technicallly the time complexity remains O(n) even if we
	// loop twice over `splits` but maybe consider not looping twice?
	for _, split := range splits {
		if strings.Contains(split, "type") {
			typeStr := strings.Split(split, "=")[1]
			tagKind := stringToKind(typeStr)

			if tagKind == reflect.Invalid {
				err := fmt.Errorf("%s has an invalid/unsupported type '%s' in tag", field.Name, typeStr)
				return tag, err
			}

			tag.Kind = tagKind
		}

	}

	// No explicit `type` was provided in the tag.
	if tag.Kind == reflect.Invalid {
		tag.Kind = field.Type.Kind()
	}

	if tag.Kind != field.Type.Kind() && field.Type.Kind() != reflect.Interface {
		err := fmt.Errorf("type mismatch: %s type in struct is '%s' and in tag is '%s'", field.Name, field.Type, tag.Kind)
		return tag, err
	}

	// NOTE: `field.ValueType.Kind()` panics when `field.Value` is `nil` !!!
	if field.Value != nil {
		if tag.Kind != field.ValueType.Kind() && field.Type.Kind() != reflect.Interface {
			err := fmt.Errorf("%s should be of type %s but got %s", field.Name, tag.Kind, field.ValueType.Kind())
			return tag, err
		}
	}

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
			if strings.Contains(split, "name") {
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
			fmt.Println("Field with a Type that is of another type:", Type)

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
	wrongTypeErr := fmt.Errorf("%s - minLength: rule can only be applied to type string but was applied to type %s", field.Name, TypeOf(field.Value))

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
