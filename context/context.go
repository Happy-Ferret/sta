/*
Package context provides type Context which represents a container for a lookable room or item.
*/
package context

import (
	"errors"
	"regexp"
	"strings"
)

// Context is the type representing the current player’s position
type Context struct {
	Name        string
	Description string
	Container   *Context
	Contents    []*Context
	Links       []*Link
	Commands    map[string]CommandFunc
	Properties  map[string]string
}

// New returns a default context intialized with just a name and a look command.
func New(name string) *Context {
	c := &Context{
		Name:     name,
		Commands: commandFuncs,
		Properties: map[string]string{
			"lookable": "lookable",
		},
	}
	return c
}

// Exec executes a context command.
func (c *Context) Exec(player *Context, args []string) (out string, err error) {
	if command, ok := c.HasCommand(args[0]); ok {
		return c.Commands[command](c, player, args)
	}
	return "", errors.New("c.Exec: no such command " + args[0])
}

// HasCommand gets whether the command exists in given context.
func (c *Context) HasCommand(cmd string) (command string, ok bool) {
	for command := range c.Commands {
		if matched, err := regexp.Match("^"+cmd+".*$", []byte(command)); err == nil && matched {
			return command, true
		}
	}
	return "", false
}

// Pick finds a *Context from a Context.Contents by name without modifying the slice.
// Returns i the context index in the Contents slice,
// ctx the found *Context,
// ok a boolean indicating whether the contex was found
func (c *Context) Pick(name string) (i int, ctx *Context, ok bool) {
	// looks through all of c.Contents if we find the rightly named context
	name = strings.TrimSpace(name)
	for i, ctx := range c.Contents {
		ok, err := regexp.Match(".*"+name+".*", []byte(ctx.Name))
		if err == nil && ok {
			return i, ctx, ok
		}
	}
	// name was not found
	return 0, nil, false
}

// MakeLookable allows or forbids command ‘look c’.
func (c *Context) MakeLookable(val bool) {
	if val {
		c.Properties["lookable"] = "lookable"
	} else {
		delete(c.Properties, "lookable")
	}
}

// MakeTakeable allows or forbids command ‘take c’
func (c *Context) MakeTakeable(val bool) {
	if val {
		c.Properties["takeable"] = "takeable"
	} else {
		delete(c.Properties, "takeable")
	}
}
