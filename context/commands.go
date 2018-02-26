package context

import (
	"strings"
)

type commandFunc func(c *Context, player *Context, cmd []string) (out string, err error)

// Look command gives a description of the context and available un-hidden links.
func Look(c *Context, player *Context, cmd []string) (out string, err error) {
	// err is always nil.
	// display context description
	out = c.Description

	// display context contents
	if len(c.Contents) > 0 {
		out += "\nThere is "
		for i, item := range c.Contents {
			if i > 0 && i < len(c.Contents)-1 {
				out += ", "
			} else if i > 0 && i == len(c.Contents)-1 {
				out += " and "
			}
			out += item.Name
		}
		out += "."
	}

	// display context links
	if len(c.Links) > 0 {
		out += "\nLinks: "
		for i, l := range c.Links {
			if i > 0 {
				out += ", "
			}
			out += l.Name
		}
	}

	return
}

// Take puts an item into player’s bag.
func Take(from *Context, to *Context, cmd []string) (out string, err error) {
	if i, ctx, ok := from.Pick(strings.Join(cmd[1:], " ")); ok {
		from.Contents = append(from.Contents[0:i], from.Contents[i+1:]...)
		to.Contents = append(to.Contents, ctx)
		ctx.Container = to
		return "You take " + ctx.Name + " from " + from.Name + ".", nil
	}
	return "There is no such thing here", nil
}
