package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func (r *Reviews) WriteMetadata(filename string) error {
	type ReviewMetadata struct {
		ID        int    `json:"id"`
		Permalink string `json:"permalink"`
	}
	type ReviewMetadatas []ReviewMetadata

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	i, mds := 0, make(ReviewMetadatas, len(*r))
	for id, review := range *r {
		mds[i] = ReviewMetadata{
			ID:        id,
			Permalink: review.Permalink,
		}
		i++
	}

	buf, err := json.Marshal(mds)
	if err != nil {
		return err
	}
	_, err = f.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

func (r *Reviews) WriteScores(filename string) error {
	type ScoresStructure struct {
		Authors []map[string]string `json:"aaData"`
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	ac := r.AuthorCount()
	i, ss := 0, ScoresStructure{
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
		ss.Authors[i] = m
		i++
	}

	buf, err := json.Marshal(ss)
	if err != nil {
		return err
	}
	_, err = f.Write(buf)
	if err != nil {
		return err
	}

	return nil
}
