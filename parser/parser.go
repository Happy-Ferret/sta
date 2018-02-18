package parser

import (
	"github.com/ribacq/sta/commands"
	"github.com/ribacq/sta/context"
	"strings"
)

// builtin commands
var (
	builtins = map[string]context.CommandAction {
		"quit": context.CommandAction {
			func (c *context.Context, args []string) string {
				return "Fare well…"
			},
			"This command allows you to exit the game."},
		"help": context.CommandAction {
			func (c *context.Context, args []string) string {
				return "help…"
			},
			"This command gives you the list of available command or help for a specific command.",
		},
	}
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

// get all builtin commands
func CommandList() []string {
	cmds := make([]string, 0, len(builtins))
	for cmd := range builtins {
		cmds = append(cmds, cmd)
	}
	return cmds
}

// get whether a command is a builtin
func IsBuiltin(cmd *commands.Command) bool {
	for c := range builtins {
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
	if (cmd.Cmd == "help") {
		if len(cmd.Args) > 0 {
			arg := cmd.Args[0]
			if IsBuiltin(&commands.Command{arg, nil}) {
				helpStr := "Built-in: " + builtins[arg].Help
				return &ParserOutput{NilFlag, helpStr}, nil
			} else if c.HasCommand(&commands.Command{arg, nil}) {
				helpStr := c.CommandActions[arg].Help
				return &ParserOutput{NilFlag, helpStr}, nil
			} else {
				return nil, commands.UnknownCommandError{commands.Command{arg, nil}}
			}
		}
		helpStr := "Available commands: " + strings.Join(CommandList(), ", ") + strings.Join(c.CommandList(), ", ")
		return &ParserOutput{NilFlag, helpStr}, nil
	} else if IsBuiltin(cmd) {
		return &ParserOutput{QuitFlag, builtins[cmd.Cmd].Action(c, cmd.Args)}, nil
	}
	return nil, commands.UnknownCommandError{*cmd}
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
