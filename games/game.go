package games

import (
	"bufio"
	"github.com/ribacq/sta/context"
)

// Game type with name of player and current context
type Game struct {
	Name    string
	Context *context.Context
	Reader  *bufio.Reader
	Writer  *bufio.Writer
}

// New returns a new game with given player name and current context
func New(name string, context *context.Context, reader *bufio.Reader, writer *bufio.Writer) *Game {
	return &Game{name, context, reader, writer}
}

// UseLink change current context using link with given name in current context
func (g *Game) UseLink(name string) error {
	l, err := g.Context.GetLink(name)
	if err != nil {
		return err
	}

	g.Context = l.GetTarget()
	return nil
}
