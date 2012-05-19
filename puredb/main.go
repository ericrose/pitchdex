package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	dbFile      *string = flag.String("db", "pitchdex.db", "database file")
	jsonFile    *string = flag.String("import", "", "JSON file to import (optional)")
	overwrite   *bool   = flag.Bool("overwrite", false, "overwrite existing reviews with JSON data")
	rescore     *bool   = flag.Bool("rescore", false, "rescore everything")
	dictFile    *string = flag.String("dict", "/usr/share/dict/words", "dict file")
	authorsFile *string = flag.String("authors", "data/authors.json", "authors output file")
	serve       *bool   = flag.Bool("serve", true, "serve HTTP")
	httpHost    *string = flag.String("http-host", "0.0.0.0", "HTTP host")
	httpPort    *int    = flag.Int("http-port", 8585, "HTTP port")
)

func main() {
	flag.Parse()

	// LOADING
	db, err := GetDB(*dbFile)
	if err != nil {
		log.Fatalf("%s: %s", *dbFile, err)
	}
	if *jsonFile != "" {
		reviews, err := GetReviewsFromJSON(*jsonFile)
		if err != nil {
			log.Fatalf("%s: %s", *jsonFile, err)
		}
		log.Printf("%s: contained %d Reviews", *jsonFile, len(reviews))
		if *overwrite {
			log.Printf("making Inserts with overwrite")
		} else {
			log.Printf("making Inserts without overwrite")
		}
		if err := db.InsertReviews(reviews, *overwrite); err != nil {
			log.Fatalf("Insert: %s", err)
		}
	}
	log.Printf("loading done, moving on to scoring")

	// SCORING
	// For every score-name in our little internal list,
	scores, count := map[int]map[string]int{}, 0
	for _, scoreName := range ReviewScoreNames {
		//   reviewsToScore =
		//     rescore all ? all IDs : IDs that don't have that score yet
		reviews := Reviews{}
		var err error
		if *rescore {
			if reviews, err = db.SelectAllReviews(); err != nil {
				log.Fatalf("Select: %s", err)
			}
			log.Printf(
				"%s: rescoring all Reviews (%d)",
				scoreName,
				len(reviews),
			)
		} else {
			if reviews, err = db.SelectAllReviewsWithout(scoreName); err != nil {
				log.Fatalf("Select: %s", err)
			}
			log.Printf(
				"%s: rescoring Reviews without this score (%d)",
				scoreName,
				len(reviews),
			)
		}
		//   for every reviewToScore,
		//     load and score
		for id, review := range reviews {
			f := IndexDefinitions[scoreName]
			score := f(review)
			if _, ok := scores[id]; !ok {
				scores[id] = map[string]int{}
			}
			scores[id][scoreName] = score
			count++
		}
	}
	log.Printf("Inserting %d scores for %d Reviews", count, len(scores))
	db.InsertReviewScores(scores, true)

	// SERVING
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
