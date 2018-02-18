package context

import (
	"github.com/ribacq/sta/parser"
)

// types
type Context struct {
	Name string
	CommandActions map[string]func (args []string) string
}

// a default context
func Default() *Context {
	var c Context
	c.Name = "the palace"
	c.CommandActions = make(map[string]func (args []string) string)
	c.CommandActions["look"] = func (args []string) string {
		return "This is a wonderful palace."
	}
	return &c
}

// whether the command exists in given context
func (c *Context) HasCommand(cmd *parser.Command) bool {
	for str, _ := range c.CommandActions {
		if cmd.Cmd == str {
			return true
		}
	}
	return false
}

// execute a command
func (c *Context) ExecCommand(cmd *parser.Command) (string, error) {
	if !c.HasCommand(cmd) {
		return "", parser.UnknownCommandError{*cmd}
	}
	return c.CommandActions[cmd.Cmd](cmd.Args), nil
}
