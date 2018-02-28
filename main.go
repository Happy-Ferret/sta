/*
Package main is the entry point for sta.
The main function lauches an SSH server listening on localhost:2222.
*/
package main

import (
	"github.com/gliderlabs/ssh"
	"github.com/ribacq/sta/context"
	"github.com/ribacq/sta/display"
	"github.com/ribacq/sta/games"
	"log"
	//"regexp"
)

// gameHandler treats an ssh session with one user.
func gameHandler(sess ssh.Session) {
	log.Println(sess.User() + " connected.")
	defer log.Println(sess.User() + " disconnected.")

	// game with player name and first context
	game := games.New(sess.User(), context.Entrance())

	// display
	disp := display.New(sess)
	// main loop
	line := "look"
	for {
		// clear screen
		err := disp.Clear()
		if err != nil {
			log.Print(err.Error())
		}

		// execute command and print output
		out, err := game.Exec(line)
		if err != nil {
			disp.WriteLine(err.Error())
		} else if len(out) > 0 {
			disp.WriteLine(out)
		}

		// exit?
		if game.Quit() {
			break
		}

		// prompt with context name
		disp.CompleteWith(game.AllCommands())
		line, err = disp.ReadLine(game.Context.Name)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
}

// main function launches the server.
func main() {
	log.Println("server started")
	log.Fatal(ssh.ListenAndServe(":2222", gameHandler, ssh.HostKeyFile("./id_rsa")))
}
