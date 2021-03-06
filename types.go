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

type JSONReview struct {
	Author    string `json:"reviewers"`
	Body      string `json:"editorial"`
	Permalink string `json:"key"`
}

type JSONReviews []JSONReview

func (r Reviews) ImportJSON(filename string, reimport bool) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	jsonReviews := JSONReviews{}
	if err := json.NewDecoder(f).Decode(&jsonReviews); err != nil {
		return err
	}
	for _, jsonReview := range jsonReviews {
		id, err := strconv.ParseInt(strings.Split(jsonReview.Permalink, "-")[0], 10, 64)
		if err != nil {
			return fmt.Errorf("%s: %s", jsonReview.Permalink, err)
		}
		if jsonReview.Author == "" || jsonReview.Body == "" {
			return fmt.Errorf(
				"%s: author %dB, body %dB",
				jsonReview.Permalink,
				len(jsonReview.Author),
				len(jsonReview.Body),
			)
		}
		if _, ok := r[int(id)]; !ok || reimport {
			r[int(id)] = Review{
				ID:        int(id),
				Author:    jsonReview.Author,
				Body:      jsonReview.Body,
				Permalink: jsonReview.Permalink,
				Scores:    map[string]int{},
			}
		}
	}
	return nil
}

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

type IndexMap map[string]ScoringFunction
type ScoringFunction func(Review) int

//
//
//

type Dict map[string]struct{}

func NewDict(filename string) Dict {
	d := Dict{}
	d.Load(filename)
	return d
}

func (d Dict) Load(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		return
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
}

func (d Dict) Count() int {
	return len(d)
}

func (d Dict) Has(s string) bool {
	_, ok := d[s]
	return ok
}
