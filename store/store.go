package store

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"
	"os"
	"sort"
	"strings"
)

func NewStore() Store {
	// handle file not existing
	// this is inelegant, probably could get io.Reader from the get go
	data, err := os.ReadFile("data")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Store{
				Documents: map[int]Document{},
				Index:     index{Terms: map[string]documentRefs{}},
			}
		} else {
			log.Fatal(err)
		}
	}

	r := bytes.NewReader(data)
	dec := gob.NewDecoder(r)
	var store Store
	err = dec.Decode(&store)
	if err != nil {
		log.Fatal("decode error:", err)
	}

	return store
}

type Store struct {
	Documents map[int]Document
	Index     index
}

func (s *Store) Search(term string) []searchResult {
	res := []searchResult{}
	for _, ref := range s.Index.Terms[strings.ToLower(term)] {
		doc := s.Documents[ref.Id]
		res = append(res, searchResult{Document: doc, fields: ref.Fields})
	}

	sort.Sort(ByFieldCount(res))

	return res
}

type ByFieldCount []searchResult

func (s ByFieldCount) Len() int           { return len(s) }
func (s ByFieldCount) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByFieldCount) Less(i, j int) bool { return len(s[i].fields) > len(s[j].fields) }

func (s *Store) Add(doc Document) {
	id := len(s.Documents)
	s.Documents[id] = doc
	s.Index.addDoc(doc, id)

	var data bytes.Buffer
	enc := gob.NewEncoder(&data)
	err := enc.Encode(s)
	if err != nil {
		log.Fatal("encode error:", err)
	}

	f, err := os.OpenFile("data", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write(data.Bytes()); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

type Document map[string]interface{}

type documentRef struct {
	Id     int
	Fields []string
}

type documentRefs []documentRef

func (d documentRefs) findRefForDoc(id int) (*int, *documentRef) {
	for i, ref := range d {
		if ref.Id == id {
			return &i, &ref
		}
	}

	return nil, nil
}

type index struct {
	Terms map[string]documentRefs
}

func (i *index) addDoc(doc Document, id int) {
	for field, value := range doc {
		for _, term := range strings.Split(value.(string), " ") {
			t := strings.ToLower(term)
			refs := i.Terms[t]
			docRefPos, ref := refs.findRefForDoc(id)

			if ref != nil {
				ref.Fields = append(ref.Fields, field)
				i.Terms[t][*docRefPos] = *ref
			} else {
				newRef := documentRef{Id: id, Fields: []string{field}}
				i.Terms[t] = append(refs, newRef)
			}
		}
	}
}

type searchResult struct {
	Document Document
	fields   []string
}
