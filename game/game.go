/*
Package games provides type Game which represent data for a player’s current game.
*/
package game

import (
	"github.com/ribacq/sta/context"
	"regexp"
	"strings"
)

// Game type with name of player and current context
type Game struct {
	Player   *context.Context
	Commands map[string]commandFunc
	Quit     chan bool
}

// New returns a new game with given player name, current context and an empty bag.
func New(name string, ctx *context.Context) *Game {
	player := context.NewPlayer(name, ctx)
	g := &Game{
		Player: player,
		Commands: map[string]commandFunc{
			"help": help,
			"quit": quit,
			"me":   me,
		},
		Quit: make(chan bool),
	}
	g.Player.Container.Contents = append(g.Player.Container.Contents, player)
	return g
}

// ExecCommand executes a command provided as a string.
func (g *Game) Exec(cmd string) error {
	// first we’ll remove any excessive blank space
	cmd = strings.TrimSpace(cmd)

	// exit on empty command
	if len(cmd) == 0 {
		return nil
	}

	// then we’ll push seperate non blank args into a slice
	var args []string
	for _, arg := range strings.Split(cmd, " ") {
		if arg != "" {
			args = append(args, arg)
		}
	}

	// now let’s execute the command
	if l, err := g.Player.Container.GetLink(args[0]); err == nil {
		// link to another context
		if l.Locked() {
			g.Player.OutCH <- "!|You cannot go this way."
			return nil
		}
		for i, ctx := range g.Player.Container.Contents {
			if ctx == g.Player {
				g.Player.Container.Contents = append(g.Player.Container.Contents[:i], g.Player.Container.Contents[i+1:]...)
				break
			}
		}
		g.Player.Container.EventsCH <- context.Event{g.Player, context.CharacterDoesEvent, "leaves"}
		g.Player.Container = l.Target()
		g.Player.Container.EventsCH <- context.Event{g.Player, context.CharacterDoesEvent, "comes this way"}
		g.Player.Container.Contents = append(g.Player.Container.Contents, g.Player)
		g.Player.Container = g.Player.Container
		return context.Look(g.Player.Container, g.Player, args)
	} else if _, ok := g.Player.Container.HasCommand(args[0]); ok {
		// context command
		return g.Player.Container.Exec(g.Player, args)
	} else if command, ok := g.HasCommand(args[0]); ok {
		// game command
		return g.Commands[command](g, args)
	}
	// error: command not found
	g.Player.OutCH <- "!|Command ‘" + args[0] + "’ is not allowed."
	return nil
}

// HasCommand returns whether a command exists in the Game variable with no ambiguity.
func (g *Game) HasCommand(cmd string) (command string, ok bool) {
	for testedCommand := range g.Commands {
		if matched, err := regexp.Match("^"+cmd+".*$", []byte(testedCommand)); err == nil && matched {
			if ok {
				return "", false
			}
			command = testedCommand
			ok = true
		}
	}
	return
}

// AllCommands returns a slice of all currently accessible commands.
func (g *Game) AllCommands() (cmds []string) {
	for cmd := range g.Commands {
		cmds = append(cmds, cmd)
	}
	for cmd := range g.Player.Container.Commands {
		cmds = append(cmds, cmd)
	}
	for _, l := range g.Player.Container.Links {
		cmds = append(cmds, l.Name())
	}
	return
}
