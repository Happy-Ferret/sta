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
)

// gameHandler treats an ssh session with one user.
func gameHandler(sess ssh.Session) {
	log.Println(sess.User() + " connected.")
	// terminal variables
	rev := string(27) + "[1;7m"
	hl := string(27) + "[3m"
	norm := string(27) + "[0m"
	term := terminal.NewTerminal(sess, "")

	// game with player name and first context
	game := games.New(sess.User(), context.Entrance())

	// welcome message and look on context
	term.Write([]byte("Welcome " + game.Player.Name + "!\n\n"))
	out, err := game.Exec("look")
	if err != nil {
		return
	}
	term.Write([]byte(out + "\n"))

	// main loop
	for !game.Quit() {
		// prompt with context name
		term.SetPrompt("\n" + rev + game.Context.Name + " " + hl + ">" + norm + " ")
		line, err := term.ReadLine()
		if err != nil {
			log.Println(err.Error())
			return
		}

		// execute command and print output
		out, err := game.Exec(line)
		if err != nil {
			term.Write([]byte(err.Error() + "\n"))
		} else if len(out) > 0 {
			term.Write([]byte(out + "\n"))
		}
	}
	log.Println(game.Player.Name + " disconnected.")
}

// main function launches the server.
func main() {
	log.Println("server started")
	log.Fatal(ssh.ListenAndServe(":2222", gameHandler, ssh.HostKeyFile("./id_rsa")))
}
