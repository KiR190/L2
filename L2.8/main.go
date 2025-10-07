package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	server := flag.String("server", "ru.pool.ntp.org", "NTP server to query (host or host:port)")
	timeout := flag.Duration("timeout", 5*time.Second, "Timeout for the NTP query")
	format := flag.String("format", time.RFC3339Nano, "Output time format")
	flag.Parse()

	opts := ntp.QueryOptions{
		Timeout: *timeout,
	}

	// Выполняем запрос с опциями
	resp, err := ntp.QueryWithOptions(*server, opts)
	if err != nil {
		// Печатаем ошибку в STDERR и выходим с кодом 1
		fmt.Fprintln(os.Stderr, "ntp query failed:", err)
		os.Exit(1)
	}

	if err := resp.Validate(); err != nil {
		fmt.Fprintln(os.Stderr, "ntp response validation failed:", err)
		os.Exit(1)
	}

	precise := time.Now().Add(resp.ClockOffset).UTC()

	// Выводим в заданном формате на STDOUT
	fmt.Println(precise.Format(*format))
}
