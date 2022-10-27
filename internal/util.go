package internal

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Pretty print non primitive types like struct, map, array, slice.
func PrettyPrint(v any) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}

func TypeOf(v any) string {
	return reflect.ValueOf(v).Kind().String()
}
