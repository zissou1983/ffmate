package ffmpeg

import (
	"fmt"
	"strings"
	"unicode"
)

func shellwordsUnicodeSafe(input string) ([]string, error) {
	var args []string
	var current strings.Builder
	inQuotes := false
	var quoteChar rune

	for _, r := range input {
		switch {
		case unicode.IsSpace(r) && !inQuotes:
			if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		case r == '"' || r == '\'':
			if inQuotes && r == quoteChar {
				inQuotes = false
			} else if !inQuotes {
				inQuotes = true
				quoteChar = r
			} else {
				current.WriteRune(r)
			}
		default:
			current.WriteRune(r)
		}
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	if inQuotes {
		return nil, fmt.Errorf("unclosed quote")
	}
	return args, nil
}
