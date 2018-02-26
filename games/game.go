package games

import (
	"github.com/ribacq/sta/context"
	"golang.org/x/crypto/ssh/terminal"
)

// Game type with name of player and current context
type Game struct {
	Name    string
	Context *context.Context
	*terminal.Terminal
}

// New returns a new game with given player name and current context
func New(name string, context *context.Context, term *terminal.Terminal) *Game {
	return &Game{name, context, term}
}

// Prompt user for a line
func (g *Game) Prompt() (string, error) {
	g.SetPrompt(g.Context.Name + " > ")
	line, err := g.ReadLine()
	if err != nil {
		return "", err
	}
	return line, nil
}

func (g *Game) WriteString(str string) (int, error) {
	return g.Write([]byte(str))
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
