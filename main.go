/*
Package main is the entry point for sta.
The main function lauches an SSH server listening on localhost:2222.
*/
package main

import (
	"github.com/gliderlabs/ssh"
	"github.com/ribacq/sta/context"
	"github.com/ribacq/sta/display"
	"github.com/ribacq/sta/game"
	"log"
)

// gameHandler treats an ssh session with one user.
func gameHandler(sess ssh.Session) {
	log.Println(sess.User() + " connected.")
	defer log.Println(sess.User() + " disconnected.")

	// display
	disp := display.New(sess)

	// game with player name and first context
	g := game.New(sess.User(), context.Entrance())

	// main loop
	line := "look"
	oldctx := g.Context
	disp.AppendComplete(g.AllCommands())
	for {
		// execute command and print output
		out, err := g.Exec(line)
		if oldctx != g.Context {
			oldctx = g.Context
			disp.ResetComplete()
			disp.AppendComplete(g.AllCommands())
		}
		if err != nil {
			if err = disp.WriteLine(err.Error()); err != nil {
				log.Println(err.Error())
				return
			}
		} else if len(out) > 0 {
			cmds, err := disp.WriteParsed(out)
			if err != nil {
				log.Println(err.Error())
				return
			}
			disp.AppendComplete(cmds)
		}

		// quit?
		if g.Quit() {
			break
		}

		// prompt with context name
		line, err = disp.ReadLine(g.Context.Name)
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
