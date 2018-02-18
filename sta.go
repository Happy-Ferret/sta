package main

import (
	"bufio"
	"fmt"
	"github.com/ribacq/sta/context"
	"github.com/ribacq/sta/parser"
	"os"
	"strings"
)

func main() {
	// setting I/O
	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)

	// creating context
	c := context.Default()

	fmt.Fprintf(w, "Hello world!\n")
	w.Flush()
mainLoop:
	for {
		line, err := prompt(c, r, w)
		if err != nil {
			return
		}
		cmd := parser.Parse(line)

		if cmd.Cmd == "" {
			// command was nothing, skip
			continue
		} else if parser.IsBuiltin(cmd) {
			// command is a builtin, call parser
			out, err := parser.ExecBuiltin(cmd, c)
			if err != nil {
				fmt.Fprintf(w, err.Error()+"\n")
			}
			fmt.Fprintf(w, out.Message+"\n")
			switch out.Flag {
			case parser.QuitFlag:
				w.Flush()
				break mainLoop
			}
		} else if c.HasCommand(cmd) {
			// command is from context
			str, err := c.ExecCommand(cmd)
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
	fmt.Fprintf(w, c.Name+" > ")
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
