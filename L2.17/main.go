package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"telnet/internal/client"
)

func main() {
	var timeout = flag.Duration("timeout", 10*time.Second, "connection timeout")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(),
			"Usage: %s [--timeout duration] host port\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	host := flag.Arg(0)
	port := flag.Arg(1)
	addr := fmt.Sprintf("%s:%s", host, port)

	// Контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	// Ctrl+C и SIGTERM через signal.NotifyContext
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Создаем и запускаем клиет
	cl := client.New(addr, os.Stdin, os.Stdout)
	if err := cl.Run(ctx); err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	log.Println("Соединение закрыто. Завершение работы.")
}
