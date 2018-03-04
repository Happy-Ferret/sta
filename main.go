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

	// main loops
	var err error
	var line string
	if err := g.Exec("look"); err != nil {
		log.Println(err.Error())
		return
	}
	oldctx := g.Context
	disp.AppendComplete(g.AllCommands())

	// output loop
	go func() {
		for {
			select {
			case out := <-g.Context.OutCH:
				// if the context receives an event, forward it to all contained players
				for _, ctx := range g.Context.Contents {
					if _, ok := ctx.Properties["player"]; ok {
						ctx.OutCH <- out
					}
				}
			case out := <-g.Player.OutCH:
				// if the player receives an event, write to output
				if len(out) > 0 {
					cmds, err := disp.WriteParsed(out)
					if err != nil {
						log.Println(err.Error())
						return
					}
					disp.AppendComplete(cmds)
				}
			}
		}
	}()

	// input loop: read line, exec and reset autocomplete
	go func() {
		for {
			line, err = disp.ReadLine(g.Player.Name + " | " + g.Context.Name)
			if err != nil {
				log.Println(err.Error())
				return
			}
			if err = g.Exec(line); err != nil {
				log.Println(err.Error())
				return
			}
			if oldctx != g.Context {
				oldctx = g.Context
				disp.ResetComplete()
				disp.AppendComplete(g.AllCommands())
			}
		}
	}()

	<-g.Quit
}

// main function launches the server.
func main() {
	log.Println("server started")
	log.Fatal(ssh.ListenAndServe(":2222", gameHandler, ssh.HostKeyFile("./id_rsa")))
}
