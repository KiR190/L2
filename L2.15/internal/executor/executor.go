package executor

import (
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	. "minishell/internal/types"
)

func ExecutePipeline(p *Pipeline, builtinCommands map[string]BuiltinCommand) error {
	var lastStatus error
	var lastOp string

	for i := 0; i < len(p.Commands); i++ {
		cmd := p.Commands[i]

		// Условное выполнение (&& / ||)
		if lastOp == "&&" && lastStatus != nil {
			continue 
		}
		if lastOp == "||" && lastStatus == nil {
			continue
		}

		lastStatus = executeChain(cmd, builtinCommands)
		lastOp = cmd.NextOp
	}

	return lastStatus
}

func prepareStdin(cmd *Command, prevReader io.Reader) (io.Reader, func() error, error) {
	if cmd.Input != "" {
		f, err := os.Open(cmd.Input)
		if err != nil {
			return nil, nil, err
		}
		return f, f.Close, nil
	} else if prevReader != nil {
		return prevReader, nil, nil
	} else {
		return os.Stdin, nil, nil
	}
}

func prepareStdout(cmd *Command) (io.Writer, io.ReadCloser, func() error, error) {
	var pipeWriter *io.PipeWriter
	var nextReader io.ReadCloser
	var closer func() error

	if cmd.NextPipe != nil {
		r, w := io.Pipe()
		pipeWriter = w
		nextReader = r
		closer = func() error { return w.Close() }
		return pipeWriter, nextReader, closer, nil
	} else if cmd.Output != "" {
		var f *os.File
		var err error
		if cmd.Append {
			f, err = os.OpenFile(cmd.Output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		} else {
			f, err = os.Create(cmd.Output)
		}
		if err != nil {
			return nil, nil, nil, err
		}
		closer = f.Close
		return f, nil, closer, nil
	} else {
		return os.Stdout, nil, nil, nil
	}
}

func runCommand(cmd *Command, builtins map[string]BuiltinCommand, stdin io.Reader, stdout io.Writer) error {
	if cmd.Type == Builtin {
		return builtins[cmd.Name](cmd.Args, stdin, stdout)
	} else {

		var execCmd *exec.Cmd

		if runtime.GOOS == "windows" {
			fullCmd := cmd.Name + " " + strings.Join(cmd.Args, " ")
			execCmd = exec.Command("cmd", "/C", fullCmd)
			/*fullCmd := cmd.Name
			if len(cmd.Args) > 0 {
				fullCmd += " " + joinArgsForCmd(cmd.Args)
			}
			execCmd = exec.Command("cmd", "/C", fullCmd)*/
		} else {
			execCmd = exec.Command(cmd.Name, cmd.Args...)
		}

		execCmd.Stdin = stdin
		execCmd.Stdout = stdout
		execCmd.Stderr = os.Stderr
		return execCmd.Run()
	}
}

/*func joinArgsForCmd(args []string) string {
	var out []string
	for _, a := range args {
		// если есть пробелы или кавычки — оборачиваем в кавычки и экранируем внутренние кавычки
		if strings.ContainsAny(a, " \"") {
			a = strings.ReplaceAll(a, `"`, `\"`)
			a = `"` + a + `"`
		}
		out = append(out, a)
	}
	return strings.Join(out, " ")
}*/

func executeChain(start *Command, builtins map[string]BuiltinCommand) error {
	current := start
	var prevReader io.Reader
	var wg sync.WaitGroup
	var execErr error
	var mu sync.Mutex

	for current != nil {
		stdin, closeIn, err := prepareStdin(current, prevReader)
		if err != nil {
			return err
		}

		stdout, nextReader, closeOut, err := prepareStdout(current)
		if err != nil {
			if closeIn != nil {
				closeIn()
			}
			return err
		}

		cmd := current
		in := stdin
		out := stdout
		closeReader := closeIn
		closeWriter := closeOut

		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				if closeWriter != nil {
					closeWriter()
				}
				if closeReader != nil {
					closeReader()
				}
			}()
			if err := runCommand(cmd, builtins, in, out); err != nil {
				mu.Lock()
				if execErr == nil {
					execErr = err
				}
				mu.Unlock()
			}
		}()

		prevReader = nextReader
		current = current.NextPipe
	}

	wg.Wait()
	return execErr
}