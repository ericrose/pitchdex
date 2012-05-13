package main

import (
	"bytes"
	"github.com/peterbourgon/exp-html"
	"strings"
)

var IndexDefinitions = IndexMap{
	"Reviews":               SimpleCount,
	"Pitchformulaity":       Pitchformulaity,
	"Naïve sentence length": NaïveSentenceLength,
	"Words invented":        InventedWordsFunc(*dictFile),
	"Character count":       CharacterCount,
	"Word count":            WordCount,
	"Word length":           AverageWordLength,
}

const BullshitScore = "Overall Bullshit Score"

//
//
//

func SimpleCount(r Review) int { return 1 }

func WordCount(r Review) int {
	return len(tokenize(r.Body))
}

func CharacterCount(r Review) int {
	return len(stripHTML(r.Body))
}

func AverageWordLength(r Review) int {
	return int(float64(CharacterCount(r)) / float64(WordCount(r)))
}

func InventedWordsFunc(dictfile string) func(r Review) int {
	dict := NewDict(dictfile)
	return func(r Review) int {
		count := 0
		for _, word := range tokenize(r.Body) {
			if !dict.Has(word) {
				// fmt.Printf("invented '%s'\n", baseWord(word))
				count++
			}
		}
		return count
	}
}

func NaïveSentenceLength(r Review) int {
	i, sentences, b := 0, 0, stripHTML(r.Body)
	for {
		j := strings.Index(b[i:], ".")
		if j < 0 {
			break
		}
		sentences++
		i = i + j + 1
	}
	return int(float64(WordCount(r)) / float64(sentences))
}

func Pitchformulaity(r Review) int {
	score := 0
	for _, word := range tokenize(r.Body) {
		if n, ok := PitchformulaWords[word]; ok {
			score += n
		}
	}
	return score
}

//
//
//

func tokenize(body string) []string {
	toks := strings.Split(stripHTML(body), " ")
	for i, tok := range toks {
		toks[i] = baseWord(tok)
	}
	return toks
}

func baseWord(word string) string {
	return strings.ToLower(strings.Trim(word, ` ,.;!?"-`))
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
