package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"

	"minishell/internal/builtins"
	"minishell/internal/executor"
	"minishell/internal/parser"
)

// Специальная ошибка для выхода из шелла
var ErrExit = errors.New("exit")

func main() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)

	go func() {
		for range sigCh {
			fmt.Println()
			fmt.Print("minishell> ")
		}
	}()

	builtinCommands := builtins.InitBuiltins()
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("minishell> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println()
				break
			}
			log.Printf("read error: %v\n", err)
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		pipeline, err := parser.ParseLine(line, builtinCommands)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
			continue
		}

		if err := executor.ExecutePipeline(pipeline, builtinCommands); err != nil {
			if errors.Is(err, ErrExit) {
				break
			}
			fmt.Fprintf(os.Stderr, "execution error: %v\n", err)
		}
	}
}
