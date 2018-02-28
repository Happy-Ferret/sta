package games

import (
	"github.com/ribacq/sta/context"
	"strings"
)

type commandFunc func(g *Game, cmd []string) (out string, err error)

func help(g *Game, cmd []string) (out string, err error) {
	// err is always nil

	// if help is asked on a specific topic
	if len(cmd) > 1 {
		out = "Help on ‘" + strings.Join(cmd[1:], " ") + "’."
		return
	}

	// print all available commands
	out = "Available commands: " + strings.Join(g.AllCommands(), ", ")
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
