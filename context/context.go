/*
Package context provides type Context which represents a container for a lookable room or item.
*/
package context

import (
	"errors"
	"regexp"
)

// Context is the type representing the current player’s position
type Context struct {
	Name           string
	Description    string
	commandActions map[string]commandFunc
	Container      *Context
	Contents       []*Context
	Links          []Link
	Properties     map[string]interface{}
}

// New returns a default context intialized with just a name and a look command.
func New(name string) *Context {
	c := &Context{
		Name: name,
		commandActions: map[string]commandFunc{
			"look": Look,
			"take": Take,
		},
		Properties: map[string]interface{}{
			"lookable": true,
		},
	}
	return c
}

// Exec executes a context command.
func (c *Context) Exec(player *Context, cmd []string) (out string, err error) {
	if f, ok := c.commandActions[cmd[0]]; ok {
		return f(c, player, cmd)
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

// Pick finds a *Context from a Context.Contents by name without modifying the slice.
// Returns i the context index in the Contents slice,
// ctx the found *Context,
// ok a boolean indicating whether the contex was found
func (c *Context) Pick(name string) (i int, ctx *Context, ok bool) {
	// looks through all of c.Contents if we find the rightly named context
	for i, ctx := range c.Contents {
		ok, err := regexp.Match(".*"+name+".*", []byte(ctx.Name))
		if err != nil || !ok {
			return 0, nil, false
		}
		return i, ctx, ok
	}
	// name was not found
	return 0, nil, false
}
