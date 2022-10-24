package vx

import (
	"fmt"
	"vx/internal"
)

func ValidateStruct(v any) (ok bool, errs []error) {
	structFields, err := internal.ParseStruct(v)

	if err != nil {
		errs = append(errs, err)
		return false, errs
	}

	for _, field := range structFields {
		tag, err := internal.MakeTag(field.Tag)

		if tag.Type != field.Type {
			errs = append(errs, fmt.Errorf("type mismatch - field '%s' type in struct is '%s' and type in tag is '%s'", field.Name, field.Type, tag.Type))
			return false, errs
		}

		if err != nil {
			errs = append(errs, err)
			return false, errs
		}

		for _, rule := range tag.Rules {
			err := rule.Exec(field.Value)

			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	return true, errs
}
