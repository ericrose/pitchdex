package pindex

import (
	"sync"
	"testing"
)

const testDatabase = "test.db"

var once sync.Once

func setup(t *testing.T) *Database {
	db, err := NewDatabase(testDatabase)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if db.Authors() == 0 {
		db.Initialize()
		if err := db.LoadFile("first-batch.json"); err != nil {
			t.Fatalf("%s", err)
		}
	}
	t.Logf("%d Authors, %d Reviews", db.Authors(), db.Reviews())
	m := IndexMap{
		"Simple word count": WordCount,
	}
	once.Do(func() { db.ScoreExistingReviews(m) })
	return db
}

func TestLoad(t *testing.T) {
	db := setup(t)
	if c := db.ReviewsBy("Drew Daniel"); c != 1 {
		t.Errorf("'Drew Daniel' had %d reviews, expected %d", c, 1)
	}
	if s := db.AverageScore("Drew Daniel", "Simple word count"); s != 3104 {
		t.Errorf("Drew's word count was %d, expected %d", s, 3104)
	}

	if c := db.ReviewsBy("Lindsay Zoladz"); c != 33 {
		t.Errorf("'Lindsay Zoladz' had %d reviews, expected %d", c, 33)
	}
	if s := db.AverageScore("Lindsay Zoladz", "Simple word count"); s != 747 {
		t.Errorf("Lindsay's word count was %d, expected %d", s, 747)
	}
}

func TestAverageRanking(t *testing.T) {
	db := setup(t)
	r := db.AverageRanking("Simple word count")
	if len(r) != db.Authors() {
		t.Fatalf("expected %d, got %d", db.Authors(), len(r))
	}
	expected := map[int]string{
		0:          "Drew Daniel",
		1:          "Hank Shteamer",
		len(r) - 1: "Sam Hockley-Smith",
	}
	for position, expectedAuthor := range expected {
		if r[position].Author != expectedAuthor {
			t.Errorf(
				"position %d: got '%s', expected '%s'",
				position,
				r[position].Author,
				expectedAuthor,
			)
		}
	}
}
