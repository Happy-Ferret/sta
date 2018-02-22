package games

import (
	"bufio"
	"fmt"
	"github.com/ribacq/sta/commands"
	"github.com/ribacq/sta/context"
	"strings"
)

// Game type with name of player and current context
type Game struct {
	Name    string
	Context *context.Context
	Bag     []*context.Context
	reader  *bufio.Reader
	writer  *bufio.Writer
}

// New returns a new game with given player name and current context
func New(name string, ctx *context.Context, reader *bufio.Reader, writer *bufio.Writer) *Game {
	return &Game{name, ctx, make([]*context.Context, 16), reader, writer}
}

// Write writes a string to game.Writer and flushes the output
func (g *Game) Write(str string) error {
	if _, err := fmt.Fprintf(g.writer, str); err != nil {
		return err
	}
	return g.writer.Flush()
}

// Writeln writes a string to game.writer with a newline at the end
func (g *Game) Writeln(str string) error {
	return g.Write(str + "\n")
}

// Prompt user for a line
func (g *Game) Prompt() (string, error) {
	fmt.Fprintf(g.writer, "\n"+g.Context.Name+" > ")
	err := g.writer.Flush()
	if err != nil {
		return "", err
	}
	line, err := g.reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	fmt.Fprintf(g.writer, "\n")
	err = g.writer.Flush()
	if err != nil {
		return "", err
	}
	return strings.Trim(line, " \n"), nil
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
