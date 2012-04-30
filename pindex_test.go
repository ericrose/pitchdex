package pindex

import (
	"strings"
	"testing"
)

var (
	corpus = []Review{
		Review{
			Author: "Lindsay Zoladz",
			Body:   []byte(`Last year, a video made the rounds of St. Vincent covering Big Black's "Kerosene" with a thrashing, punk intensity unlike what we'd come to expect from Annie Clark.`),
		},
		Review{
			Author: "Lindsay Zoladz",
			Body:   []byte(`The gods of power pop are perennially unkind. Sure, it's never been music's most commercially viable subgenre, but I suspect an element of cosmic doom, too.`),
		},
		Review{
			Author: "Lindsay Zoladz",
			Body:   []byte(`With a frontman iconically haughty enough to be a feasible answer to a New York Times crossword puzzle clue (L.A.-by-way-of-East Village curmudgeon; 14 letters)-- it has been especially easy during this run of albums to dismiss the Magnetic Fields as Ivory Tower pop, wrapped up in the cleverness of their own ideas and out of touch with the world below.`),
		},
		Review{
			Author: "Ned Raggett",
			Body:   []byte(`The Leeds, England-based Hood, now some years into a relaxed hiatus following 2005's Outside Closer, are one of many bands that should have been deservedly famous.`),
		},
		Review{
			Author: "Ned Raggett",
			Body:   []byte(`Throughout all three albums, Opeth are about explicit formalism as stirring power via the rock gods-- the goal is far from new, but it's done so expertly that it's hard not to be impressed.`),
		},
		Review{
			Author: "Ned Raggett",
			Body:   []byte(`Everything said is said with a sense of loss, and everything you hear on Where the Sands Turn to Gold is something that is heard two ways, as the expression of intense feeling thrillingly captured and as the mark of personal destruction.`),
		},
		Review{
			Author: "Mark Richardson",
			Body:   []byte(`"I hear a lot of music that's just lazy-- you know, people in their bedrooms singing some shit into the microphone." That's California singer and songwriter Julia Holter, talking to Pitchfork recently.`),
		},
		Review{
			Author: "Mark Richardson",
			Body:   []byte(`Beal's debut album, Acousmatic Sorcery, which consists of of home-recorded songs stretching back over the last few years, doesn't answer this question. But it does suggest that the answer, when it finally comes, may well be fascinating.`),
		},
		Review{
			Author: "Mark Richardson",
			Body:   []byte(`Electronic music was once the domain of academics and researchers with access to vast rooms filled with pulsing tubes and clusters of snaking cables. Only those with a commission were allowed anywhere near the machinery.`),
		},
	}
)

func TestSimpleCount(t *testing.T) {
	idx := NewIndex()
	for _, r := range corpus {
		idx.Add(r)
	}
	countMusic := func(r Review) int {
		return strings.Count(strings.ToLower(string(r.Body)), "music")
	}
	results := idx.MapAdditive(countMusic)
	expected := map[string]int{
		"Mark Richardson": 2,
		"Lindsay Zoladz":  1,
		"Ned Raggett":     0,
	}
	for author, score := range expected {
		actualScore, ok := results[author]
		if !ok {
			t.Fatalf("%s not represented in results", author)
		}
		t.Logf("%s: %d", author, actualScore)
		if score != actualScore {
			t.Errorf("%s: expected score %d, got %d", author, score, actualScore)
		}
	}
}
