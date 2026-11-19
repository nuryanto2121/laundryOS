package bot

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/nuryanto2121/laundryos/apps/api/internal/model"
)

var whitespaceRegex = regexp.MustCompile(`\s+`)
var introNameRegex = regexp.MustCompile(`(?i)(?:nama saya|saya|aku)\s+([a-zA-Z]{3,15})`)

var commonNonNames = map[string]struct{}{
	"mau": {}, "ingin": {}, "minta": {}, "sudah": {}, "udh": {}, "bisa": {},
	"lagi": {}, "dari": {}, "ada": {}, "baru": {}, "pelanggan": {}, "customer": {},
	"di": {}, "ke": {}, "dan": {}, "yg": {}, "yang": {},
}

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

// ExtractCustomerName inspects text for self-introduction patterns (e.g. "saya nana teman'a nini mba").
func ExtractCustomerName(text string) string {
	matches := introNameRegex.FindStringSubmatch(text)
	if len(matches) > 1 {
		candidate := strings.ToLower(strings.TrimSpace(matches[1]))
		if _, isStopword := commonNonNames[candidate]; !isStopword && len(candidate) >= 3 {
			// Title case the candidate
			return strings.ToUpper(candidate[:1]) + candidate[1:]
		}
	}
	return ""
}

// ServiceKeywordMap maps common customer colloquial terms to canonical service names.
var ServiceKeywordMap = []struct {
	CanonicalName string
	Keywords      []string
}{
	// Self Service (Mesin Tabung)
	{
		CanonicalName: "Cuci Basah (Self Service)",
		Keywords:      []string{"cuci basah", "self service basah", "basah", "tabung basah"},
	},
	{
		CanonicalName: "Cuci Kering (Self Service)",
		Keywords:      []string{"cuci kering", "self service kering", "self service", "cuci mesin", "tabung", "cuci tabung"},
	},

	// Cuci Setrika
	{
		CanonicalName: "Cuci Setrika Turbo",
		Keywords:      []string{"cuci turbo", "setrika turbo", "turbo", "cuci 6 jam"},
	},
	{
		CanonicalName: "Cuci Setrika Super Kilat",
		Keywords:      []string{"super kilat", "cuci super kilat", "12 jam"},
	},
	{
		CanonicalName: "Cuci Setrika Kilat",
		Keywords:      []string{"cuci kilat", "cuci setrika kilat", "kilat", "24 jam", "1 hari"},
	},
	{
		CanonicalName: "Cuci Setrika Express",
		Keywords:      []string{"cuci express", "cuci ekspres", "express", "ekspres", "2 hari"},
	},
	{
		CanonicalName: "Cuci Komplit Reguler",
		Keywords:      []string{"cuci regular", "cuci reguler", "cuci setrika", "komplit reguler", "kiloan reguler", "reguler", "regular"},
	},

	// Setrika Saja
	{
		CanonicalName: "Setrika Saja",
		Keywords:      []string{"setrika", "setrika saja", "gosok", "gosok saja", "hanya setrika"},
	},

	// Bad Cover
	{
		CanonicalName: "Cuci Bed Cover",
		Keywords:      []string{"bad cover", "badcover", "bed cover", "bedcover", "selimut", "sprei", "bedding"},
	},

	// Karpet Permadani
	{
		CanonicalName: "Karpet Permadani",
		Keywords:      []string{"karpet", "permadani", "cuci karpet", "kasur busa", "kasur lantai", "dacron"},
	},

	// Boneka
	{
		CanonicalName: "Cuci Boneka",
		Keywords:      []string{"boneka", "cuci boneka", "boneka besar", "boneka kecil", "boneka sedang", "stroller"},
	},
}

// MatchServiceIntent checks if the customer's message expresses an intention to select a specific service.
// Returns the matching ServiceItem if found, or nil.
func MatchServiceIntent(text string, services []model.ServiceItem) *model.ServiceItem {
	norm := Normalize(text)
	if norm == "" {
		return nil
	}

	for _, entry := range ServiceKeywordMap {
		matched := false
		for _, kw := range entry.Keywords {
			if strings.Contains(norm, kw) {
				matched = true
				break
			}
		}

		if matched {
			// Look for matching service in dynamic list
			for _, s := range services {
				if strings.Contains(strings.ToLower(s.ServiceName), strings.ToLower(entry.CanonicalName)) ||
					strings.Contains(strings.ToLower(entry.CanonicalName), strings.ToLower(s.ServiceName)) {
					item := s
					return &item
				}
			}

			// Fallback placeholder item if DB services not populated
			return &model.ServiceItem{
				ServiceName: entry.CanonicalName,
				PricePerKg:  8000,
				DurationDay: 2,
			}
		}
	}

	return nil
}

var (
	hashOrderRegex = regexp.MustCompile(`(?i)#([0-9]{10,14})`)
	ordOrderRegex  = regexp.MustCompile(`(?i)\b(ORD[0-9A-Z]{6,16})\b`)
	pureNumRegex   = regexp.MustCompile(`\b(2[0-9]{11})\b`)
)

// ExtractOrderNumber extracts and normalizes an order number from text,
// supporting "#261119000012", "ORD261119000012", or "261119000012".
// Always returns the canonical database format e.g. "ORD261119000012".
func ExtractOrderNumber(text string) string {
	if m := hashOrderRegex.FindStringSubmatch(text); len(m) > 1 {
		return "ORD" + m[1]
	}
	if m := ordOrderRegex.FindStringSubmatch(text); len(m) > 1 {
		return strings.ToUpper(m[1])
	}
	if m := pureNumRegex.FindStringSubmatch(text); len(m) > 1 {
		return "ORD" + m[1]
	}
	return ""
}

// ParseWeightCommand parses weight values like "TIMBANG ORD... 9.5", "timbang 5,5", "8 kg", or "9.5".
func ParseWeightCommand(text string) (float64, bool) {
	norm := strings.TrimSpace(text)
	norm = strings.ReplaceAll(norm, ",", ".")

	numRegex := regexp.MustCompile(`\b([0-9]+(?:\.[0-9]+)?)\s*(?:kg)?\b`)
	matches := numRegex.FindAllStringSubmatch(norm, -1)
	for _, m := range matches {
		if len(m) > 1 {
			val, err := strconv.ParseFloat(m[1], 64)
			if err == nil && val > 0 && val < 500 {
				return val, true
			}
		}
	}
	return 0, false
}

// IsWeighCommand checks if the text indicates a weighing operation.
func IsWeighCommand(text string) bool {
	lower := strings.ToLower(text)
	return strings.Contains(lower, "timbang") || strings.Contains(lower, "berat")
}

// IsCompleteCommand checks if the text indicates order completion.
func IsCompleteCommand(text string) bool {
	lower := strings.ToLower(text)
	return strings.Contains(lower, "selesai") || strings.Contains(lower, "done") || strings.Contains(lower, "ready")
}
