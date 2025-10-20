package types

import "io"

type CommandType int
type ConditionalType int

const (
	Builtin CommandType = iota
	External
)

type Command struct {
	Name     string
	Args     []string
	Type     CommandType
	CondType ConditionalType

	Input    string
	Output   string
	Append   bool
	NextPipe *Command
	NextCond *Command
	NextOp   string
}

type Pipeline struct {
	Commands []*Command
}

// Тип встроенной команды
type BuiltinCommand func(args []string, input io.Reader, output io.Writer) error
