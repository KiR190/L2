package parser

import (
	"os"
	"regexp"
	"strings"

	. "minishell/internal/types"
)

func ParseLine(line string, builtins map[string]BuiltinCommand) (*Pipeline, error) {
	tokens := splitByConditional(line)
	pipeline := &Pipeline{}

	var prevCondLast *Command // последняя команда предыдущей условной части
	prevTokenOp := ""         // оператор, который следовал за предыдущей частью (&&, ||, или "")

	for _, token := range tokens {
		part := strings.TrimSpace(token.Text)
		if part == "" {
			prevTokenOp = token.Op
			continue
		}

		pipeParts := strings.Split(part, "|")
		var firstPipe, prevPipe *Command

		for _, p := range pipeParts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}

			fields := strings.Fields(p)			
			if len(fields) == 0 {
				continue
			}

			cmd := &Command{
				Name: fields[0],
				Args: expandVariables(fields[1:]),
			}

			if _, ok := builtins[cmd.Name]; ok {
				cmd.Type = Builtin
			} else {
				cmd.Type = External
			}

			// Редиректы
			for i := 0; i < len(cmd.Args); i++ {
				switch cmd.Args[i] {
				case ">":
					if i+1 < len(cmd.Args) {
						cmd.Output = cmd.Args[i+1]
						cmd.Args = append(cmd.Args[:i], cmd.Args[i+2:]...)
						i--
					}
				case ">>":
					if i+1 < len(cmd.Args) {
						cmd.Output = cmd.Args[i+1]
						cmd.Append = true
						cmd.Args = append(cmd.Args[:i], cmd.Args[i+2:]...)
						i--
					}
				case "<":
					if i+1 < len(cmd.Args) {
						cmd.Input = cmd.Args[i+1]
						cmd.Args = append(cmd.Args[:i], cmd.Args[i+2:]...)
						i--
					}
				}
			}

			// Связываем пайпы
			if prevPipe != nil {
				prevPipe.NextPipe = cmd
			} else {
				firstPipe = cmd
			}
			prevPipe = cmd
		}

		// Добавляем первую команду этой цепочки в pipeline (если есть)
		if firstPipe != nil {
			pipeline.Commands = append(pipeline.Commands, firstPipe)
		}

		// Связываем условные блоки: используем prevCondLast и prevTokenOp
		if prevCondLast != nil && prevTokenOp != "" {
			prevCondLast.NextOp = prevTokenOp
			prevCondLast.NextCond = firstPipe
		}

		// Последняя команда этой части становится prevCondLast
		prevCondLast = prevPipe

		// Сохраняем Op текущего токена как "предыдущий" для следующей итерации
		prevTokenOp = token.Op
	}

	return pipeline, nil
}

func expandVariables(args []string) []string {
	re := regexp.MustCompile(`\$(\w+)|\$\{(\w+)\}`)
	res := make([]string, len(args))

	for i, arg := range args {
		res[i] = re.ReplaceAllStringFunc(arg, func(match string) string {
			varName := ""
			if match[1] == '{' {
				varName = match[2 : len(match)-1]
			} else {
				varName = match[1:]
			}
			val := os.Getenv(varName)
			return val
		})
	}

	return res
}

type ConditionalPart struct {
	Text string
	Op   string // "", "&&", "||"
}

// Разделяет строку на куски по && и ||
func splitByConditional(line string) []ConditionalPart {
	var parts []ConditionalPart
	var current strings.Builder

	for i := 0; i < len(line); i++ {
		if i+1 < len(line) {
			if line[i] == '&' && line[i+1] == '&' {
				parts = append(parts, ConditionalPart{Text: current.String(), Op: "&&"})
				current.Reset()
				i++
				continue
			}
			if line[i] == '|' && line[i+1] == '|' {
				parts = append(parts, ConditionalPart{Text: current.String(), Op: "||"})
				current.Reset()
				i++
				continue
			}
		}
		current.WriteByte(line[i])
	}

	parts = append(parts, ConditionalPart{Text: current.String(), Op: ""})
	return parts
}
