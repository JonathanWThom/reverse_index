package main

import (
	"encoding/json"
	"net/http"

	"github.com/jonathanwthom/reverse_index/store"
)

const (
	Get  = "GET"
	Post = "POST"
)

func add(w http.ResponseWriter, r *http.Request) {
	if r.Method != Post {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

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
	if r.Method != Get {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	res := s.Search(query)
	json.NewEncoder(w).Encode(res)
}
