package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	jsonFile    *string = flag.String("import", "", "JSON file to import (optional)")
	rescore     *bool   = flag.Bool("rescore", false, "rescore everything")
	dictFile    *string = flag.String("dict", "/usr/share/dict/words", "dict file")
	authorsFile *string = flag.String("authors", "data/authors.json", "authors output file")
	serve       *bool   = flag.Bool("serve", true, "serve HTTP")
	httpHost    *string = flag.String("http-host", "0.0.0.0", "HTTP host")
	httpPort    *int    = flag.Int("http-port", 8585, "HTTP port")
)

func main() {
	flag.Parse()

	// Load

	reviews := Reviews{}
	if *jsonFile != "" {
		log.Printf("importing %s", *jsonFile)
		if err := reviews.ImportJSON(*jsonFile); err != nil {
			log.Printf("%s", err)
		}
	}
	if len(reviews) <= 0 {
		log.Fatalf("no reviews loaded")
	}
	log.Printf("%d reviews loaded", len(reviews))

	// Calculate review-scores
	log.Printf("calculating regular scores...")
	count := 0
	for indexName, scoringFunc := range IndexDefinitions {
		for id, review := range reviews {
			if _, ok := review.Scores[indexName]; !ok || *rescore {
				reviews[id].Scores[indexName] = scoringFunc(review)
				count++
			}
		}
	}
	log.Printf("calculating %s...", BullshitScore)
	allStats := GatherAll(reviews)
	for id, review := range reviews {
		if _, ok := review.Scores[BullshitScore]; *rescore || !ok {
			reviews[id].Scores[BullshitScore] = calculateBullshit(review, allStats)
			count++
		}
	}
	log.Printf("calculated %d scores", count)

	// Calculate author-scores
	authors := map[string]map[string]int{} // author -> avg scores
	for author, count := range reviews.AuthorCount() {
		authors[author] = map[string]int{
			"Reviews": count,
		}
		ids := reviews.By(func(r Review) bool { return r.Author == author })
		// The Bullshit of a given review is impacted by global stats.
		// But the Bullshit of an author is independent of other authors,
		// in the first-order sense.
		supplementedIndexDefinitions := IndexDefinitions
		supplementedIndexDefinitions[BullshitScore] = nil
		for indexName, _ := range IndexDefinitions {
			total := 0 // total score for this author for indexName
			for _, id := range ids {
				total += reviews[id].Scores[indexName]
			}
			averageScore := int(float64(total) / float64(len(ids))) // simple
			authors[author][indexName] = averageScore
		}
	}

	// Write
	if err := WriteAuthors(authors, *authorsFile); err != nil {
		log.Fatalf("%s", err)
	}

	// Serve HTTP
	staticDirs := []string{"js", "css", "img", "ico", "data"}
	for _, d := range staticDirs {
		route := fmt.Sprintf("/%s/", d)
		strip := fmt.Sprintf("/%s", d)
		serve := fmt.Sprintf("./%s/", d)
		http.Handle(
			route,
			http.StripPrefix(
				strip,
				http.FileServer(http.Dir(serve)),
			),
		)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf(
			"serving client %s (via %s) -- %s",
			r.RemoteAddr,
			func() string {
				if r.Referer() == "" {
					return "direct"
				}
				return r.Referer()
			}(),
			r.RequestURI,
		)
		http.ServeFile(w, r, "index.html")
	})

	endpoint := fmt.Sprintf("%s:%d", *httpHost, *httpPort)
	log.Printf("serving on %s", endpoint)
	log.Fatalf("%s", http.ListenAndServe(endpoint, nil))
}
