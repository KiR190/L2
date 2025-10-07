package main

import (
	"fmt"
	"sort"
	"strings"
)

func findAnagrams(words []string) map[string][]string {
	groups := make(map[string][]string)

	for _, w := range words {
		w = strings.ToLower(w)
		key := sortRunes(w)
		groups[key] = append(groups[key], w)
	}

	result := make(map[string][]string)

	for _, group := range groups {
		if len(group) < 2 {
			continue
		}
		sort.Strings(group)
		result[group[0]] = group
	}

	return result
}

func sortRunes(s string) string {
	runes := []rune(s)
	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})
	return string(runes)
}

func main() {
	words := []string{"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}

	anagrams := findAnagrams(words)

	for k, v := range anagrams {
		fmt.Printf("%q: %q\n", k, v)
	}
}
