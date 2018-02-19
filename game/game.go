package game

import (
	"github.com/ribacq/sta/context"
)

// type Game with name of player and current context
type Game struct {
	Name string
	Ct *context.Context
}

// returns a new game with given player name and current context
func New(name string, ct *context.Context) *Game {
	return &Game{name, ct}
}

// change current context using link with given name in current context
func (g *Game) UseLink(name string) error {
	l, err := g.Ct.GetLink(name)
	if err != nil {
		return err
	}

	g.Ct = l.GetTarget()
	return nil
}
