package ocr

import (
	"strings"
	"unicode"
)

func words(text string) []string {
	text = strings.ReplaceAll(strings.ToLower(text), "ё", "е")
	return strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
}

func hasBlacklistedWord(text string, blacklist []string, threshold float64) bool {
	inputWords := words(text)
	blacklistWords := words(strings.Join(blacklist, " "))
	for _, word := range inputWords {
		for _, blocked := range blacklistWords {
			if wordSimilarity(word, blocked) >= threshold {
				return true
			}
		}
	}
	return false
}

// Levenshtein distance
func wordSimilarity(a, b string) float64 {
	x, y := []rune(a), []rune(b)
	if len(x) < len(y) {
		x, y = y, x
	}
	if len(x) == 0 {
		return 1
	}
	row := make([]int, len(y)+1)
	for j := range row {
		row[j] = j
	}
	for i, left := range x {
		diagonal := row[0]
		row[0] = i + 1
		for j, right := range y {
			previous := row[j+1]
			cost := 0
			if left != right {
				cost = 1
			}
			row[j+1] = min(row[j]+1, previous+1, diagonal+cost)
			diagonal = previous
		}
	}
	return 1 - float64(row[len(y)])/float64(len(x))
}
