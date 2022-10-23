package main

import (
	"fmt"
	"strings"
)

// Note to self: Not working correctly when matching on mutliple fields

func main() {
	s := newStore()

	d1 := document{
		"name":    "Jonathan",
		"city":    "Ellensburg",
		"tagline": "Here's to feeling good all the time",
	}
	fmt.Printf("Adding document: %#v\n", d1)
	s.add(d1)

	d2 := document{
		"name":    "Jerry",
		"city":    "New York",
		"tagline": "What's the deal with time",
	}
	fmt.Printf("Adding document: %#v\n", d2)
	s.add(d2)

	d3 := document{
		"name":    "Kramer",
		"city":    "New York",
		"tagline": "These pretzels are making me thirsty",
	}
	fmt.Printf("Adding document: %#v\n", d3)
	s.add(d3)

	d4 := document{
		"name":    "Marvin Martian",
		"city":    "Mars",
		"tagline": "Love 2 mess with Bugs on Mars",
	}
	fmt.Printf("Adding document: %#v\n\n", d4)
	s.add(d4)

	fmt.Println("Searching for term 'good'...")
	result1 := s.search("good")
	fmt.Printf("%s\n\n", result1)

	fmt.Println("Searching for term 'new'...")
	result2 := s.search("new")
	fmt.Printf("%s\n\n", result2)

	fmt.Println("Searching for term 'Time'...")
	result3 := s.search("Time")
	fmt.Printf("%s\n\n", result3)

	fmt.Println("Searching for term 'mars'...")
	result4 := s.search("mars")
	fmt.Printf("%s\n", result4)
}

func newStore() store {
	return store{
		documents: map[int]document{},
		index:     index{terms: map[string]documentRefs{}},
	}
}

type store struct {
	documents map[int]document
	index     index
}

func (s *store) search(term string) []searchResult {
	res := []searchResult{}
	for _, ref := range s.index.terms[strings.ToLower(term)] {
		doc := s.documents[ref.id]
		res = append(res, searchResult{document: doc, fields: ref.fields})
	}

	return res
}

func (s *store) add(doc document) {
	id := len(s.documents)
	s.documents[id] = doc

	s.index.addDoc(doc, id)
}

type document map[string]string // could change to interface{}

type documentRef struct {
	id     int
	fields []string
}

type documentRefs []documentRef

func (d documentRefs) findRefForDoc(id int) *documentRef {
	for _, ref := range d {
		if ref.id == id {
			return &ref
		}
	}

	return nil
}

type index struct {
	terms map[string]documentRefs
}

func (i *index) addDoc(doc document, id int) {
	for field, value := range doc {
		for _, term := range strings.Split(value, " ") {
			t := strings.ToLower(term)
			refs := i.terms[t]
			ref := refs.findRefForDoc(id)

			if ref != nil {
				ref.fields = append(ref.fields, field)
			} else {
				newRef := documentRef{id: id, fields: []string{field}}
				i.terms[t] = append(refs, newRef)
			}
		}
	}
}

type searchResult struct {
	document document
	fields   []string
}
