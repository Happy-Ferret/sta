package parser

import (
	"github.com/ribacq/sta/commands"
	"github.com/ribacq/sta/context"
	"strings"
)

// builtin commands
var (
	builtins = map[string]context.CommandAction{
		"quit": context.CommandAction{
			func(c *context.Context, args []string) string {
				return "Fare well…"
			},
			"This command allows you to exit the game."},
		"help": context.CommandAction{
			func(c *context.Context, args []string) string {
				return "help…"
			},
			"This command gives you the list of available command or help for a specific command.",
		},
	}
)

// enumeratinon of flags used in ParserOutput
const (
	NilFlag = iota
	QuitFlag
	ErrFlag
)

// Output is an output type for the parser commands
type Output struct {
	Flag    int
	Message string
}

// CommandList gets all builtin commands names into an array
func CommandList() []string {
	cmds := make([]string, 0, len(builtins))
	for cmd := range builtins {
		cmds = append(cmds, cmd)
	}
	return cmds
}

// IsBuiltin gets whether a command is a builtin
func IsBuiltin(cmd *commands.Command) bool {
	for c := range builtins {
		if cmd.Cmd == c {
			return true
		}
	}
	return false
}

// ExecBuiltin executes builtin command
func ExecBuiltin(cmd *commands.Command, c *context.Context) (*Output, error) {
	if !IsBuiltin(cmd) {
		return nil, commands.UnknownCommandError{*cmd}
	}
	if cmd.Cmd == "help" {
		if len(cmd.Args) > 0 {
			// we ask for help on a specific command
			arg := cmd.Args[0]
			if IsBuiltin(&commands.Command{arg, nil}) {
				// help on a built-in
				helpStr := builtins[arg].Help + " (built-in command)"
				return &Output{NilFlag, helpStr}, nil
			} else if c.HasCommand(&commands.Command{arg, nil}) {
				// help on a context command
				helpStr := c.CommandActions[arg].Help
				return &Output{NilFlag, helpStr}, nil
			} else {
				// help on a command that does not exist
				return nil, commands.UnknownCommandError{commands.Command{arg, nil}}
			}
		}
		// just list all existing commands
		helpStr := "Available commands: " + strings.Join(CommandList(), ", ") + strings.Join(c.CommandList(), ", ")
		for _, l := range c.Links {
			helpStr += ", " + l.Name
		}
		return &Output{NilFlag, helpStr}, nil
	} else if IsBuiltin(cmd) {
		return &Output{QuitFlag, builtins[cmd.Cmd].Action(c, cmd.Args)}, nil
	}
	return nil, commands.UnknownCommandError{*cmd}
}

// Parse returns a commands.Command type from a line string
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
