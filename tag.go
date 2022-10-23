package vx

import (
	"strings"
)

type rule interface {
	exec() (ok bool, err []error)
}

type Tag struct {
	Type  VxType
	rules []rule
}

func MakeTagFromString(s string) Tag {
	splits := strings.Split(s, ",")

	// fmt.Println("SPLITS:")
	// PrettyPrint(splits)

	t := Tag{}

	for _, split := range splits {
		if strings.Contains(split, "type") {
			typeStr := strings.Split(split, "=")[1]

			switch typeStr {
			case "byte":
				t.Type = TYPE_BYTE
			case "float":
				t.Type = TYPE_FLOAT
			case "int":
				t.Type = TYPE_INT
			case "rune":
				t.Type = TYPE_RUNE
			case "string":
				t.Type = TYPE_STRING
			}
		}

	}

	return t
}
