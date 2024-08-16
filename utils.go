package pgs

import (
	"strings"
	"unicode"
)

func toSnakeCase(str string) string {
	var result string
	var words []string
	var lastPosition int
	runes := []rune(str)

	for i, character := range runes {
		if i == 0 {
			continue
		}

		if unicode.IsUpper(character) {
			words = append(words, str[lastPosition:i])
			lastPosition = i
		}
	}

	// Добавить последнее слово в slice
	words = append(words, str[lastPosition:])

	for _, word := range words {
		if result != "" {
			result += "_"
		}
		result += strings.ToLower(word)
	}

	return result
}
