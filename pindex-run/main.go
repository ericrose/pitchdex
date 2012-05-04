package main

import (
	"flag"
	"log"
	"pindex"
)

var (
	infile   *string = flag.String("infile", "in.json", ".json file to load reviews")
	outfile  *string = flag.String("outfile", "out.json", ".json file to store scores")
	dictfile *string = flag.String("dictfile", "/usr/share/dict/words", "dict file")
)

func main() {
	flag.Parse()
	idx := pindex.NewIndex()
	if err := idx.LoadFile(*infile); err != nil {
		log.Fatalf("%s", err)
	}
	log.Printf(
		"%s: %d authors, %d reviews represented",
		*infile,
		idx.Authors(),
		idx.Reviews(),
	)
	dict := pindex.NewDict(*dictfile)
	log.Printf("%s: %d words", *dictfile, dict.Count())

	log.Printf("Pitchformulaity:")
	fr := pindex.SortedResults(idx.MapAverage(pindex.Pitchformulaity))
	printTopN(5, fr)

	log.Printf("Average naïve sentence length:")
	lr := pindex.SortedResults(idx.MapAverage(pindex.NaïveSentenceLength))
	printTopN(5, lr)

	log.Printf("Average number of invented words:")
	ir := pindex.SortedResults(idx.MapAverage(pindex.InventedWordsFunc(*dictfile)))
	printTopN(5, ir)
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
