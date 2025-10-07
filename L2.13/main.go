package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// parseFields парсит строку из -f (например, "1,3-5") в map[int]bool
func parseFields(spec string) (map[int]bool, error) {
	fields := make(map[int]bool)
	parts := strings.Split(spec, ",")
	for _, p := range parts {
		if strings.Contains(p, "-") {
			// диапазон
			r := strings.SplitN(p, "-", 2)
			if len(r) != 2 {
				return nil, fmt.Errorf("invalid range: %s", p)
			}
			start, err1 := strconv.Atoi(r[0])
			end, err2 := strconv.Atoi(r[1])
			if err1 != nil || err2 != nil || start <= 0 || end <= 0 || start > end {
				return nil, fmt.Errorf("invalid range: %s", p)
			}
			for i := start; i <= end; i++ {
				fields[i] = true
			}
		} else {
			// одно число
			n, err := strconv.Atoi(p)
			if err != nil || n <= 0 {
				return nil, fmt.Errorf("invalid field: %s", p)
			}
			fields[n] = true
		}
	}
	return fields, nil
}

func main() {
	fieldSpec := flag.String("f", "", "fields (columns) to select, e.g. 1,3-5")
	delimiter := flag.String("d", "\t", "delimiter (default tab)")
	separated := flag.Bool("s", false, "only print lines with delimiter")

	flag.Parse()

	if *fieldSpec == "" {
		fmt.Fprintln(os.Stderr, "usage: cut -f fields [-d delim] [-s]")
		os.Exit(1)
	}

	fields, err := parseFields(*fieldSpec)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid -f:", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		// если -s и нет разделителя
		if *separated && !strings.Contains(line, *delimiter) {
			continue
		}

		cols := strings.Split(line, *delimiter)
		out := []string{}

		for i := 1; i <= len(cols); i++ {
			if fields[i] {
				out = append(out, cols[i-1])
			}
		}

		if len(out) > 0 {
			fmt.Println(strings.Join(out, *delimiter))
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
		os.Exit(1)
	}
}
