package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"vx"
)

type AssocFilter map[string]interface{}

type BonusFilter map[string]interface{}

type AssocOrBonusFilter map[string]interface{}

type user struct {
	Name                     any `vx:"name=name, type=string, required, minLength=3"`
	Age                      any `vx:"name=age, type=float64, required"`
	Location                 any `vx:"name=location, type=[]string"`
	AssocOrBonus             any `vx:"type=map[string]string"`
	HowWeWannaDoAssocOrBonus AssocOrBonusFilter
}

func test(w http.ResponseWriter, req *http.Request) {
	var u user

	err := json.NewDecoder(req.Body).Decode(&u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, ok := vx.ValidateStruct(u)
	fmt.Println("ok:", ok)
	fmt.Println("res:", res)

	w.Header().Set("Content-Type", "application/json")
	var jsonRes []byte = []byte{}

	if !ok {
		jsonRes, err = json.Marshal(res.StringArray())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonRes)
		return
	}

	if len(res.Errors) > 0 {
		jsonRes, err = json.Marshal(res.StringArray())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonRes)
		return
	}

	// fmt.Fprintf(w, "User: %+v", u)

	jsonRes, err = json.Marshal(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonRes)
}

func main() {
	http.HandleFunc("/test", test)
	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
