package timeSync

import (
	"fmt"
	"os"
	"time"

	"github.com/beevik/ntp"
)

// GetCurrentTimeNTP получает точное текущее время с указанного NTP-сервера.
func GetCurrentTimeNTP(server string, timeout time.Duration) (time.Time, error) {
	opts := ntp.QueryOptions{
		Timeout: timeout,
	}

	resp, err := ntp.QueryWithOptions(server, opts)
	if err != nil {
		return time.Time{}, fmt.Errorf("ntp query failed: %w", err)
	}

	if err := resp.Validate(); err != nil {
		return time.Time{}, fmt.Errorf("ntp response validation failed: %w", err)
	}

	precise := time.Now().Add(resp.ClockOffset)
	return precise, nil
}

// PrintCurrentTime получает время с NTP-сервера и печатает в консоль.
func PrintCurrentTime(server string, timeout time.Duration, format string) {
	currentTime, err := GetCurrentTimeNTP(server, timeout)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(currentTime.Local().Format(format))
}
