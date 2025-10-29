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

	// Котекст для завершения
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		select {
		case <-sigCh:
			log.Println("Получен сигнал завершения, закрываю соединение...")
			cancel()
		case <-ctx.Done():
		}
	}()

	// Создаем и запускаем клиет
	cl := client.New(addr, os.Stdin, os.Stdout)
	if err := cl.Run(ctx); err != nil {
		log.Fatalf("Ошибка: %v", err)
	}

	log.Println("Соединение закрыто. Завершение работы.")
}
