package main

import (
	"bufio"
	"fmt"
	"github.com/ribacq/sta/context"
	"github.com/ribacq/sta/game"
	"github.com/ribacq/sta/parser"
	"os"
	"strings"
)

func main() {
	// setting I/O
	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)

	// creating context
	hall := context.New("Palace hall")
	hall.Description = "You are inside the hall of a palace. There is a door to the left."
	kitchen := context.New("Kitchen")
	kitchen.Description = "You are inside a shiny kitchen. There is a camembert on the table and a door to your right."

	hall.AddLink(kitchen, "door")
	kitchen.AddLink(hall, "door")

	g := game.New("Jirsad", hall)
	fmt.Fprintf(w, "Hello " + g.Name + "!\n\n")
	fmt.Fprintf(w, g.Ct.Look() + "\n")
	w.Flush()
mainLoop:
	for {
		line, err := prompt(g.Ct, r, w)
		if err != nil {
			return
		}
		cmd := parser.Parse(line)

		if cmd.Cmd == "" {
			// command was nothing, skip
			continue
		} else if parser.IsBuiltin(cmd) {
			// command is a builtin, call parser
			out, err := parser.ExecBuiltin(cmd, g.Ct)
			if err != nil {
				fmt.Fprintf(w, err.Error()+"\n")
				w.Flush()
				continue
			}
			fmt.Fprintf(w, out.Message+"\n")
			switch out.Flag {
			case parser.QuitFlag:
				w.Flush()
				break mainLoop
			}
		} else if _, err := g.Ct.GetLink(cmd.Cmd) ; err == nil {
			g.UseLink(cmd.Cmd)
			fmt.Fprintf(w, g.Ct.Look() + "\n")
		} else if g.Ct.HasCommand(cmd) {
			// command is from context
			str, err := g.Ct.ExecCommand(cmd)
			if err != nil {
				fmt.Fprintf(w, err.Error()+"\n")
				w.Flush()
				continue
			}
			fmt.Fprintf(w, str+"\n")
		} else {
			// other, command does not exist
			fmt.Fprintf(w, "Unknown command: %v %q\n", cmd.Cmd, cmd.Args)
		}
		w.Flush()
	}
}

//prompt user for a line
func prompt(c *context.Context, r *bufio.Reader, w *bufio.Writer) (string, error) {
	fmt.Fprintf(w, "\n" + c.Name+" > ")
	err := w.Flush()
	if err != nil {
		return "", err
	}
	line, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.Trim(line, " \n"), nil
}
