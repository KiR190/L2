package client

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type Client struct {
	addr string
	in   io.Reader
	out  io.Writer
}

// New создаёт новый telnet-клиент
func New(addr string, in io.Reader, out io.Writer) *Client {
	return &Client{
		addr: addr,
		in:   in,
		out:  out,
	}
}

// Run запуск клиента
func (c *Client) Run(ctx context.Context) error {
	dialer := net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", c.addr)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к %s: %w", c.addr, err)
	}
	defer conn.Close()

	log.Printf("Подключено к %s\n", c.addr)

	// Канал для сигнала о завершении
	done := make(chan struct{})
	var wg sync.WaitGroup

	// Горутина stdin
	wg.Go(func() {
		reader := bufio.NewReader(c.in)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				line, err := reader.ReadBytes('\n')
				if err != nil {
					// Ctrl+D или EOF
					if err == io.EOF {
						log.Println("EOF от stdin — закрываю соединение.")
					} else {
						log.Printf("Ошибка чтения stdin: %v\n", err)
					}
					return
				}
				if _, err := conn.Write(line); err != nil {
					log.Printf("Ошибка записи в сокет: %v\n", err)
					return
				}
			}
		}
	})

	// Горутина stdout
	wg.Go(func() {
		reader := bufio.NewReader(conn)
		writer := bufio.NewWriter(c.out)
		defer writer.Flush()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				data, err := reader.ReadBytes('\n')
				if err != nil {
					if err == io.EOF {
						log.Println("Сервер закрыл соединение.")
					} else {
						log.Printf("Ошибка чтения из сокета: %v\n", err)
					}
					close(done)
					return
				}
				if _, err := writer.Write(data); err != nil {
					log.Printf("Ошибка записи в stdout: %v\n", err)
					close(done)
					return
				}
				writer.Flush()
			}
		}
	})

	// Ожидание завершения
	select {
	case <-ctx.Done():
		log.Println("Таймаут или отмена контекста, закрываю соединение.")
	case <-done:
		// Сервер закрыл соединение
	}

	conn.Close()
	wg.Wait()
	return nil
}
