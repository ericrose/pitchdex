package main

import (
	"net/http"
	"log"
	"encoding/json"
	"strings"
	"fmt"
)

// QueryReviews takes a server-side-processing HTTP request
// and returns the JSON data that serves DataTables.
func QueryReviews(db DB, r *http.Request) ([]byte, error) {
	p, err := GetProcessingParams(r)
	if err != nil {
		return []byte{}, fmt.Errorf("parsing request: %s", err)
	}
	ids, total, matching, err := db.SelectReviewsBy(p)
	if err != nil {
		return []byte{}, fmt.Errorf("selecting IDs: %s", err)
	}
	reviews, err := db.SelectReviews(ids)
	if err != nil {
		return []byte{}, fmt.Errorf("selecting reviews: %s", err)
	}
	response := ReviewsToResponse(ids, reviews, p.Echo, total, matching)
	buf, err := json.Marshal(response)
	if err != nil {
		return []byte{}, fmt.Errorf("marshaling response: %s", err)
	}
	return buf, nil
}

type ProcessingParams struct {
	Echo string
	SortColumn string
	SortDirection string
	SearchTerm string
	Offset string
	Count string

	rawSortColumn string
	columnNames map[string]string // int -> name
}

func GetProcessingParams(r *http.Request) (ProcessingParams, error) {
	p := ProcessingParams{
		columnNames: map[string]string{},
	}
	if err := r.ParseForm(); err != nil {
		return p, err
	}
	log.Printf(">> %s", r.URL.Query())
	for key, values := range r.URL.Query() {
		if len(values) != 1 {
			return p, fmt.Errorf("%s: %d values", key, len(values))
		}
		switch {
		case key == "sEcho":
			p.Echo = values[0]
		case key == "sSearch":
			p.SearchTerm = values[0]
		case key == "iDisplayStart":
			p.Offset = values[0]
		case key == "iDisplayLength":
			p.Count = values[0]
		default:
			if toks := strings.Split(key, "_"); len(toks) == 2 {
				if toks[0] == "mDataProp" {
					p.columnNames[toks[1]] = values[0]
				} else if toks[0] == "iSortCol" && toks[1] == "0" {
					p.rawSortColumn = values[0]
				} else if toks[0] == "sSortDir" && toks[1] == "0" {
					p.SortDirection = values[0]
				}
			}
		}
	}
	if p.Count == "" {
		p.Count = "10"
	}
	if p.Offset == "" {
		p.Offset = "0"
	}
	if p.SortDirection == "" {
		p.SortDirection = "desc"
	}
	if !(p.SortDirection == "asc" || p.SortDirection == "desc") {
		return p, fmt.Errorf("'%s': invalid sort direction", p.SortDirection)
	}
	if p.rawSortColumn == "" {
		return p, fmt.Errorf("no sort column specified")
	}
	sortCol, ok := p.columnNames[p.rawSortColumn]
	if !ok {
		return p, fmt.Errorf("unknown rawSortColumn %s", p.rawSortColumn)
	}
	p.SortColumn = sortCol

	// Whoo boy
	validate := []string{
		p.SearchTerm,
		p.SortColumn,
		p.SortDirection,
		p.Count,
		p.Offset,
	}
	for _, v := range validate {
		if strings.ContainsAny(v, ";\\'") {
			return p, fmt.Errorf("'%s': invalid input", v)
		}
	}

	return p, nil
}

func (db DB) SelectReviewsBy(p ProcessingParams) (ids IDSlice, total int, matching int, err error) {

	// we build the query from a set of clauses
	// we provide "SELECT id" or "SELECT Count(id)" as neccessary
	// then, fromClause + matchClause + orderClause + limitClause

	// Basic stuff
	total = db.SelectReviewCount()
	fromClause := "FROM reviews r, review_scores rs, authors a, authorship x"
	whereClause := "WHERE r.id = rs.review_id AND r.id = x.review_id AND a.name = x.author_name"

	// matchClause and orderClause depend on what the sort column is
	var matchClause, orderClause string
	switch p.SortColumn {
	case "Pitchformulaity", "Naïve sentence length", "Word count", "Words invented":
		matchClause = fmt.Sprintf(
			"AND rs.name = '%s'",
			p.SortColumn,
		)
		orderClause = fmt.Sprintf(
			"ORDER BY rs.score %s",
			p.SortDirection,
		)
	case "Author":
		orderClause = fmt.Sprintf("ORDER BY a.name %s", p.SortDirection)
	case "ID":
		orderClause = fmt.Sprintf("ORDER BY r.id %s", p.SortDirection)
	case "Title":
		orderClause = fmt.Sprintf("ORDER BY r.title %s", p.SortDirection)
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
					"AND (a.name LIKE '%%%s%%' OR r.title LIKE '%%%s%%')",
					p.SearchTerm,
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
				"SELECT Count(id)",
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

type ReviewResponse struct {
	Echo string `json:"sEcho"`
	TotalRecords int `json:"iTotalRecords"`
	TotalDisplayRecords int `json:"iTotalDisplayRecords"`
	Results []map[string]string `json:"aaData"`
}

func ReviewsToResponse(ids IDSlice, reviews Reviews, echo string, total, matching int) ReviewResponse {
	rr := ReviewResponse{
		Echo: echo,
		TotalRecords: total,
		TotalDisplayRecords: matching,
		Results: make([]map[string]string, len(reviews)),
	}
	log.Printf("ReviewsToResponse: %d Reviews", len(reviews))
	for i, id := range ids {
		review, ok := reviews[id]
		if !ok {
			panic("impossible")
		}
		// These column names should match mDataProps in the .js
		result := map[string]string{
			"ID": fmt.Sprintf("%d", review.ID),
			"Title": review.Permalink,
			"Author": review.Author,
			"Pitchformulaity": fmt.Sprintf("%d", review.Scores["Pitchformulaity"]),
			"Naïve sentence length": fmt.Sprintf("%d", review.Scores["Naïve sentence length"]),
			"Word count": fmt.Sprintf("%d", review.Scores["Word count"]),
			"Words invented": fmt.Sprintf("%d", review.Scores["Words invented"]),
		}
		log.Printf("ReviewsToResponse: %d/%d: %d", i, len(reviews), review.ID)
		rr.Results[i] = result
	}
	return rr
}