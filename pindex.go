package pindex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/peterbourgon/exp-html"
	"io/ioutil"
	"sort"
	"strings"
)

type Review struct {
	Author string `json:"reviewers"`
	Body   string `json:"editorial"`
}

type BunchOfReviews map[string]Review

type Index struct {
	repr map[string][]Review
}

func NewIndex() *Index {
	return &Index{
		repr: map[string][]Review{},
	}
}

func (me *Index) LoadFile(filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("ReadFile: %s", err)
	}
	reviews := BunchOfReviews{}
	err = json.Unmarshal(buf, &reviews)
	if err != nil {
		return fmt.Errorf("Unmarshal: %s", err)
	}
	for _, review := range reviews {
		me.Add(review)
	}
	return nil
}

func (me *Index) Add(r Review) {
	r.Body = stripHTML(r.Body)
	if reviews, ok := me.repr[r.Author]; ok {
		me.repr[r.Author] = append(reviews, r)
	} else {
		me.repr[r.Author] = []Review{r}
	}
}

func (me *Index) Authors() int {
	return len(me.repr)
}

func (me *Index) Reviews() int {
	count := 0
	for _, reviews := range me.repr {
		count += len(reviews)
	}
	return count
}

func (me *Index) MapAll(f func([]Review) int) map[string]int {
	// Split work
	data := map[string]chan int{}
	for author, reviews := range me.repr {
		c := make(chan int)
		data[author] = c
		go func(c chan int, reviews []Review) {
			c <- f(reviews)
			close(c)
		}(c, reviews)
	}
	// Aggregate results
	scores := map[string]int{}
	for author, c := range data {
		scores[author] = <-c
	}
	return scores
}

func (me *Index) MapAdditive(f func(Review) int) map[string]int {
	return me.MapAll(
		func(reviews []Review) int {
			value := 0
			for _, review := range reviews {
				value += f(review)
			}
			return value
		},
	)
}

func (me *Index) MapAverage(f func(Review) int) map[string]int {
	return me.MapAll(
		func(reviews []Review) int {
			value := 0
			for _, review := range reviews {
				value += f(review)
			}
			return int(float64(value) / float64(len(reviews)))
		},
	)
}

func stripHTML(s string) string {
	z := html.NewTokenizer(bytes.NewBufferString(s))
	results := []string{}
	done := false
	for !done {
		tt := z.Next()
		switch tt {
		case html.TextToken:
			results = append(results, string(z.Text()))
		case html.ErrorToken:
			done = true
		}
	}
	return strings.Join(results, "")
}

type AuthorScore struct {
	Author string
	Score  int
}

type AuthorScoreList []AuthorScore

func (l AuthorScoreList) Len() int           { return len(l) }
func (l AuthorScoreList) Less(i, j int) bool { return l[i].Score > l[j].Score }
func (l AuthorScoreList) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

func SortedResults(results map[string]int) AuthorScoreList {
	i, l := 0, make(AuthorScoreList, len(results))
	for k, v := range results {
		l[i] = AuthorScore{k, v}
		i += 1
	}
	sort.Sort(l)
	return l
}
