package context

import (
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
		"say":    Say,
	}
)

type CommandFunc func(c, player *Context, cmd []string) error

// GetCommandFunc returns the function corresponding to the command.
// The command name must be exact.
// This function exists to protect commandFuncs from exterior modification.
func GetCommandFunc(cmd string) (f CommandFunc, ok bool) {
	f, ok = commandFuncs[cmd]
	return
}

// Look command gives a description of the context and available un-hidden links.
func Look(c, player *Context, cmd []string) error {
	// maybe look for something else than c
	if len(cmd) > 1 {
		_, ctx, ok := c.Pick(strings.Join(cmd[1:], " "))
		if !ok {
			player.OutCH <- "!|No such thing here."
			return nil
		} else {
			return Look(ctx, player, []string{})
		}
	}

	// notify the looked context
	c.EventsCH <- Event{player, LookEvent, c}

	// display context description
	out := c.Description()

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
				out += "*" + item.Name() + "*"
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

	player.OutCH <- out
	return nil
}

// Take puts an item into or out of the  player’s bag.
func Take(c, player *Context, cmd []string) error {
	// set to and from depending on calling command (default is take)
	drop := false
	from, to := c, player
	if matched, err := regexp.Match("^"+cmd[0]+".*", []byte("drop")); err == nil && matched {
		from, to = to, from
		drop = true
	}

	// transfer the object
	if i, ctx, ok := from.Pick(strings.Join(cmd[1:], " ")); ok {
		if _, ok := ctx.Properties["takeable"]; !ok {
			if drop {
				player.OutCH <- "!|You cannot drop " + ctx.Name() + "."
			} else {
				player.OutCH <- "!|You cannot take " + ctx.Name() + "."
			}
			return nil
		}
		from.Contents = append(from.Contents[0:i], from.Contents[i+1:]...)
		to.Contents = append(to.Contents, ctx)
		ctx.SetContainer(to)
		c.EventsCH <- Event{player, TakeDropEvent, takeDropEventContent{ctx, drop}}
		return nil
	}
	player.OutCH <- "!|There is no such thing here."
	return nil
}

// Lock locks or unlocks a link given in argument if the accurate key is owned
// Command called with ‘lock’ or ‘unlock’.
func Lock(c, player *Context, cmd []string) error {
	// get whether we want to lock or unlock
	action := true
	if matched, err := regexp.Match("^"+cmd[0]+".*", []byte("unlock")); err == nil && matched {
		action = false
	}

	// get link or return if not found
	l, err := c.GetLink(strings.Join(cmd[1:], " "))
	if err != nil {
		player.OutCH <- "!|No such link found."
		return nil
	}

	// return if link has no key
	if !l.HasKey() {
		player.OutCH <- "The link has no lock."
		return nil
	}

	// try to lock or unlock now
	for _, ctx := range player.Contents {
		if val, ok := ctx.Properties["key"]; ok && val == l.key {
			// return if there is nothing to do
			if l.locked == action {
				if action {
					player.OutCH <- "!|" + l.Name() + " is already locked."
				} else {
					player.OutCH <- "!|" + l.Name() + " is already unlocked."
				}
				return nil
			}

			// lock or unlock l and its slaves
			l.locked = action
			for _, slave := range l.slaves {
				slave.locked = action
			}
			if action {
				c.OutCH <- l.Name() + " locked!"
			} else {
				c.OutCH <- l.Name() + " unlocked!"
			}
			return nil
		}
	}

	// return if required key was not found in player.Contents
	player.OutCH <- "!|Required key not found."
	return nil
}

// Say writes a message publicly to the current context.
func Say(c, player *Context, cmd []string) error {
	if len(cmd) < 2 {
		player.OutCH <- "!|You have to say something..."
		return nil
	}
	c.OutCH <- "*" + player.Name() + "*: `" + strings.Join(cmd[1:], " ") + "`"
	return nil
}
