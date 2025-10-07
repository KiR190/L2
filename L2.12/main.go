package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Options cтруктура флагов
type Options struct {
	after     int
	before    int
	countOnly bool
	ignore    bool
	invert    bool
	fixed     bool
	lineNum   bool
}

// compilePattern подготавливает функцию проверки строки
func compilePattern(pattern string, opts Options) (func(string) bool, error) {
	if opts.ignore {
		pattern = "(?i)" + pattern
	}
	if opts.fixed {
		if opts.ignore {
			pattern = strings.ToLower(pattern)
			return func(s string) bool {
				return strings.Contains(strings.ToLower(s), pattern)
			}, nil
		}
		return func(s string) bool {
			return strings.Contains(s, pattern)
		}, nil
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}
	return func(s string) bool {
		return re.MatchString(s)
	}, nil
}

func main() {
	after := flag.Int("A", 0, "print N lines After match")
	before := flag.Int("B", 0, "print N lines Before match")
	context := flag.Int("C", 0, "print N lines of Context")
	countOnly := flag.Bool("c", false, "print only count of matching lines")
	ignore := flag.Bool("i", false, "ignore case")
	invert := flag.Bool("v", false, "invert match")
	fixed := flag.Bool("F", false, "fixed string (literal match)")
	lineNum := flag.Bool("n", false, "print line number")

	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: grep [options] PATTERN [FILE]")
		os.Exit(1)
	}

	pattern := flag.Arg(0)
	var filename string
	if flag.NArg() > 1 {
		filename = flag.Arg(1)
	}

	// -C эквивалентен -A и -B
	if *context > 0 {
		*after = *context
		*before = *context
	}

	opts := Options{
		after:     *after,
		before:    *before,
		countOnly: *countOnly,
		ignore:    *ignore,
		invert:    *invert,
		fixed:     *fixed,
		lineNum:   *lineNum,
	}

	matchFunc, err := compilePattern(pattern, opts)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid pattern:", err)
		os.Exit(1)
	}

	var scanner *bufio.Scanner
	if filename != "" {
		file, err := os.Open(filename)
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
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
		os.Exit(1)
	}

	printGrep(lines, matchFunc, opts)
}

func printGrep(lines []string, matchFunc func(string) bool, opts Options) {
	matched := make([]bool, len(lines))

	for i, line := range lines {
		ok := matchFunc(line)
		if opts.invert {
			ok = !ok
		}
		if ok {
			matched[i] = true
		}
	}

	if opts.countOnly {
		count := 0
		for _, m := range matched {
			if m {
				count++
			}
		}
		fmt.Println(count)
		return
	}

	printed := make(map[int]bool)
	for i := range lines {
		if !matched[i] {
			continue
		}

		start := i - opts.before
		if start < 0 {
			start = 0
		}
		end := i + opts.after
		if end >= len(lines) {
			end = len(lines) - 1
		}

		for j := start; j <= end; j++ {
			if printed[j] {
				continue
			}
			printed[j] = true
			if opts.lineNum {
				fmt.Printf("%d:%s\n", j+1, lines[j])
			} else {
				fmt.Println(lines[j])
			}
		}
	}
}
