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
	name        string
	description string
	container   *Context `json:"-"`
	contents    []*Context
	links       []*Link
	commands    map[string]CommandFunc
	Properties  map[string]string
	EventsCH    chan Event  `json:"-"`
	OutCH       chan string `json:"-"`
}

// New returns a default context intialized with just a name and a look command.
func New(name string) *Context {
	c := &Context{
		name:     name,
		commands: commandFuncs,
		Properties: map[string]string{
			"lookable": "lookable",
		},
		EventsCH: make(chan Event),
		OutCH:    make(chan string, 256),
	}

	go c.handleEvents()

	return c
}

//////
// Getters and setters
//////

// Name is a getter for context name.
func (c *Context) Name() string {
	return c.name
}

// SetName updates a context’s name.
func (c *Context) SetName(name string) {
	c.name = name
}

// Description is a getter for context description.
func (c *Context) Description() string {
	return c.description
}

// SetDescription updates a context’s description.
func (c *Context) SetDescription(description string) {
	c.description = description
}

// Container is a getter for a context’s container.
func (c *Context) Container() *Context {
	return c.container
}

// SetContainer is a setter for a context’s container.
func (c *Context) SetContainer(container *Context) {
	c.container = container
}

// Contents is a getter for a context’s contents.
func (c *Context) Contents() []*Context {
	return c.contents
}

// Content is a getter for an item in a context’s contents.
func (c *Context) Content(i int) *Context {
	if i < 0 || i >= len(c.contents) {
		return nil
	}
	return c.contents[i]
}

// AppendContent appends to a context’s contents.
func (c *Context) AppendContent(ctxs ...*Context) {
	c.contents = append(c.contents, ctxs...)
}

// RemoveContent removes a context from another, based on it’s index in the contents slice.
func (c *Context) RemoveContent(i int) {
	if i < 0 || i >= len(c.contents) {
		return
	}
	c.contents = append(c.contents[:i], c.contents[i+1:]...)
}

// RemoveFromContainer removes a context from its container.
func (c *Context) RemoveFromContainer() bool {
	if c.container == nil {
		return false
	}
	for i, ctx := range c.container.Contents() {
		if ctx == c {
			c.container.RemoveContent(i)
			return true
		}
	}
	return false
}

// Links is a getter for a context’s links.
func (c *Context) Links() []*Link {
	return c.links
}

// CommandNames returns the names of available commands for this context.
func (c *Context) CommandNames() (names []string) {
	for cmd := range c.commands {
		names = append(names, cmd)
	}
	return
}

//////
// Other actions
//////

// Exec executes a context command.
func (c *Context) Exec(player *Context, args []string) error {
	if command, ok := c.HasCommand(args[0]); ok {
		return c.commands[command](c, player, args)
	}
	return errors.New("c.Exec: no such command " + args[0])
}

// HasCommand gets whether the command exists in given context with no ambiguity.
func (c *Context) HasCommand(cmd string) (command string, ok bool) {
	for testedCommand := range c.commands {
		if matched, err := regexp.Match("^"+cmd+".*$", []byte(testedCommand)); err == nil && matched {
			if ok {
				return "", false
			}
			command = testedCommand
			ok = true
		}
	}
	return
}

// Pick finds a *Context from a Context.contents by name without modifying the slice.
// Returns i the context index in the contents slice,
// ctx the found *Context,
// ok a boolean indicating whether the contex was found
func (c *Context) Pick(name string) (i int, ctx *Context, ok bool) {
	// trim spaces and return on empty name
	name = strings.TrimSpace(name)
	if name == "" {
		return
	}

	// looks through all of c.contents if we find the rightly named context
	for i, ctx := range c.contents {
		ok, err := regexp.Match(".*"+name+".*", []byte(ctx.name))
		if err == nil && ok {
			return i, ctx, ok
		}
	}

	// name was not found
	return
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
