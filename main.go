/*
Package main is the entry point for sta.
The main function lauches an SSH server listening on localhost:2222.
*/
package main

import (
	"github.com/gliderlabs/ssh"
	"github.com/ribacq/sta/context"
	//_ "github.com/ribacq/sta/db"
	"github.com/ribacq/sta/games"
	"golang.org/x/crypto/ssh/terminal"
	"log"
)

func gameHandler(sess ssh.Session) {
	log.Println(sess.User() + " connected.")
	// Game variable with name, context
	term := terminal.NewTerminal(sess, "> ")
	game := games.New(sess.User(), context.Entrance())
	term.Write([]byte("Hello " + game.Player.Name + "!\n"))
	out, err := game.Exec("look")
	if err != nil {
		return
	}
	term.Write([]byte("\n" + out + "\n\n"))
	for !game.Quit() {
		term.SetPrompt(game.Context.Name + " > ")
		line, err := term.ReadLine()
		if err != nil {
			log.Println(err.Error())
			return
		}
		out, err := game.Exec(line)
		if err != nil {
			term.Write([]byte("\n" + err.Error() + "\n\n"))
		} else if len(out) > 0 {
			term.Write([]byte("\n" + out + "\n\n"))
		}
	}
	log.Println(game.Player.Name + " disconnected.")
}

func main() {
	log.Println("server started")
	log.Fatal(ssh.ListenAndServe(":2222", gameHandler, ssh.HostKeyFile("./id_rsa")))
}
