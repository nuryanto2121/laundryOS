package bot

import (
	"regexp"
	"strings"
	"unicode"
)

var whitespaceRegex = regexp.MustCompile(`\s+`)

// Matcher performs rule-based keyword matching on customer messages.
type Matcher struct {
	items []FAQItem
}

// NewMatcher creates a new Matcher with the given FAQ knowledge base items.
func NewMatcher(items []FAQItem) *Matcher {
	return &Matcher{items: items}
}

// Normalize cleans and standardizes input text for reliable matching.
func Normalize(text string) string {
	lower := strings.ToLower(text)

	// Replace common punctuation with spaces
	cleaned := strings.Map(func(r rune) rune {
		if unicode.IsPunct(r) || unicode.IsSymbol(r) {
			return ' '
		}
		return r
	}, lower)

	// Collapse multiple whitespaces and trim
	trimmed := strings.TrimSpace(whitespaceRegex.ReplaceAllString(cleaned, " "))
	return trimmed
}

// Match evaluates input text against all FAQ categories and returns the best matching category and its score.
func (m *Matcher) Match(text string) (category string, score int) {
	norm := Normalize(text)
	if norm == "" {
		return CategoryUnknown, 0
	}

	tokens := strings.Fields(norm)
	tokenSet := make(map[string]struct{}, len(tokens))
	for _, t := range tokens {
		tokenSet[t] = struct{}{}
	}

	bestCategory := CategoryUnknown
	maxScore := 0

	for _, item := range m.items {
		itemScore := 0

		for _, kw := range item.Keywords {
			normKw := Normalize(kw)
			if normKw == "" {
				continue
			}

			// 1. Exact match with single keyword (e.g. user typed "harga" or "menu")
			if norm == normKw {
				itemScore += 20
				continue
			}

			// 2. Multi-word phrase contained in input (e.g. "dry clean", "antar jemput")
			if strings.Contains(normKw, " ") {
				if strings.Contains(norm, normKw) {
					itemScore += 12
				}
				continue
			}

			// 3. Single-word token exact match
			if _, exists := tokenSet[normKw]; exists {
				itemScore += 5
				continue
			}

			// 4. Substring match for longer words (min 4 chars) to catch basic typos / affixes
			if len(normKw) >= 4 && strings.Contains(norm, normKw) {
				itemScore += 2
			}
		}

		if itemScore > maxScore {
			maxScore = itemScore
			bestCategory = item.Category
		}
	}

	if maxScore == 0 {
		return CategoryUnknown, 0
	}

	return bestCategory, maxScore
}
