package main

import (
	"bufio"
	"os"
	"strings"
	"fmt"
)

func main() {
	// setting I/O
	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)

	fmt.Fprintf(w, "Bienvenue au Royaume de Reosal!\n")
	w.Flush()
	for {
		fmt.Fprintf(w, "> ")
		w.Flush()
		line, _ := r.ReadString('\n')
		line = strings.Trim(line, " \n")
		args := strings.Split(line, " ")

		if strings.Compare(args[0], "quit") == 0 {
			break
		} else {
			fmt.Fprintf(w, "%v %v\n", args[0], args[1:])
		}
		w.Flush()
	}

}

