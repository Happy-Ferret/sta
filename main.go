package main

import (
	"bufio"
	"fmt"
	"github.com/ribacq/sta/context"
	"github.com/ribacq/sta/games"
	"github.com/ribacq/sta/parser"
	"os"
)

func main() {
	// default rooms
	// TODO: use a database
	hall := context.New("Palace hall")
	hall.Description = "You are inside the hall of a palace. There is a door to the left."
	kitchen := context.New("Kitchen")
	kitchen.Description = "You are inside a shiny kitchen. There is a camembert on the table and a door to your right."

	hall.AddLink(kitchen, "door")
	kitchen.AddLink(hall, "door")

	// Game variable with name, context, reader and writer
	game := games.New("Jirsad", hall, bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout))
	game.Writeln("Hello" + game.Name)
	game.Writeln("\n" + game.Context.Look())
mainLoop:
	for {
		line, err := game.Prompt()
		if err != nil {
			return
		}
		cmd := parser.Parse(line)

		if cmd.Cmd == "" {
			// command was nothing, skip
			continue
		} else if parser.IsBuiltin(cmd) {
			// command is a builtin, call parser
			out, err := parser.ExecBuiltin(cmd, game.Context)
			if err != nil {
				game.Writeln(err.Error())
				continue
			}
			game.Writeln(out.Message)
			switch out.Flag {
			case parser.QuitFlag:
				break mainLoop
			}
		} else if _, err := game.Context.GetLink(cmd.Cmd); err == nil {
			game.UseLink(cmd.Cmd)
			game.Writeln(game.Context.Look())
		} else if game.Context.HasCommand(cmd) {
			// command is from context
			str, err := game.Context.ExecCommand(cmd)
			if err != nil {
				game.Writeln(err.Error())
				continue
			}
			game.Writeln(str)
		} else {
			// other, command does not exist
			game.Writeln(fmt.Sprintf("Unknown command: %v %q\n", cmd.Cmd, cmd.Args))
		}
	}
}
