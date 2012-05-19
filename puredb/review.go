package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Review struct {
	ID        int
	Author    string
	Body      string
	Permalink string
	Scores    map[string]int
}

type Reviews map[int]Review

func GetReviewsFromJSON(filename string) (Reviews, error) {
	type JSONReview struct {
		ID        int
		Author    string `json:"reviewers"`
		Body      string `json:"editorial"`
		Permalink string `json:"key"`
	}
	type JSONReviews []JSONReview
	jsonReviews := JSONReviews{}

	// Read file
	f, err := os.Open(filename)
	if err != nil {
		return Reviews{}, err
	}
	defer f.Close()

	// Unmarshal
	if err := json.NewDecoder(f).Decode(&jsonReviews); err != nil {
		return Reviews{}, err
	}

	// Convert
	reviews := make(Reviews, len(jsonReviews))
	for i, r := range jsonReviews {
		id, err := strconv.ParseInt(strings.Split(r.Permalink, "-")[0], 10, 64)
		if err != nil {
			return reviews, fmt.Errorf("%s: %s", r.Permalink, err)
		}
		reviews[i] = Review{
			ID:        int(id),
			Author:    r.Author,
			Body:      r.Body,
			Permalink: r.Permalink,
			Scores:    map[string]int{},
		}
	}

	return reviews, nil
}

//
//
//

type Filter func(Review) bool
type IDSlice []int

func (r Reviews) By(f Filter) IDSlice {
	matching := IDSlice{}
	for id, review := range r {
		if f(review) {
			matching = append(matching, id)
		}
	}
	return matching
}

func (r Reviews) AuthorCount() map[string]int {
	m := map[string]int{}
	for _, review := range r {
		if _, ok := m[review.Author]; ok {
			m[review.Author]++
		} else {
			m[review.Author] = 1
		}
	}
	return m
}

func (r Reviews) TotalScore(ids IDSlice, indexName string) int {
	v := 0
	for _, id := range ids {
		v += r[id].Scores[indexName]
	}
	return v
}

func (r Reviews) AverageScore(ids IDSlice, indexName string) int {
	return int(float64(r.TotalScore(ids, indexName)) / float64(len(ids)))
}

//
//
//

// A data structure to hold a key/value pair.
type Pair struct {
	Key   string
	Value int
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type PairList []Pair

func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value > p[j].Value }

// A function to turn a map into a PairList, then sort and return it.
func sortMapByValue(m map[string]int) PairList {
	p := make(PairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = Pair{k, v}
		i++
	}
	sort.Sort(p)
	return p
}

//
//
//

type Dict map[string]struct{}

func NewDict(filename string) Dict {
	d := Dict{}
	f, err := os.Open(filename)
	if err != nil {
		return d
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		d[strings.TrimSpace(strings.ToLower(line))] = struct{}{}
	}
	return d
}

func (d Dict) Count() int {
	return len(d)
}

func (d Dict) Has(s string) bool {
	_, ok := d[s]
	return ok
}
