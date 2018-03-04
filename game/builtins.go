package game

import (
	"github.com/ribacq/sta/context"
	"strings"
)

type commandFunc func(g *Game, cmd []string) error

func help(g *Game, cmd []string) error {
	// err is always nil

	// if help is asked on a specific topic
	if len(cmd) > 1 {
		g.Player.OutCH <- "Help on ‘" + strings.Join(cmd[1:], " ") + "’."
		return nil
	}

	// print all available commands
	g.Player.OutCH <- "Available commands: " + strings.Join(g.AllCommands(), ", ")
	return nil
}

func quit(g *Game, cmd []string) error {
	g.Quit <- true
	g.Player.OutCH <- "Fare well..."
	return nil
}

func me(g *Game, cmd []string) error {
	g.Player.OutCH <- "You are looking at yourself."
	return context.Look(g.Player, g.Player, cmd)
}
