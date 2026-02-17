package search

// Package search provides full-text search utilities.
// The actual FTS5 queries are implemented in the db/repository.go
// This package provides query preprocessing and result ranking helpers.

import "strings"

// PrepareQuery cleans and prepares a search query for FTS5.
func PrepareQuery(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	// Split into tokens, add prefix matching to last token
	tokens := strings.Fields(raw)
	if len(tokens) == 0 {
		return ""
	}

	// Escape FTS5 special characters
	for i, t := range tokens {
		t = strings.ReplaceAll(t, "\"", "")
		t = strings.ReplaceAll(t, "'", "")
		tokens[i] = t
	}

	// Add prefix matching to the last token (for type-ahead search)
	tokens[len(tokens)-1] = tokens[len(tokens)-1] + "*"

	return strings.Join(tokens, " ")
}
