package main

import (
	"flag"
	"github.com/peterbourgon/bonus"
	"log"
)

var (
	gobFile      *string = flag.String("gob", "data.gob", ".gob index for persistence")
	jsonFile     *string = flag.String("json", "", ".json to load to index (optional)")
	dictFile     *string = flag.String("dict", "/usr/share/dict/words", "dict file")
	rescore      *bool   = flag.Bool("rescore", false, "force rescoring of all articles")
	metadataFile *string = flag.String("metadata", "metadata.json", "output file for review metadata")
	scoresFile   *string = flag.String("scores", "scores.json", "output file for scores data")
)

func main() {
	flag.Parse()
	log.Printf("Derpin' it up")

	// Load reviews to memory
	r := loadReviews(*gobFile, *jsonFile)
	defer func() {
		if err := r.Persist(*gobFile); err != nil {
			log.Printf("persist: %s", err)
		} else {
			log.Printf("persisted to %s", *gobFile)
		}
	}()
	log.Printf("%s loaded", bonus.Pluralize(len(*r), "review"))

	// Perform any necessary scoring
	log.Printf("Scoring...")
	n := 0
	for indexName, scoringFunc := range IndexDefinitions {
		for id, review := range *r {
			if _, ok := review.Scores[indexName]; *rescore || !ok {
				(*r)[id].Scores[indexName] = scoringFunc(review)
				n++
			}
		}
	}
	log.Printf("calculated %d scores", n)

	// Dump statistics
	printStats(r)
	log.Printf("writing metadata to %s...", *metadataFile)
	r.WriteMetadata(*metadataFile)
	log.Printf("writing scores to %s...", *scoresFile)
	r.WriteScores(*scoresFile)
}

func loadReviews(gob, json string) *Reviews {
	r := Reviews{}
	if err := r.ImportGob(*gobFile); err != nil {
		log.Printf("%s: %s", *gobFile, err)
	}
	if *jsonFile != "" {
		if err := r.ImportJSON(*jsonFile); err != nil {
			log.Printf("%s: %s", *jsonFile, err)
		}
	}
	return &r
}

var IndexDefinitions = IndexMap{
	"Pitchformulaity":       Pitchformulaity,
	"Naïve sentence length": NaïveSentenceLength,
	"Words invented":        InventedWordsFunc(*dictFile),
	"Character count":       CharacterCount,
	"Word count":            WordCount,
	"Word length":           WordLength,
}

func printStats(r *Reviews) {
	authorCount := r.AuthorCount()
	log.Printf(
		"%s, %s",
		bonus.Pluralize(len(*r), "review"),
		bonus.Pluralize(len(authorCount), "author"),
	)

	for indexName, _ := range IndexDefinitions {
		m := map[string]int{}
		log.Printf("%s:", indexName)
		for author, _ := range authorCount {
			matchAuthor := func(rv Review) bool { return rv.Author == author }
			m[author] = r.AverageScore(r.By(matchAuthor), indexName)
		}
		printTopN(m, 5)
	}
}

func printTopN(m map[string]int, n int) {
	sorted := sortMapByValue(m)
	for i := 0; i < n && i < len(sorted); i++ {
		log.Printf(
			" %d/%d) %s (%d)",
			i+1,
			len(m),
			sorted[i].Key,
			sorted[i].Value,
		)
	}
}
