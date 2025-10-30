package client

import (
	"context"
	"errors"
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

	log.Printf("Подключено к %s\n", c.addr)

	// Канал для сигнала о завершении
	go func() {
		<-ctx.Done()
		conn.Close()
	}()

	done := make(chan struct{})
	var once sync.Once
	closeOnce := func() {
		once.Do(func() {
			conn.Close()
			close(done)
		})
	}

	go func() {
		if _, err := io.Copy(conn, c.in); err != nil && !errors.Is(err, net.ErrClosed) {
			if err != io.EOF {
				log.Printf("Ошибка записи в сокет: %v", err)
			} else {
				log.Println("EOF от stdin — закрываю соединение.")
			}
		}
		closeOnce()
	}()

	// Горутина stdout
	go func() {
		if _, err := io.Copy(c.out, conn); err != nil && !errors.Is(err, net.ErrClosed) {
			if err == io.EOF {
				log.Println("Сервер закрыл соединение.")
			} else {
				log.Printf("Ошибка чтения из сокета: %v", err)
			}
		}
		closeOnce()
	}()

	<-done
	return nil
}
