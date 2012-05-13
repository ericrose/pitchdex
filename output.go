package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func WriteReviews(reviews Reviews, filename string) error {
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
		Reviews: make([]map[string]string, len(reviews)),
	}
	for id, review := range reviews {
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

func WriteAuthors(authors map[string]map[string]int, filename string) error {
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
	count := len(authors)
	i, as := 0, AuthorsStructure{make([]map[string]string, count)}
	for author, scores := range authors {
		as.Authors[i] = map[string]string{
			"Author": author,
		}
		for indexName, score := range scores {
			as.Authors[i][indexName] = fmt.Sprintf("%d", score)
		}
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
