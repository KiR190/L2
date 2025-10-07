package unpack

import (
	"errors"
	"strings"
	"unicode"
)

func Unpack(input string) (string, error) {
	if input == "" {
		return "", nil
	}

	var b strings.Builder
	var prevRune rune
	var hasPrev bool
	escaped := false

	for _, r := range input {
		switch {
		case escaped:
			b.WriteRune(r)
			prevRune = r
			hasPrev = true
			escaped = false
		case r == '\\':
			escaped = true
		case unicode.IsDigit(r):
			if !hasPrev {
				return "", errors.New("invalid string: starts with digit")
			}
			count := int(r - '0')
			if count == 0 {
				result := []rune(b.String())
				b.Reset()
				b.WriteString(string(result[:len(result)-1]))
				continue
			}

			for i := 1; i < count; i++ {
				b.WriteRune(prevRune)
			}

		default:
			b.WriteRune(r)
			prevRune = r
			hasPrev = true
		}
	}

	if escaped {
		return "", errors.New("invalid string: dangling escape")
	}

	return b.String(), nil
}
