package store

import (
	"sort"
	"strings"
)

func NewStore() Store {
	return Store{
		documents: map[int]Document{},
		index:     index{terms: map[string]documentRefs{}},
	}
}

type Store struct {
	documents map[int]Document
	index     index
}

func (s *Store) Search(term string) []searchResult {
	res := []searchResult{}
	for _, ref := range s.index.terms[strings.ToLower(term)] {
		doc := s.documents[ref.id]
		res = append(res, searchResult{Document: doc, fields: ref.fields})
	}

	sort.Sort(ByFieldCount(res))

	return res
}

type ByFieldCount []searchResult

func (s ByFieldCount) Len() int           { return len(s) }
func (s ByFieldCount) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByFieldCount) Less(i, j int) bool { return len(s[i].fields) > len(s[j].fields) }

func (s *Store) Add(doc Document) {
	id := len(s.documents)
	s.documents[id] = doc
	s.index.addDoc(doc, id)
}

type Document map[string]interface{}

type documentRef struct {
	id     int
	fields []string
}

type documentRefs []documentRef

func (d documentRefs) findRefForDoc(id int) (*int, *documentRef) {
	for i, ref := range d {
		if ref.id == id {
			return &i, &ref
		}
	}

	return nil, nil
}

type index struct {
	terms map[string]documentRefs
}

func (i *index) addDoc(doc Document, id int) {
	for field, value := range doc {
		for _, term := range strings.Split(value.(string), " ") {
			t := strings.ToLower(term)
			refs := i.terms[t]
			docRefPos, ref := refs.findRefForDoc(id)

			if ref != nil {
				ref.fields = append(ref.fields, field)
				i.terms[t][*docRefPos] = *ref
			} else {
				newRef := documentRef{id: id, fields: []string{field}}
				i.terms[t] = append(refs, newRef)
			}
		}
	}
}

type searchResult struct {
	Document Document
	fields   []string
}
