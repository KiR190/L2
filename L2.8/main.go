package main

import (
	"flag"
	"time"

	"ntp-time/internal/timeSync"
)

func main() {
	server := flag.String("server", "ru.pool.ntp.org", "NTP server to query (host or host:port)")
	timeout := flag.Duration("timeout", 5*time.Second, "Timeout for the NTP query")
	format := flag.String("format", time.RFC3339Nano, "Output time format")
	flag.Parse()

	timeSync.PrintCurrentTime(*server, *timeout, *format)
}
