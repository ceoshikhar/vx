package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"vx"
)

type user struct {
	Name string      `json:"name" vx:"name=name, required"`
	Age  interface{} `json:"age" vx:"name=age, type=int, required"`
}

func test(w http.ResponseWriter, req *http.Request) {
	var u user

	err := json.NewDecoder(req.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, ok := vx.ValidateStruct(u)

	if !ok {
		http.Error(w, res.String(), http.StatusInternalServerError)
		return
	}

	if len(res.Errors) > 0 {
		http.Error(w, res.String(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "User: %+v", u)
}

func main() {
	http.HandleFunc("/test", test)

	http.ListenAndServe(":8080", nil)
}
