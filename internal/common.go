package internal

import (
	"fmt"
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
	TYPE_UNKNOWN     VxType = "vx_unknown"     // We do not understand the underlying type.
	TYPE_UNSUPPORTED VxType = "vx_unsupported" // We understand the underlying type but don't support it yet.
	TYPE_STRING      VxType = "vx_string"
)

func MakeVxType(s string) (VxType, error) {
	switch s {
	case "string":
		return TYPE_STRING, nil
	case "bool", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "float32", "float64", "complex64", "complex128", "int", "uint", "uintptr", "interface {}":
		{
			fmt.Println("MakeVxType: got unsupported type:", s)
			return TYPE_UNSUPPORTED, nil
		}
	default:
		return TYPE_UNKNOWN, fmt.Errorf("couldn't figure out VxType from the string: %s", s)
	}
}

// This interface should be implemented by everything but "type" in the "vx" tag.
type rule interface {
	Exec(v any) error
}
