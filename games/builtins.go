package games

import (
	"github.com/ribacq/sta/context"
)

func help(g *Game, cmd []string) (out string, err error) {
	return "help...", nil
}

func quit(g *Game, cmd []string) (out string, err error) {
	g.quit = true
	return "quit...", nil
}

func me(g *Game, cmd []string) (out string, err error) {
	out, err = context.Look(g.Player, g.Player, cmd)
	out = "You are looking at yourself.\n" + out
	return
}
