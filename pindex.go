package pindex

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type Review struct {
	ID     int
	Author string `json:"reviewers"`
	Body   string `json:"editorial"`
}

type BunchOfReviews map[string]Review

//
//
//

type Database struct {
	db *sql.DB
}

func NewDatabase(filename string) (*Database, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	return &Database{db}, nil
}

func (d *Database) Initialize() {
	tables := []string{
		`CREATE TABLE reviews (
			id INT PRIMARY KEY,
			author VARCHAR(64),
			body TEXT
		)`,
		`CREATE INDEX reviews_author ON reviews(author)`,
		`CREATE TABLE scores (
			review_id,
			index_name VARCHAR(64),
			score INT
		)`,
		`CREATE INDEX scores_index ON scores(index_name)`,
		`CREATE INDEX scores_review ON scores(review_id)`,
	}
	for _, createTable := range tables {
		if _, err := d.db.Exec(createTable); err != nil {
			log.Printf("%s", err)
		}
	}
}

type IndexMap map[string]ScoringFunction
type ScoringFunction func(Review) int

func (d *Database) LoadFile(filename string) error {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("ReadFile: %s", err)
	}
	reviews := BunchOfReviews{}
	err = json.Unmarshal(buf, &reviews)
	if err != nil {
		return fmt.Errorf("Unmarshal: %s", err)
	}
	for permalink, review := range reviews {
		i, err := strconv.ParseInt(strings.Split(permalink, "-")[0], 10, 64)
		if err != nil {
			log.Printf("%s: %s", permalink, err)
			continue
		}
		review.ID = int(i)

		//log.Printf("Loading %d: %s (%dB)", review.ID, review.Author, len(review.Body))
		if err := d.AddReview(review); err != nil {
			return err
		}
	}
	return nil
}

func (d *Database) AddReview(r Review) error {
	_, err := d.db.Exec(
		"INSERT INTO reviews (id, author, body) VALUES (?, ?, ?)",
		r.ID,
		r.Author,
		r.Body,
	)
	return err
}

func (d *Database) AddScore(r Review, index string, score int) error {
	_, err := d.db.Exec(
		"INSERT INTO scores (review_id, index_name, score) VALUES (?, ?, ?)",
		r.ID,
		index,
		score,
	)
	return err
}

func (d *Database) ScoreExistingReviews(m IndexMap) error {
	rows, err := d.db.Query(`SELECT id, author, body FROM reviews`)
	if err != nil {
		return err
	}
	log.Printf("Score existing reviews: reading...")
	for rows.Next() {
		var r Review
		rows.Scan(&r.ID, &r.Author, &r.Body)
		for indexName, scoreFunc := range m {
			defer d.AddScore(r, indexName, scoreFunc(r))
		}
	}
	log.Printf("Score existing reviews: writing...")
	return nil
}

func (d *Database) Reviews() int {
	row := d.db.QueryRow(`SELECT Count(*) FROM reviews`)
	var count int
	row.Scan(&count)
	return count
}

func (d *Database) Authors() int {
	row := d.db.QueryRow(`SELECT Count(DISTINCT author) FROM reviews`)
	var count int
	row.Scan(&count)
	return count
}

func (d *Database) ReviewsBy(author string) int {
	row := d.db.QueryRow(`SELECT Count(*) FROM reviews WHERE author = ?`, author)
	var count int
	row.Scan(&count)
	return count
}

func (d *Database) TotalScore(author, index string) int {
	row := d.db.QueryRow(
		`SELECT SUM(score) FROM scores
		 WHERE index_name = ?
		 AND review_id IN (
		   SELECT id FROM reviews WHERE author = ?
		 )`,
		index,
		author,
	)
	var sum int
	row.Scan(&sum)
	return sum
}

func (d *Database) AverageScore(author, index string) int {
	row := d.db.QueryRow(
		`SELECT AVG(score) FROM scores
		 WHERE index_name = ?
		 AND review_id IN (
		   SELECT id FROM reviews WHERE author = ?
		 )`,
		index,
		author,
	)
	var avg float64
	row.Scan(&avg)

	return int(avg)
}

type AuthorScore struct {
	Author string
	Score  int
}

func (d *Database) AverageRanking(index string) []AuthorScore {
	authorScores := []AuthorScore{}
	rows, err := d.db.Query(
		`SELECT r.author, Avg(s.score) AS avg_score
		 FROM reviews r, scores s
		 WHERE s.review_id = r.id
		   AND s.index_name = ?
		 GROUP BY r.author
		 ORDER by avg_score DESC`,
		index,
	)
	if err != nil {
		return authorScores
	}
	for rows.Next() {
		as := AuthorScore{}
		var f float64
		rows.Scan(&as.Author, &f)
		as.Score = int(f)
		authorScores = append(authorScores, as)
	}
	return authorScores
}

//
//
//

type Dict map[string]struct{}

func NewDict(filename string) *Dict {
	d := &Dict{}
	d.Load(filename)
	return d
}

func (d *Dict) Load(filename string) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		(*d)[strings.TrimSpace(strings.ToLower(line))] = struct{}{}
	}
}

func (d *Dict) Count() int {
	return len(*d)
}

func (d *Dict) Has(s string) bool {
	_, ok := (*d)[s]
	return ok
}
