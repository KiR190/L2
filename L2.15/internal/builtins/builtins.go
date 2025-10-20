package builtins

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	. "minishell/internal/types"
)

// Специальная ошибка для выхода из шелла
var ErrExit = errors.New("exit")

// Map встроенных команд
func InitBuiltins() map[string]BuiltinCommand {
	return map[string]BuiltinCommand{
		"cd":   makeCd(),
		"pwd":  makePwd(),
		"echo": makeEcho(),
		"kill": makeKill(),
		"ps":   makePs(),
		"exit": makeExit(),
	}
}

func ExecuteBuiltinWithIO(builtins map[string]BuiltinCommand, name string, args []string, input io.Reader, output io.Writer) error {
	fn, ok := builtins[name]
	if !ok {
		return fmt.Errorf("builtin command not found: %s", name)
	}
	return fn(args, input, output)
}

func makeCd() BuiltinCommand {
	return func(args []string, input io.Reader, output io.Writer) error {
		var dir string
		if len(args) == 0 {
			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("cd: cannot determine home directory: %w", err)
			}
			dir = home
		} else {
			dir = args[0]
		}
		if err := os.Chdir(dir); err != nil {
			return fmt.Errorf("cd: %w", err)
		}
		return nil
	}
}

func makePwd() BuiltinCommand {
	return func(args []string, input io.Reader, output io.Writer) error {
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("pwd: %w", err)
		}
		fmt.Fprintln(output, dir)
		return nil
	}
}

func makeEcho() BuiltinCommand {
	return func(args []string, input io.Reader, output io.Writer) error {
		fmt.Fprintln(output, strings.Join(args, " "))
		return nil
	}
}

func makeKill() BuiltinCommand {
	return func(args []string, input io.Reader, output io.Writer) error {
		if len(args) == 0 {
			return fmt.Errorf("kill: missing pid")
		}
		pidStr := args[0]
		pid, err := strconv.Atoi(pidStr)
		if err != nil {
			return fmt.Errorf("kill: invalid pid %q", pidStr)
		}

		proc, err := os.FindProcess(pid)
		if err != nil {
			return fmt.Errorf("kill: %w", err)
		}

		if err := proc.Kill(); err != nil {
			return fmt.Errorf("kill: %w", err)
		}
		return nil
	}
}

func makePs() BuiltinCommand {
	return func(args []string, input io.Reader, output io.Writer) error {
		cmd := exec.Command("ps")
		cmd.Stdout = output
		cmd.Stderr = os.Stderr
		cmd.Stdin = input
		return cmd.Run()
	}
}

func makeExit() BuiltinCommand {
	return func(args []string, input io.Reader, output io.Writer) error {
		return ErrExit
	}
}