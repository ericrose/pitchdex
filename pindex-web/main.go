package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	httpHost     *string = flag.String("host", "0.0.0.0", "HTTP host to bind to")
	httpPort     *int    = flag.Int("port", 8585, "HTTP port to bind to")
	metadataFile *string = flag.String("metadata", "metadata.json", "input file for review metadata")
	scoresFile   *string = flag.String("scores", "scores.json", "input file for scores data")
)

func main() {
	flag.Parse()

	reviewsBuf, err := loadReviews(*metadataFile)
	if err != nil {
		log.Fatalf("%s", err)
	}
	scoresBuf, err := loadScores(*scoresFile)
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

func loadReviews(filename string) ([]byte, error) {
	type ReviewMetadata struct {
		ID        int    `json:"id"`
		Permalink string `json:"permalink"`
	}
	type ReviewMetadatas []ReviewMetadata

	f, err := os.Open(filename)
	if err != nil {
		return []byte{}, nil
	}
	defer f.Close()

	m := ReviewMetadatas{}
	if err := json.NewDecoder(f).Decode(&m); err != nil {
		return []byte{}, err
	}

	type Data struct {
		Rows [][]string `json:"aaData"`
	}
	d := Data{[][]string{}}
	for _, md := range m {
		d.Rows = append(d.Rows, []string{
			fmt.Sprintf("%s", md.ID),
			md.Permalink,
		})
	}

	buf, err := json.Marshal(d)
	if err != nil {
		return []byte{}, err
	}
	return buf, nil
}

func loadScores(filename string) ([]byte, error) {
	type AuthorScores struct {
		Author string         `json:"author"`
		Scores map[string]int `json:"scores"`
	}
	type AuthorScoresArray []AuthorScores

	f, err := os.Open(filename)
	if err != nil {
		return []byte{}, nil
	}
	defer f.Close()

	m := AuthorScoresArray{}
	if err := json.NewDecoder(f).Decode(&m); err != nil {
		return []byte{}, err
	}

	type Data struct {
		Rows [][]string `json:"aaData"`
	}
	d := Data{[][]string{}}
	for _, as := range m {
		d.Rows = append(d.Rows, []string{
			as.Author,
			fmt.Sprintf("%d", as.Scores["Pitchformulaity"]),
			fmt.Sprintf("%d", as.Scores["Naïve sentence length"]),
			fmt.Sprintf("%d", as.Scores["Words invented"]),
			fmt.Sprintf("%d", as.Scores["Character count"]),
			fmt.Sprintf("%d", as.Scores["Word count"]),
			fmt.Sprintf("%d", as.Scores["Word length"]),
		})
	}

	buf, err := json.Marshal(d)
	if err != nil {
		return []byte{}, err
	}
	return buf, nil
}
