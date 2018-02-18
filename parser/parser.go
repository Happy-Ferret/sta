package parser

import (
	"strings"
)

// builtin commands
const (
	Builtins = "quit help"
)

// types
type Command struct {
	Cmd  string
	Args []string
}

type UnknownCommandError struct {
	Cmd Command
}

func (err UnknownCommandError) Error() string {
	return "Unknown command: " + err.Cmd.Cmd
}

const (
	NilFlag = iota
	QuitFlag
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
func (cmd Command) Exists() bool {
	for _, c := range CommandList() {
		if cmd.Cmd == c {
			return true
		}
	}
	return false
}

// get whether a command is a builtin
func (cmd Command) IsBuiltin() bool {
	for _, c := range strings.Split(Builtins, " ") {
		if cmd.Cmd == c {
			return true
		}
	}
	return false
}

// execute builtin command
func ExecBuiltin(cmd *Command) (*ParserOutput, error) {
	if !cmd.IsBuiltin() {
		return nil, UnknownCommandError{*cmd}
	}
	switch cmd.Cmd {
	case "help":
		return &ParserOutput{NilFlag, strings.Join(CommandList(), ", ")}, nil
	case "quit":
		return &ParserOutput{QuitFlag, "Fare well."}, nil
	}
	return nil, nil
}

// prompts the user for a command
func Parse(line string) *Command {
	sections := strings.Split(line, " ")
	for i := range sections {
		sections[i] = strings.Trim(sections[i], " ")
	}
	cmd := &Command{sections[0], sections[1:]}
	return cmd
}
