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
	var completionOptions []string
	var err error
	var out string
	for {
		// execute command and print output
		out, err = game.Exec(line)
		if err != nil {
			if err = disp.WriteLine(err.Error()); err != nil {
				log.Println(err.Error())
				return
			}
		} else if len(out) > 0 {
			completionOptions, err = disp.WriteParsed(out)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}

		// quit?
		if game.Quit() {
			break
		}

		// prompt with context name
		disp.CompleteWith(append(game.AllCommands(), completionOptions...))
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
