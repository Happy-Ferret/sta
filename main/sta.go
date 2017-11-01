package main

import (
	"sta/server"
	"fmt"
	"strings"
)

func main() {
	server.Launch();
	game := server.NewGame("Jirsad")
	var cmd string
	for {
		fmt.Print("> ")
		fmt.Scanln(&cmd)
		args := strings.Split(cmd, " ")
		fmt.Println()
		if len(args) > 0 && args[0] != "quit" {
			fmt.Println(game.C.T.Call(args[0]))
		} else {
			break
		}
	}
	fmt.Println("Please come again soon!")
}

