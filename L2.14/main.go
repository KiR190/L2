package main

import (
	"fmt"
	"time"
)

var or func(channels ...<-chan interface{}) <-chan interface{}

func init() {
	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {
		case 0:
			c := make(chan interface{})
			close(c)
			return c
		case 1:
			return channels[0]
		}

		orDone := make(chan interface{})
		go func() {
			defer close(orDone)

			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				m := len(channels) / 2
				select {
				case <-or(channels[:m]...):
				case <-or(channels[m:]...):
				}
			}
		}()
		return orDone
	}
}

func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func main() {
	start := time.Now()

	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("done after %v", time.Since(start))
}
