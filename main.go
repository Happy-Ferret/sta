package main

import (
	"bufio"
	"fmt"
	"github.com/ribacq/sta/context"
	"github.com/ribacq/sta/games"
	"github.com/ribacq/sta/parser"
	"os"
	"strings"
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
	fmt.Fprintf(game.Writer, "Hello " + game.Name + "!\n\n")
	fmt.Fprintf(game.Writer, game.Context.Look() + "\n")
	game.Writer.Flush()
mainLoop:
	for {
		line, err := prompt(game)
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
				fmt.Fprintf(game.Writer, err.Error()+"\n")
				game.Writer.Flush()
				continue
			}
			fmt.Fprintf(game.Writer, out.Message+"\n")
			switch out.Flag {
			case parser.QuitFlag:
				game.Writer.Flush()
				break mainLoop
			}
		} else if _, err := game.Context.GetLink(cmd.Cmd) ; err == nil {
			game.UseLink(cmd.Cmd)
			fmt.Fprintf(game.Writer, game.Context.Look() + "\n")
		} else if game.Context.HasCommand(cmd) {
			// command is from context
			str, err := game.Context.ExecCommand(cmd)
			if err != nil {
				fmt.Fprintf(game.Writer, err.Error()+"\n")
				game.Writer.Flush()
				continue
			}
			fmt.Fprintf(game.Writer, str+"\n")
		} else {
			// other, command does not exist
			fmt.Fprintf(game.Writer, "Unknown command: %v %q\n", cmd.Cmd, cmd.Args)
		}
		game.Writer.Flush()
	}
}

//prompt user for a line
func prompt(game *games.Game) (string, error) {
	fmt.Fprintf(game.Writer, "\n" + game.Context.Name+" > ")
	err := game.Writer.Flush()
	if err != nil {
		return "", err
	}
	line, err := game.Reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	fmt.Fprintf(game.Writer, "\n")
	err = game.Writer.Flush()
	if err != nil {
		return "", err
	}
	return strings.Trim(line, " \n"), nil
}
