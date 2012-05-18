package main

import (
	"os"
	"testing"
)

func TestInitialize(t *testing.T) {
	os.Remove("testing.db")
	db, err := GetDB("testing.db")
	if err != nil {
		t.Fatalf("%s", err)
	}
	if err := Initialize(db); err != nil {
		t.Fatalf("%s", err)
	}
}

func TestInsertSelect(t *testing.T) {
	// Initialize
	os.Remove("testing.db")
	db, err := GetDB("testing.db")
	if err != nil {
		t.Fatalf("%s", err)
	}
	if err := Initialize(db); err != nil {
		t.Fatalf("%s", err)
	}

	// Insert first
	r1 := Review{
		ID:        123,
		Body:      "This is the review body.",
		Author:    "Joe Reviewer",
		Permalink: "123-foo-bar",
		Scores:    map[string]int{"Foo": 7},
	}
	if err := InsertReview(db, r1); err != nil {
		t.Fatalf("%s", err)
	}

	// Select 1
	body, err := SelectBody(db, r1.ID)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if body != r1.Body {
		t.Errorf("got '%s', expected '%s'", body, r1.Body)
	}

	// Insert second
	r2 := Review{
		ID:        456,
		Body:      "Second body.",
		Author:    "Frank Reviewer",
		Permalink: "456-baz",
		Scores:    map[string]int{},
	}
	if err := InsertReview(db, r2); err != nil {
		t.Fatalf("%s", err)
	}

	// Select and check 2 Bodys
	m, err := SelectBodys(db, []int{123, 456, 789})
	if err != nil {
		t.Fatalf("%s", err)
	}
	body123, ok123 := m[123]
	if !ok123 {
		t.Errorf("Review 123 not Selected")
	}
	if body123 != r1.Body {
		t.Errorf("got '%s', expected '%s'", body123, r1.Body)
	}
	body456, ok456 := m[456]
	if !ok456 {
		t.Errorf("Review 456 not Selected")
	}
	if body456 != r2.Body {
		t.Errorf("got '%s', expected '%s'", body456, r2.Body)
	}
	if _, ok789 := m[789]; ok789 {
		t.Errorf("ID 789 improperly returned")
	}

	// Select and check 2 Reviews
	reviews, err := SelectReviews(db, []int{123, 456, 789})
	if err != nil {
		t.Fatalf("%s", err)
	}
	review123, ok123 := reviews[123]
	if !ok123 {
		t.Errorf("Review 123 not Selected")
	}
	if review123.Author != r1.Author {
		t.Errorf("got '%s', expected '%s'", review123.Author, r1.Author)
	}
	if review123.Scores["Foo"] != r1.Scores["Foo"] {
		t.Errorf("got %d, expected %d", review123.Scores["Foo"], r1.Scores["Foo"])
	} else {
		t.Logf("%v: Foo score was %d", review123, review123.Scores["Foo"])
	}
}

func TestScoring(t *testing.T) {
}
