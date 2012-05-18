package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

func GetDB(filename string) (*sql.DB, error) {
	return sql.Open("sqlite3", filename)
}

func Initialize(db *sql.DB) error {
	statements := []string{
		"CREATE TABLE reviews (id INT PRIMARY KEY, body TEXT)",
		"CREATE TABLE authors (name TEXT PRIMARY KEY)",
		"CREATE TABLE review_scores (review_id INT, name STRING, score INT)",
		"CREATE TABLE author_scores (author_name STRING, name STRING, score INT)",
		"CREATE TABLE authorship (review_id INT, author_name STRING)",
		"CREATE INDEX review_score_name ON review_scores (name)",
		"CREATE INDEX author_score_name ON author_scores (name)",
		"CREATE INDEX review_score_nsc ON review_scores (name, score)",
		"CREATE INDEX author_score_nsc ON author_scores (name, score)",
	}
	for _, statement := range statements {
		db.Exec(statement) // Best-effort is.. best.. effort.
	}
	return nil
}

func InsertReview(db *sql.DB, review Review) error {
	_, err := db.Exec(
		"INSERT INTO reviews VALUES (?, ?)",
		review.ID,
		review.Body,
	)
	if err != nil {
		return err
	}

	// best-effort is OK: author may exist
	db.Exec(
		"INSERT INTO authors VALUES (?)",
		review.Author,
	)

	_, err = db.Exec(
		"INSERT INTO authorship VALUES (?, ?)",
		review.ID,
		review.Author,
	)
	if err != nil {
		return err
	}
	scoreMap := map[int]map[string]int{
		review.ID: review.Scores,
	}
	if err := InsertReviewScores(db, scoreMap, true); err != nil {
		return err
	}
	return nil
}

func InsertReviews(db *sql.DB, reviews Reviews) error {
	for _, review := range reviews {
		if err := InsertReview(db, review); err != nil {
			return err
		}
	}
	return nil
}

func SelectBody(db *sql.DB, id int) (string, error) {
	row := db.QueryRow("SELECT body FROM reviews WHERE id = ?", id)
	if row == nil {
		return "", fmt.Errorf("fail to query")
	}
	var s string
	if err := row.Scan(&s); err != nil {
		return "", err
	}
	return s, nil
}

func SelectBodys(db *sql.DB, ids []int) (map[int]string, error) {
	m := map[int]string{}
	strs := make([]string, len(ids))
	for i, id := range ids {
		strs[i] = fmt.Sprintf("%d", id)
	}
	clause := strings.Join(strs, ",")
	rows, err := db.Query(
		fmt.Sprintf(
			"SELECT id, body FROM reviews WHERE id IN (%s)",
			clause,
		),
	)
	if err != nil {
		return m, err
	}
	for rows.Next() {
		var id int
		var body string
		if err := rows.Scan(&id, &body); err != nil {
			return m, fmt.Errorf("SELECT error: %s", err)
		}
		m[id] = body
	}
	return m, nil
}

func SelectReviews(db *sql.DB, ids []int) (Reviews, error) {
	reviews := Reviews{}
	strs := make([]string, len(ids))
	for i, id := range ids {
		strs[i] = fmt.Sprintf("%d", id)
	}
	clause := strings.Join(strs, ",")
	rows, err := db.Query(
		fmt.Sprintf(
			`SELECT r.id, a.name, r.body
			 FROM reviews r, authors a, authorship x
			 WHERE r.id IN (%s)
			 AND x.review_id == r.id
			 AND x.author_name == a.name
			`,
			clause,
		),
	)
	if err != nil {
		return reviews, err
	}
	for rows.Next() {
		var id int
		var author string
		var body string
		if err := rows.Scan(&id, &author, &body); err != nil {
			return reviews, fmt.Errorf("SELECT review error: %s", err)
		}
		reviews[id] = Review{
			ID:     id,
			Author: author,
			Body:   body,
			Scores: map[string]int{},
		}
	}
	rows, err = db.Query(
		fmt.Sprintf(
			`SELECT review_id, name, score
			 FROM review_scores
			 WHERE review_id IN (%s)
			`,
			clause,
		),
	)
	if err != nil {
		return reviews, err
	}
	for rows.Next() {
		var id int
		var scoreName string
		var scoreValue int
		if err := rows.Scan(&id, &scoreName, &scoreValue); err != nil {
			return reviews, fmt.Errorf("SELECT score error: %s", err)
		}
		reviews[id].Scores[scoreName] = scoreValue
	}
	return reviews, nil
}

func SelectAllReviews(db *sql.DB) (Reviews, error) {
	reviews := Reviews{}
	ids := []int{}
	rows, err := db.Query("SELECT id FROM reviews")
	if err != nil {
		return reviews, err
	}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return reviews, err
		}
		ids = append(ids, id)
	}
	return SelectReviews(db, ids)
}

func InsertReviewScores(db *sql.DB, scores map[int]map[string]int, overwrite bool) error {
	for reviewId, scoreMap := range scores {
		for scoreName, scoreValue := range scoreMap {
			if overwrite {
				db.Exec(
					"DELETE FROM review_scores WHERE review_id = ? AND name = ?",
					reviewId,
					scoreName,
				)
			}
			_, err := db.Exec(
				"INSERT INTO review_scores VALUES (?, ?, ?)",
				reviewId,
				scoreName,
				scoreValue,
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func InsertAuthorScores(db *sql.DB, scores map[string]map[string]int, overwrite bool) error {
	for authorName, scoreMap := range scores {
		for scoreName, scoreValue := range scoreMap {
			if overwrite {
				db.Exec(
					"DELETE FROM author_scores WHERE author_name = ? AND name = ?",
					authorName,
					scoreName,
				)
			}
			_, err := db.Exec(
				"INSERT INTO author_scores VALUES (?, ?, ?)",
				authorName,
				scoreName,
				scoreValue,
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
