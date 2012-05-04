package main

import (
	"flag"
	"log"
	"pindex"
)

var (
	infile  *string = flag.String("infile", "in.json", ".json file to load reviews")
	outfile *string = flag.String("outfile", "out.json", ".json file to store scores")
)

func main() {
	flag.Parse()
	idx := pindex.NewIndex()
	if err := idx.LoadFile(*infile); err != nil {
		log.Fatalf("%s", err)
	}
	log.Printf("%d authors, %d reviews represented", idx.Authors(), idx.Reviews())

	log.Printf("Pitchformulaity:")
	fr := pindex.SortedResults(idx.MapAverage(pindex.Pitchformulaity))
	printTopN(5, fr)

	log.Printf("Average naïve sentence length:")
	lr := pindex.SortedResults(idx.MapAverage(pindex.NaïveSentenceLength))
	printTopN(5, lr)
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
