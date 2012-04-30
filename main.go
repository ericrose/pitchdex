package pindex

import (
	"crypto/md5"
	"fmt"
)

type Review struct {
	Author string
	Body   []byte
}

func (r *Review) Hash() string {
	h := md5.New()
	h.Write([]byte(r.Author))
	h.Write(r.Body)
	return fmt.Sprintf("%x", h)
}

type Index struct {
	repr map[string][]Review
}

func NewIndex() *Index {
	return &Index{
		repr: map[string][]Review{},
	}
}

func (me *Index) Add(r Review) error {
	if reviews, ok := me.repr[r.Author]; ok {
		me.repr[r.Author] = append(reviews, r)
	} else {
		me.repr[r.Author] = []Review{r}
	}
	return nil
}

func (me *Index) MapAll(f func([]Review) int) map[string]int {
	scores := map[string]int{}
	for author, reviews := range me.repr {
		scores[author] = f(reviews)
	}
	return scores
}

func (me *Index) MapAdditive(f func(Review) int) map[string]int {
	return me.MapAll(
		func(reviews []Review) int {
			score := 0
			for _, review := range reviews {
				score += f(review)
			}
			return score
		},
	)
}
