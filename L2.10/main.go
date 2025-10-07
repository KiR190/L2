package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Месяцы для -M
var months = map[string]int{
	"Jan": 1, "Feb": 2, "Mar": 3, "Apr": 4,
	"May": 5, "Jun": 6, "Jul": 7, "Aug": 8,
	"Sep": 9, "Oct": 10, "Nov": 11, "Dec": 12,
}

func expandShortFlags(args []string) []string {
	var result []string
	for _, arg := range args {
		if len(arg) > 2 && arg[0] == '-' && arg[1] != '-' {
			for _, ch := range arg[1:] {
				result = append(result, "-"+string(ch))
			}
		} else {
			result = append(result, arg)
		}
	}
	return result
}

func main() {
	col := flag.Int("k", 0, "column to sort by (1-based)")
	numeric := flag.Bool("n", false, "sort by numeric value")
	reverse := flag.Bool("r", false, "reverse order")
	unique := flag.Bool("u", false, "unique lines only")
	monthSort := flag.Bool("M", false, "sort by month names")
	ignoreTrailing := flag.Bool("b", false, "ignore trailing blanks")
	checkSorted := flag.Bool("c", false, "check if sorted")
	human := flag.Bool("h", false, "sort by human-readable sizes (e.g., 10K, 2M)")

	expandedArgs := expandShortFlags(os.Args[1:])
	
	if err := flag.CommandLine.Parse(expandedArgs); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		flag.Usage()
		os.Exit(2)
	}

	var scanner *bufio.Scanner
	if flag.NArg() > 0 {
		file, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, "error opening file:", err)
			os.Exit(1)
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	} else {
		scanner = bufio.NewScanner(os.Stdin)
	}

	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if *ignoreTrailing {
			line = strings.TrimRight(line, " \t")
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
		os.Exit(1)
	}

	if *checkSorted {
		if isSorted(lines, *col, *numeric, *monthSort, *human, *reverse) {
			os.Exit(0)
		}
		fmt.Fprintln(os.Stderr, "input is not sorted")
		os.Exit(1)
	}

	sort.Slice(lines, func(i, j int) bool {
		return compare(lines[i], lines[j], *col, *numeric, *monthSort, *human, *reverse)
	})

	if *unique {
		lines = uniq(lines)
	}

	for _, l := range lines {
		fmt.Println(l)
	}
}

func getKey(line string, col int) string {
	if col <= 0 {
		return line
	}
	parts := strings.Split(line, "\t")
	if col-1 < len(parts) {
		return parts[col-1]
	}
	return ""
}

func compare(a, b string, col int, numeric, monthSort, human, reverse bool) bool {
	ka := getKey(a, col)
	kb := getKey(b, col)

	var less bool
	switch {
	case monthSort:
		less = months[ka] < months[kb]
	case human:
		va := parseHuman(ka)
		vb := parseHuman(kb)
		less = va < vb
	case numeric:
		na, _ := strconv.ParseFloat(ka, 64)
		nb, _ := strconv.ParseFloat(kb, 64)
		less = na < nb
	default:
		less = ka < kb
	}

	if reverse {
		return !less
	}
	return less
}

func isSorted(lines []string, col int, numeric, monthSort, human, reverse bool) bool {
	for i := 1; i < len(lines); i++ {
		if compare(lines[i], lines[i-1], col, numeric, monthSort, human, reverse) {
			return false
		}
	}
	return true
}

func uniq(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}
	res := []string{lines[0]}
	for i := 1; i < len(lines); i++ {
		if lines[i] != lines[i-1] {
			res = append(res, lines[i])
		}
	}
	return res
}

func parseHuman(s string) float64 {
	if s == "" {
		return 0
	}
	last := s[len(s)-1]
	mult := 1.0
	switch last {
	case 'K', 'k':
		mult = 1 << 10
		s = s[:len(s)-1]
	case 'M':
		mult = 1 << 20
		s = s[:len(s)-1]
	case 'G':
		mult = 1 << 30
		s = s[:len(s)-1]
	}
	v, _ := strconv.ParseFloat(s, 64)
	return v * mult
}
