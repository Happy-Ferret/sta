package context

import (
	"github.com/ribacq/sta/commands"
)

// CommandAction regroups a function with a help string for a command
type CommandAction struct {
	Action func(c *Context, args []string) string
	Help   string
}

// Context is the type representing the current player’s position
type Context struct {
	Name           string
	Description    string
	CommandActions map[string]CommandAction
	Container      *Context
	Links          []Link
}

// New returns a default context intialized with just a name and a look command
func New(name string) *Context {
	var c Context
	c.Name = name
	c.Description = ""
	c.Container = nil
	c.CommandActions = make(map[string]CommandAction)
	c.CommandActions["look"] = CommandAction{
		func(c *Context, args []string) string {
			return c.Look()
		},
		"This command allows you to look around you.",
	}
	return &c
}

// MakeTakeable adds ‘take’ command to a context
func (c *Context) MakeTakeable() {
	c.CommandActions["take"] = CommandAction{
		func(c *Context, args []string) string {
			return "You took " + c.Name + "!"
		},
		"This command allows you to take an item from the environment.",
	}
}

// CommandList extracts the command list from the CommandActions map
func (c *Context) CommandList() []string {
	cmds := make([]string, len(c.CommandActions))
	for str := range c.CommandActions {
		cmds = append(cmds, str)
	}
	return cmds
}

// HasCommand gets whether the command exists in given context
func (c *Context) HasCommand(cmd *commands.Command) bool {
	for str := range c.CommandActions {
		if cmd.Cmd == str {
			return true
		}
	}
	return false
}

// Look command gives a description of the context and available un-hidden links
func (c *Context) Look() string {
	lookStr := c.Description + "\nLinks:"
	for _, l := range c.Links {
		lookStr += " " + l.Name
	}
	return lookStr
}
