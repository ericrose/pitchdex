package main

import (
	"encoding/json"
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
	type AuthorScores struct {
		Author string         `json:"author"`
		Scores map[string]int `json:"scores"`
	}
	type AuthorScoresArray []AuthorScores

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	ac := r.AuthorCount()
	i, asa := 0, make(AuthorScoresArray, len(ac))
	for author, _ := range ac {
		as := AuthorScores{
			Author: author,
			Scores: map[string]int{},
		}

		ids := r.By(func(rv Review) bool { return rv.Author == author })
		for indexName, _ := range IndexDefinitions {
			as.Scores[indexName] = r.AverageScore(ids, indexName)
		}

		asa[i] = as
		i++
	}

	buf, err := json.Marshal(asa)
	if err != nil {
		return err
	}
	_, err = f.Write(buf)
	if err != nil {
		return err
	}

	return nil
}
