package internal

import (
	"encoding/json"
	"fmt"
)

// Pretty print non primitive types like struct, map, array, slice.
func PrettyPrint(v interface{}) (err error) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err == nil {
		fmt.Println(string(b))
	}
	return
}
