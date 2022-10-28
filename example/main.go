package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"vx"
)

type AssocFilter map[string]interface{}

type BonusFilter map[string]interface{}

type user struct {
	Name  string      `vx:"name=name, required"`
	Age   interface{} `vx:"name=age, type=int, required"`
	Assoc AssocFilter
	Bonus BonusFilter
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
	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
