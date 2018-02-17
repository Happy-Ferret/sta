package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Command struct {
	Cmd  string
	Args []string
}

func prompt(promptStr string, r *bufio.Reader, w *bufio.Writer) (cmd *Command, err error) {
	fmt.Fprintf(w, promptStr)
	err = w.Flush()
	if err != nil {
		return nil, err
	}
	line, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}
	line = strings.Trim(line, " \n")
	sections := strings.Split(line, " ")
	cmd = &Command{sections[0], sections[1:]}
	return cmd, nil
}

func main() {
	// setting I/O
	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)

	fmt.Fprintf(w, "Hello world!\n")
	w.Flush()
	for {
		cmd, err := prompt("> ", r, w)
		if err != nil {
			break
		}

		if strings.Compare(cmd.Cmd, "quit") == 0 {
			break
		} else {
			fmt.Fprintf(w, "%v %v\n", cmd.Cmd, cmd.Args)
		}
		w.Flush()
	}
}
