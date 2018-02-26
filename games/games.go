/*
Package games provides type Game which represent data for a single player’s current game.
*/
package games

import (
	"errors"
	"github.com/ribacq/sta/context"
	"strings"
)

type commandFunc func(g *Game, cmd []string) (out string, err error)

// Game type with name of player and current context
type Game struct {
	Player         *context.Context
	Context        *context.Context
	commandActions map[string]commandFunc
	quit           bool
}

// New returns a new game with given player name, current context and an empty bag.
func New(name string, ctx *context.Context) *Game {
	game := &Game{
		Player:  context.NewPlayer(name),
		Context: ctx,
		commandActions: map[string]commandFunc{
			"help": help,
			"quit": quit,
			"me":   me,
		},
		quit: false,
	}
	return game
}

// Quit returns if we want to quit the game
func (g Game) Quit() bool {
	return g.quit
}

// ExecCommand executes a command provided as a string.
func (g *Game) Exec(cmd string) (out string, err error) {
	// first we’ll remove any excessive blank space
	cmd = strings.TrimSpace(cmd)
	// then we’ll push seperate non blank args into a slice
	var args []string
	for _, arg := range strings.Split(cmd, " ") {
		if arg != "" {
			args = append(args, arg)
		}
	}

	// now let’s execute the command
	if l, err := g.Context.GetLink(args[0]); err == nil {
		// link to another content
		g.Context = l.Target()
		return context.Look(g.Context, g.Player, args)
	} else if g.Context.HasCommand(args[0]) {
		// context command
		return g.Context.Exec(g.Player, args)
	} else if g.HasCommand(args[0]) {
		// game command
		return g.commandActions[args[0]](g, args)
	}
	// error: command not found
	return "", errors.New("Command ‘" + args[0] + "’ is not allowed.")
}

// HasCommand returns whether a command exists in the Game variable.
func (g *Game) HasCommand(cmd string) bool {
	for str := range g.commandActions {
		if cmd == str {
			return true
		}
	}
	return false
}
