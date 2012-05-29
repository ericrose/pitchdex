package main

type Author struct {
	Name string
	Reviews int
	Scores map[string]int
}

type Authors map[string]Author