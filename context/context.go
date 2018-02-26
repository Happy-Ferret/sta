package context

import (
	"errors"
)

type commandFunc func(c *Context, cmd []string) (out string, err error)

// Context is the type representing the current player’s position
type Context struct {
	Name           string
	Description    string
	commandActions map[string]commandFunc
	Container      *Context
	Links          []Link
}

// New returns a default context intialized with just a name and a look command.
func New(name string) *Context {
	c := &Context{
		Name: name,
		commandActions: map[string]commandFunc{
			"look": Look,
		},
	}
	return c
}

// MakeTakeable adds ‘take’ command to a context.
func (c *Context) MakeTakeable() {
	c.commandActions["take"] = Take
}

// Exec executes a context command.
func (c *Context) Exec(cmd []string) (out string, err error) {
	if f, ok := c.commandActions[cmd[0]]; ok {
		return f(c, cmd)
	}
	return "", errors.New("c.Exec: no such command " + cmd[0])
}

// CommandList extracts the command list from the commandActions map.
func (c *Context) CommandList() []string {
	var cmds []string
	for str := range c.commandActions {
		cmds = append(cmds, str)
	}
	return cmds
}

// HasCommand gets whether the command exists in given context.
func (c *Context) HasCommand(cmd string) bool {
	for str := range c.commandActions {
		if cmd == str {
			return true
		}
	}
	return false
}
