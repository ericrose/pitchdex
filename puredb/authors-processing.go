package main

import (
	"net/http"
	"strings"
	"log"
	"encoding/json"
	"fmt"
)

// QueryAuthors takes a server-side-processing HTTP request
// and returns the JSON data that serves DataTables.
func QueryAuthors(db DB, r *http.Request) ([]byte, error) {
	p, err := GetProcessingParams(r)
	if err != nil {
		return []byte{}, fmt.Errorf("parsing request: %s", err)
	}
	authors, total, matching, err := db.SelectAuthorsBy(p)
	if err != nil {
		return []byte{}, fmt.Errorf("selecting IDs: %s", err)
	}
	response := AuthorsToResponse(authors, p.Echo, total, matching)
	buf, err := json.Marshal(response)
	if err != nil {
		return []byte{}, fmt.Errorf("marshaling response: %s", err)
	}
	return buf, nil
}

func (db DB) SelectAuthorsBy(p ProcessingParams) (authors []Author, total int, matching int, err error) {

	// we build the query from a set of clauses
	// we provide "SELECT id" or "SELECT Count(id)" as neccessary
	// then, fromClause + matchClause + orderClause + limitClause

	// Basic stuff
	total = db.SelectAuthorCount()
	fromClause := "FROM authors a, author_scores as"
	whereClause := "WHERE a.name = as.name"

	// matchClause and orderClause depend on what the sort column is
	var matchClause, orderClause string
	switch p.SortColumn {
	case "Overall Bullshit Score", "Pitchformulaity", "Naïve sentence length", "Word count", "Words invented":
		matchClause = fmt.Sprintf(
			"AND as.name = '%s'",
			p.SortColumn,
		)
		orderClause = fmt.Sprintf(
			"ORDER BY as.score %s",
			p.SortDirection,
		)
	case "Author":
		orderClause = fmt.Sprintf("ORDER BY a.name %s", p.SortDirection)
	case "Reviews":
		err = fmt.Errorf("TODO can't sort by reviews yet")
		return
	default:
		err = fmt.Errorf("'%s': invalid sort column", p.SortColumn)
		return
	}

	// Special search-related work
	if p.SearchTerm != "" {
		matchClause = strings.Join(
			[]string{
				matchClause,
				fmt.Sprintf(
					"AND a.name LIKE '%%%s%%'",
					p.SearchTerm,
				),
			},
			"\n",
		)
	}

	// Limit clause is straightforward
	limitClause := fmt.Sprintf("LIMIT %s OFFSET %s", p.Count, p.Offset)

	// First, calculate matching rows
	if p.SearchTerm == "" {
		matching = total
	} else {
		matchingQuery := strings.Join(
			[]string{
				"SELECT Count(a.name)",
				fromClause,
				whereClause,
				matchClause,
				// no order necessary
				// no limit
			},
			"\n",
		)
		if err = db.db.QueryRow(matchingQuery).Scan(&matching); err != nil {
			err = fmt.Errorf("calculating matching: %s", err)
			return
		}
	}
	log.Printf("%d matching", matching)

	// Then, fetch IDs
	idsQuery := strings.Join(
		[]string{
			"SELECT DISTINCT r.id",
			fromClause,
			whereClause,
			matchClause,
			orderClause,
			limitClause,
		},
		"\n",
	)
	log.Printf("IDs query:\n\n%s\n\n", idsQuery)
	rows, err := db.db.Query(idsQuery)
	if err != nil {
		err = fmt.Errorf("making IDs query: %s", err)
		return
	}
	ids = IDSlice{}
	for rows.Next() {
		var id int
		if err = rows.Scan(&id); err != nil {
			err = fmt.Errorf("scanning an ID: %s", err)
			return
		}
		ids = append(ids, id)
	}
	return
}
