package games

import (
	"github.com/ribacq/sta/commands"
	"github.com/ribacq/sta/context"
	"golang.org/x/crypto/ssh/terminal"
	"strings"
)

// Game type with name of player and current context
type Game struct {
	Name    string
	Context *context.Context
	Bag     []*context.Context
	*terminal.Terminal
}

// New returns a new game with given player name and current context
func New(name string, ctx *context.Context, term *terminal.Terminal) *Game {
	return &Game{name, ctx, make([]*context.Context, 0), term}
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

// ExecCommand executes a command of the Context
func (g *Game) ExecCommand(cmd *commands.Command) (string, error) {
	if !g.Context.HasCommand(cmd) {
		return "", commands.UnknownCommandError{*cmd}
	}
	return g.Context.CommandActions[cmd.Cmd].Action(g.Context, cmd.Args), nil
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
