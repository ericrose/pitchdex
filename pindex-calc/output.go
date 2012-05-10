package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func (r *Reviews) WriteAuthors(filename string) error {
	// Create the file
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Build the in-memory structure
	type AuthorsStructure struct {
		Authors []map[string]string `json:"aaData"`
	}
	ac := r.AuthorCount()
	i, as := 0, AuthorsStructure{
		Authors: make([]map[string]string, len(ac)),
	}
	for author, _ := range ac {
		ids := r.By(func(rv Review) bool { return rv.Author == author })
		m := map[string]string{
			"Author":      author,
			BullshitScore: fmt.Sprintf("%d", r.AverageScore(ids, BullshitScore)), // TODO fix per-author
		}
		for indexName, _ := range IndexDefinitions {
			m[indexName] = fmt.Sprintf("%d", r.AverageScore(ids, indexName))
		}
		as.Authors[i] = m
		i++
	}

	// Dump the structure to the file
	buf, err := json.Marshal(as)
	if err != nil {
		return err
	}
	_, err = f.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (r *Reviews) WriteReviews(filename string) error {
	// Create the file
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// Build the in-memory structure
	type ReviewsStructure struct {
		Reviews []map[string]string `json:"aaData"`
	}
	i, rs := 0, ReviewsStructure{
		Reviews: make([]map[string]string, len(*r)),
	}
	for id, review := range *r {
		m := map[string]string{
			"ID":     fmt.Sprintf("%d", id),
			"Title":  review.Permalink,
			"Author": review.Author,
		}
		for indexName, score := range review.Scores {
			m[indexName] = fmt.Sprintf("%d", score)
		}
		rs.Reviews[i] = m
		i++
	}

	// Dump the structure to the file
	buf, err := json.Marshal(rs)
	if err != nil {
		return err
	}
	_, err = f.Write(buf)
	if err != nil {
		return err
	}
	return nil
}
