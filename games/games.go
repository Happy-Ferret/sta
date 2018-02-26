/*
Package games provides type Game which represent data for a single player’s current game.
*/
package games

import (
	"errors"
	"github.com/ribacq/sta/context"
	"regexp"
	"strings"
)

type commandFunc func(g *Game, cmd []string) (out string, err error)

// Game type with name of player and current context
type Game struct {
	Player         *context.Context
	Context        *context.Context
	CommandActions map[string]commandFunc
	quit           bool
}

// New returns a new game with given player name, current context and an empty bag.
func New(name string, ctx *context.Context) *Game {
	game := &Game{
		Player:  context.NewPlayer(name),
		Context: ctx,
		CommandActions: map[string]commandFunc{
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

	// exit on empty command
	if len(cmd) == 0 {
		return "", nil
	}

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
		if target, ok := l.Try(g.Player); ok {
			g.Context = target
			return context.Look(g.Context, g.Player, args)
		}
		return "", errors.New("You do not have the required key.")
	} else if _, ok := g.Context.HasCommand(args[0]); ok {
		// context command
		return g.Context.Exec(g.Player, args)
	} else if command, ok := g.HasCommand(args[0]); ok {
		// game command
		return g.CommandActions[command](g, args)
	}
	// error: command not found
	return "", errors.New("Command ‘" + args[0] + "’ is not allowed.")
}

// HasCommand returns whether a command exists in the Game variable.
func (g *Game) HasCommand(cmd string) (command string, ok bool) {
	for command := range g.CommandActions {
		if matched, err := regexp.Match("^"+cmd+".*$", []byte(command)); err == nil && matched {
			return command, true
		}
	}
	return "", false
}
