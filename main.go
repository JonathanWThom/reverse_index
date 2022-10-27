package main

import (
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
