package vx

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

// VX supports validation for the following types.
const (
	TYPE_BYTE   VxType = "byte"
	TYPE_FLOAT  VxType = "float"
	TYPE_INT    VxType = "int"
	TYPE_RUNE   VxType = "rune"
	TYPE_STRING VxType = "string"
)
