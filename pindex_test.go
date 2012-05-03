package pindex

import (
	"strings"
	"testing"
)

var (
	corpus = []Review{
		Review{
			Author: "Lindsay Zoladz",
			Body:   `Last year, a video made the rounds of St. Vincent covering Big Black's "Kerosene" with a thrashing, punk intensity unlike what we'd come to expect from Annie Clark.`,
		},
		Review{
			Author: "Lindsay Zoladz",
			Body:   `The gods of power pop are perennially unkind. Sure, it's never been music's most commercially viable subgenre, but I suspect an element of cosmic doom, too.`,
		},
		Review{
			Author: "Lindsay Zoladz",
			Body:   `With a frontman iconically haughty enough to be a feasible answer to a New York Times crossword puzzle clue (L.A.-by-way-of-East Village curmudgeon; 14 letters)-- it has been especially easy during this run of albums to dismiss the Magnetic Fields as Ivory Tower pop, wrapped up in the cleverness of their own ideas and out of touch with the world below.`,
		},
		Review{
			Author: "Ned Raggett",
			Body:   `The Leeds, England-based Hood, now some years into a relaxed hiatus following 2005's Outside Closer, are one of many bands that should have been deservedly famous.`,
		},
		Review{
			Author: "Ned Raggett",
			Body:   `Throughout all three albums, Opeth are about explicit formalism as stirring power via the rock gods-- the goal is far from new, but it's done so expertly that it's hard not to be impressed.`,
		},
		Review{
			Author: "Ned Raggett",
			Body:   `Everything said is said with a sense of loss, and everything you hear on Where the Sands Turn to Gold is something that is heard two ways, as the expression of intense feeling thrillingly captured and as the mark of personal destruction.`,
		},
		Review{
			Author: "Mark Richardson",
			Body:   `"I hear a lot of music that's just lazy-- you know, people in their bedrooms singing some shit into the microphone." That's California singer and songwriter Julia Holter, talking to Pitchfork recently.`,
		},
		Review{
			Author: "Mark Richardson",
			Body:   `Beal's debut album, Acousmatic Sorcery, which consists of of home-recorded songs stretching back over the last few years, doesn't answer this question. But it does suggest that the answer, when it finally comes, may well be fascinating.`,
		},
		Review{
			Author: "Mark Richardson",
			Body:   `Electronic music was once the domain of academics and researchers with access to vast rooms filled with pulsing tubes and clusters of snaking cables. Only those with a commission were allowed anywhere near the machinery.`,
		},
	}
)

func TestSimpleCount(t *testing.T) {
	idx := NewIndex()
	for _, r := range corpus {
		idx.Add(r)
	}
	countMusic := func(r Review) int {
		return strings.Count(strings.ToLower(r.Body), "music")
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

func TestLoadFile(t *testing.T) {
	idx := NewIndex()
	if err := idx.LoadFile("first-batch.json"); err != nil {
		t.Fatalf("LoadFile failed: %s", err)
	}
	if idx.Authors() != 68 {
		t.Fatalf("expected %d authors, got %d", 68, idx.Authors())
	}
	wordCount := func(r Review) int {
		return strings.Count(strings.ToLower(r.Body), " ") + 1
	}
	r := SortedResults(idx.MapAverage(wordCount))
	if len(r) != idx.Authors() {
		t.Fatalf("expected %d results, got %d", idx.Authors(), len(r))
	}
	first, second, last := "Drew Daniel", "Hank Shteamer", "Sam Hockley-Smith"
	if r[0].Author != first {
		t.Errorf("expected #1 to be %s, got '%s'", first, r[0].Author)
	}
	if r[1].Author != second {
		t.Errorf("expected #2 to be %s, got '%s'", second, r[1].Author)
	}
	if r[len(r)-1].Author != last {
		t.Errorf("expected last-place to be %s, got '%s'", last, r[len(r)-1].Author)
	}
	for i, pair := range r {
		t.Logf(
			"%2d/%2d) %s (average %d words per review)",
			i+1,
			len(r),
			pair.Author,
			pair.Score,
		)
	}
}
