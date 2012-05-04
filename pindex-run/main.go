package main

import (
	"flag"
	"log"
	"pindex"
)

var (
	dbfile   *string = flag.String("dbfile", "pindex.db", "database file")
	load     *string = flag.String("load", "", ".json file to load")
	dictfile *string = flag.String("dictfile", "/usr/share/dict/words", "dict file")
	rescore  *bool   = flag.Bool("rescore", false, "attempt to rescore reviews")
)

func main() {
	flag.Parse()
	d, err := pindex.NewDatabase(*dbfile)
	if err != nil {
		log.Fatalf("%s", err)
	}
	if *load != "" {
		d.Initialize()
		if err := d.LoadFile(*load); err != nil {
			log.Fatalf("%s", err)
		}
	}
	log.Printf(
		"%d authors, %d reviews represented",
		d.Authors(),
		d.Reviews(),
	)
	dict := pindex.NewDict(*dictfile)
	log.Printf("%s: %d words", *dictfile, dict.Count())

	m := pindex.IndexMap{
		"Pitchformulaity":       pindex.Pitchformulaity,
		"Naïve sentence length": pindex.NaïveSentenceLength,
		"Words invented":        pindex.InventedWordsFunc(*dictfile),
		"Character count":       pindex.CharacterCount,
	}
	if *rescore {
		d.ScoreExistingReviews(m)
	}

	for name, _ := range m {
		log.Printf("%s:", name)
		r := d.AverageRanking(name)
		printTopN(5, r)
	}
}

func printTopN(n int, results []pindex.AuthorScore) {
	for i, authorScore := range results {
		log.Printf(
			"%d/%d) %s (%d)",
			i+1,
			len(results),
			authorScore.Author,
			authorScore.Score,
		)
		if i >= (n - 1) {
			break
		}
	}
}
