package internal

import (
	"fmt"
	"strconv"
	"strings"
)

type Tag struct {
	Type  VxType
	Rules []rule
}

func MakeTag(s string) (Tag, error) {
	tag := Tag{
		Type:  TYPE_UNKNOWN,
		Rules: []rule{},
	}

	splits := strings.Split(s, ",")
	// fmt.Println("SPLITS:")
	// PrettyPrint(splits)

	for _, split := range splits {
		if strings.Contains(split, "type") {
			Type, err := MakeVxType(strings.Split(split, "=")[1])
			if err != nil {
				return tag, err
			}

			tag.Type = Type
		}

		if strings.Contains(split, "minLength") {
			if tag.Type != TYPE_STRING {
				return tag, fmt.Errorf("minLength: rule is applicable only to value of TYPE_STRING, but got type %s", tag.Type)
			}

			v := strings.Split(split, "=")[1]

			i, err := strconv.Atoi(v)
			if err != nil {
				return tag, fmt.Errorf("MakeTag: value provided to rule minLength should be a valid integer, got %s", v)
			}

			rule := makeMinLength(i)
			tag.Rules = append(tag.Rules, rule)
		}
	}

	return tag, nil
}
