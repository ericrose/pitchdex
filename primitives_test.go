package main

import (
	"testing"
)

func TestStripHTML(t *testing.T) {
	snippets := map[string]string{
		`<p>Easy <strong>first</strong> one.</p>`:    "Easy first one.",
		`A <em><a href="xyz">trickier</a> one</em>.`: "A trickier one.",
		`<b><b>Unmatched</b> <em>entries.`:           "Unmatched entries.",
	}
	for from, expected := range snippets {
		got := stripHTML(from)
		if got != expected {
			t.Errorf("'%s': got '%s', expected '%s'", from, got, expected)
		}
	}
}

func TestTokenize(t *testing.T) {
	snippets := map[string][]string{
		`<p>Easy <em>first</em> one.</p>`: []string{"easy", "first", "one"},
		`<p><p>Un-matched thing!</p>`:     []string{"un-matched", "thing"},
	}
	for from, expected := range snippets {
		got := tokenize(from)
		if !equal(got, expected) {
			t.Errorf("'%s': got '%v', expected '%v'", from, got, expected)
		}
	}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
