package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FairyDevicesRD/macomp"
	"github.com/gorilla/mux"
)

//DoMA parses
func DoMA(w http.ResponseWriter, r *http.Request) {
	results := []macomp.MaResult{}

	var text string
	var callback string

	vars := mux.Vars(r)
	if text = vars["text"]; len(text) == 0 {
		text = r.FormValue("text")
	}
	if callback = vars["callback"]; len(callback) == 0 {
		callback = r.FormValue("callback")
	}

	if len(text) > 0 {
		results = resource.Parse(text)
	} else {
		http.Error(w, "Invalid query", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if len(callback) != 0 {
		fmt.Fprintf(w, "%s(", callback)
		if b, err := json.Marshal(results); err == nil {
			w.Write(b)
		}
		fmt.Fprintf(w, ")")
	} else {
		encoder := json.NewEncoder(w)
		encoder.Encode(results)
	}
}
