package context

import (
	"errors"
	"regexp"
	"strings"
)

var (
	commandFuncs = map[string]CommandFunc{
		"look":   Look,
		"take":   Take,
		"drop":   Take,
		"lock":   Lock,
		"unlock": Lock,
	}
)

type CommandFunc func(c, player *Context, cmd []string) (out string, err error)

// GetCommandFunc returns the function corresponding to the command.
// The command name must be exact.
// This function exists to protect commandFuncs from exterior modification.
func GetCommandFunc(cmd string) (f CommandFunc, ok bool) {
	f, ok = commandFuncs[cmd]
	return
}

// Look command gives a description of the context and available un-hidden links.
func Look(c, player *Context, cmd []string) (out string, err error) {
	// maybe look for something else than c
	if len(cmd) > 1 {
		_, ctx, ok := c.Pick(strings.Join(cmd[1:], " "))
		if !ok {
			return "", errors.New("No such thing here.")
		} else {
			return Look(ctx, player, []string{})
		}
	}

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
			if item == player {
				out += "you"
			} else {
				out += "*" + item.Name + "*"
			}
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
			out += "**" + l.Name() + "**"
			if l.locked {
				out += " (locked)"
			}
		}
	}

	return
}

// Take puts an item into or out of the  player’s bag.
func Take(c, player *Context, cmd []string) (out string, err error) {
	// set to and from depending on calling command (default is take)
	from, to := c, player
	if matched, err := regexp.Match("^"+cmd[0]+".*", []byte("drop")); err == nil && matched {
		from, to = to, from
	}

	// transfer the object
	if i, ctx, ok := from.Pick(strings.Join(cmd[1:], " ")); ok {
		if _, ok := ctx.Properties["takeable"]; !ok {
			return "You cannot take " + ctx.Name + ".", nil
		}
		from.Contents = append(from.Contents[0:i], from.Contents[i+1:]...)
		to.Contents = append(to.Contents, ctx)
		ctx.Container = to
		return ctx.Name + " --> " + to.Name, nil
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
					out = l.Name() + " is already locked."
				} else {
					out = l.Name() + " is already unlocked."
				}
				return
			}

			// lock or unlock l and its slaves
			l.locked = action
			for _, slave := range l.slaves {
				slave.locked = action
			}
			if action {
				out = l.Name() + " locked!"
			} else {
				out = l.Name() + " unlocked!"
			}
			return
		}
	}

	// return if required key was not found in player.Contents
	return "", errors.New("Required key not found.")
}
