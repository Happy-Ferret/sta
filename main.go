/*
Package main is the entry point for sta.
The main function lauches an SSH server listening on localhost:2222.
*/
package main

import (
	"github.com/gliderlabs/ssh"
	"github.com/ribacq/sta/context"
	"github.com/ribacq/sta/games"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"strings"
)

func gameHandler(sess ssh.Session) {
	// default rooms
	// TODO: use a database
	hall := context.New("Palace hall")
	hall.Description = "You are inside the hall of a palace. There is a door to the left."
	kitchen := context.New("Kitchen")
	kitchen.Description = "You are inside a shiny kitchen. There is a camembert on the table and a door to your right."
	camembert := context.New("Camembert")
	camembert.Description = "This is the most beautiful piece of dairy you’ve ever seen…"
	camembert.MakeTakeable()

	context.AddDoubleLink(hall, kitchen, "door", "")
	kitchen.AddLink(camembert, "camembert", "")
	camembert.AddLink(kitchen, "kitchen", "")

	// Game variable with name, context
	term := terminal.NewTerminal(sess, "> ")
	game := games.New(sess.User(), hall)
	term.Write([]byte("Hello "))
	term.Write(term.Escape.Red)
	term.Write([]byte(game.Name))
	term.Write(term.Escape.Reset)
	term.Write([]byte("!\n"))
	for !game.Quit() {
		term.SetPrompt(game.Context.Name + " > ")
		line, err := term.ReadLine()
		if err != nil {
			log.Fatal(err.Error())
			return
		}
		out, err := game.Exec(strings.Split(line, " "))
		if err != nil {
			term.Write([]byte("\n" + err.Error() + "\n\n"))
		} else {
			term.Write([]byte("\n" + out + "\n\n"))
		}
	}
}

func main() {
	log.Println("server started")
	log.Fatal(ssh.ListenAndServe(":2222", gameHandler, ssh.HostKeyFile("./id_rsa")))
}
