package context

import (
	"errors"
	"regexp"
	"strings"
)

type commandFunc func(c, player *Context, cmd []string) (out string, err error)

// Look command gives a description of the context and available un-hidden links.
func Look(c, player *Context, cmd []string) (out string, err error) {
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
func Take(from, to *Context, cmd []string) (out string, err error) {
	if i, ctx, ok := from.Pick(strings.Join(cmd[1:], " ")); ok {
		if _, ok := ctx.Properties["takeable"]; !ok {
			return "You cannot take " + ctx.Name + ".", nil
		}
		from.Contents = append(from.Contents[0:i], from.Contents[i+1:]...)
		to.Contents = append(to.Contents, ctx)
		ctx.Container = to
		return "You take " + ctx.Name + " from " + from.Name + ".", nil
	}
	return "", errors.New("There is no such thing here.")
}

// Lock locks or unlocks a link given in argument if the accurate key is owned
// Command called with ‘lock’ or ‘unlock’.
func Lock(c, player *Context, cmd []string) (out string, err error) {
	// get whether we want to lock or unlock
	action := true
	if matched, err := regexp.Match("^"+cmd[0]+".*", []byte("unlock")); err == nil && matched {
		action = false
	}

	// get link or return if not found
	l, err := c.GetLink(strings.Join(cmd[1:], " "))
	if err != nil {
		return "", errors.New("No such link found.")
	}

	// return if link has no key
	if !l.HasKey() {
		return "", errors.New("The link has no lock.")
	}

	// try to lock or unlock now
	for _, ctx := range player.Contents {
		if val, ok := ctx.Properties["key"]; ok && val == l.key {
			// return if there is nothing to do
			if l.locked == action {
				if action {
					out = l.Name + " is already locked."
				} else {
					out = l.Name + " is already unlocked."
				}
				return
			}

			// lock or unlock
			l.locked = action
			if action {
				out = l.Name + " locked!"
			} else {
				out = l.Name + " unlocked!"
			}
			return
		}
	}

	// return if required key was not found in player.Contents
	return "", errors.New("Required key not found.")
}
