package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	httpHost     *string = flag.String("host", "0.0.0.0", "HTTP host to bind to")
	httpPort     *int    = flag.Int("port", 8585, "HTTP port to bind to")
	metadataFile *string = flag.String("metadata", "metadata.json", "input file for review metadata")
	scoresFile   *string = flag.String("scores", "scores.json", "input file for scores data")
)

func main() {
	flag.Parse()

	reviewsBuf, err := ioutil.ReadFile(*metadataFile)
	if err != nil {
		log.Fatalf("%s", err)
	}
	scoresBuf, err := ioutil.ReadFile(*scoresFile)
	if err != nil {
		log.Fatalf("%s", err)
	}

	http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("./js/"))))
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./css/"))))
	http.Handle("/img/", http.StripPrefix("/img", http.FileServer(http.Dir("./img/"))))
	http.Handle("/ico/", http.StripPrefix("/ico", http.FileServer(http.Dir("./ico/"))))
	http.HandleFunc("/reviews", func(w http.ResponseWriter, r *http.Request) {
		w.Write(reviewsBuf)
	})
	http.HandleFunc("/scores", func(w http.ResponseWriter, r *http.Request) {
		w.Write(scoresBuf)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	endpoint := fmt.Sprintf("%s:%d", *httpHost, *httpPort)
	log.Fatalf("%s", http.ListenAndServe(endpoint, nil))
}
