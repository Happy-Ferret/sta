package parser

import (
	"github.com/ribacq/sta/commands"
	"github.com/ribacq/sta/context"
	"strings"
)

// builtin commands
const (
	Builtins = "quit help"
)

const (
	NilFlag = iota
	QuitFlag
	ErrFlag
)

type ParserOutput struct {
	Flag int
	Message string
}

// get all possible commands for current context
func CommandList() []string {
	return strings.Split(Builtins, " ")
}

// get whether a command exists
func Exists(cmd *commands.Command) bool {
	for _, c := range CommandList() {
		if cmd.Cmd == c {
			return true
		}
	}
	return false
}

// get whether a command is a builtin
func IsBuiltin(cmd *commands.Command) bool {
	for _, c := range strings.Split(Builtins, " ") {
		if cmd.Cmd == c {
			return true
		}
	}
	return false
}

// execute builtin command
func ExecBuiltin(cmd *commands.Command, c *context.Context) (*ParserOutput, error) {
	if !IsBuiltin(cmd) {
		return nil, commands.UnknownCommandError{*cmd}
	}
	switch cmd.Cmd {
	case "help":
		if len(cmd.Args) > 0 {
			arg := cmd.Args[0]
			if IsBuiltin(&commands.Command{arg, nil}) {
				return &ParserOutput{NilFlag, "Built-in: " + arg}, nil
			} else if c.HasCommand(&commands.Command{arg, nil}) {
				return &ParserOutput{NilFlag, c.CommandActions[arg].Help}, nil
			} else {
				return &ParserOutput{ErrFlag, "Unknown command: " + arg}, nil
			}
		}
		helpStr := "Available commands: " + strings.Join(CommandList(), ", ") + strings.Join(c.CommandList(), ", ")
		return &ParserOutput{NilFlag, helpStr}, nil
	case "quit":
		return &ParserOutput{QuitFlag, "Fare well."}, nil
	}
	return nil, nil
}

// returns a command from a line
func Parse(line string) *commands.Command {
	sections := strings.Split(line, " ")
	var args []string
	for _, arg := range sections[1:] {
		arg = strings.Trim(arg, " ")
		if arg != "" {
			args = append(args, arg)
		}
	}
	cmd := &commands.Command{sections[0], args}
	return cmd
}
