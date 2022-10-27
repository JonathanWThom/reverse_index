package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jonathanwthom/reverse_index/store"
)

var s store.Store

func main() {
	s = store.NewStore()

	http.HandleFunc("/add", add)
	http.HandleFunc("/search", search)

	log.Fatal(http.ListenAndServe(":3333", nil))
}

func add(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	document := store.Document{}
	for k, v := range r.Form {
		document[k] = v[0]
	}

	s.Add(document)
	w.WriteHeader(http.StatusCreated)
}

func search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	res := s.Search(query)
	json.NewEncoder(w).Encode(res)
}
