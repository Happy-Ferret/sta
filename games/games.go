package games

import (
	"errors"
	"github.com/ribacq/sta/context"
	"strings"
)

type commandFunc func(g *Game, cmd []string) (out string, err error)

// Game type with name of player and current context
type Game struct {
	Name           string
	Context        *context.Context
	Bag            []*context.Context
	commandActions map[string]commandFunc
	quit           bool
}

// New returns a new game with given player name, current context and an empty bag.
func New(name string, ctx *context.Context) *Game {
	game := &Game{
		Name:    name,
		Context: ctx,
		commandActions: map[string]commandFunc{
			"help": helpCF,
			"quit": quitCF,
		},
		quit: false,
	}
	return game
}

// Quit returns if we want to quit the game
func (g Game) Quit() bool {
	return g.quit
}

// ExecCommand executes a command provided as an array of strings.
func (g *Game) Exec(cmd []string) (out string, err error) {
	if l, err := g.Context.GetLink(cmd[0]); err == nil {
		g.Context = l.Target()
		return context.Look(g.Context, cmd)
	} else if g.Context.HasCommand(cmd[0]) {
		return g.Context.Exec(cmd)
	} else if g.HasCommand(cmd[0]) {
		return g.commandActions[cmd[0]](g, cmd)
	}
	return "", errors.New("Command " + cmd[0] + " does not exist.")
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

// Take places context in bag
func (g *Game) Take(ctx *context.Context) {
	if ctx == nil {
		return
	}

	ctx.Container = g.Context
	g.Bag = append(g.Bag, ctx)
}

// ShowBag describes bag content
func (g *Game) ShowBag() string {
	bagNames := make([]string, len(g.Bag))
	for _, ctx := range g.Bag {
		bagNames = append(bagNames, ctx.Name)
	}
	return "Bag: " + strings.Join(bagNames, ", ")
}
