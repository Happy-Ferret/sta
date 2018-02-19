package context

import (
	"github.com/ribacq/sta/commands"
)

// types
type CommandAction struct {
	Action func (c *Context, args []string) string
	Help string
}

type Context struct {
	Name string
	Description string
	CommandActions map[string]CommandAction
	Links []Link
}

// a default context
func New(name string) *Context {
	var c Context
	c.Name = name
	c.Description = ""
	c.CommandActions = make(map[string]CommandAction)
	c.CommandActions["look"] = CommandAction{
		func (c *Context, args []string) string {
			return c.Look()
		},
		"This command allows you to look around you.",
	}
	return &c
}

// extracts the command list from the CommandActions map
func (c *Context) CommandList() []string {
	cmds := make([]string, len(c.CommandActions))
	for str := range c.CommandActions {
		cmds = append(cmds, str)
	}
	return cmds
}

// whether the command exists in given context
func (c *Context) HasCommand(cmd *commands.Command) bool {
	for str, _ := range c.CommandActions {
		if cmd.Cmd == str {
			return true
		}
	}
	return false
}

// execute a command
func (c *Context) ExecCommand(cmd *commands.Command) (string, error) {
	if !c.HasCommand(cmd) {
		return "", commands.UnknownCommandError{*cmd}
	}
	return c.CommandActions[cmd.Cmd].Action(c, cmd.Args), nil
}

// look command
func (c *Context) Look() string {
	lookStr := c.Description + "\nLinks:"
	for _, l := range c.Links {
		lookStr += " " + l.Name
	}
	return lookStr
}
