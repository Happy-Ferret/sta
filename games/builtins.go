package games

import (
	"github.com/ribacq/sta/context"
	"strings"
)

func help(g *Game, cmd []string) (out string, err error) {
	// err is always nil

	// if help is asked on a specific topic
	if len(cmd) > 1 {
		out = "Help on ‘" + strings.Join(cmd[1:], " ") + "’."
		return
	}

	// print all available commands
	out = "Available commands: "
	i := 0
	for cmd := range g.Commands {
		if i > 0 {
			out += ", "
		}
		i++
		out += cmd
	}
	for cmd := range g.Context.Commands {
		if i > 0 {
			out += ", "
		}
		i++
		out += cmd
	}
	for _, l := range g.Context.Links {
		if i > 0 {
			out += ", "
		}
		i++
		out += l.Name()
	}
	return
}

func quit(g *Game, cmd []string) (out string, err error) {
	g.quit = true
	return "Fare well...", nil
}

func me(g *Game, cmd []string) (out string, err error) {
	out, err = context.Look(g.Player, g.Player, cmd)
	out = "You are looking at yourself.\n" + out
	return
}
