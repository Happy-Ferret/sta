/*
Package main is the entry point for sta.
The main function lauches an SSH server listening on localhost:2222.
*/
package main

import (
	"encoding/json"
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
	disp.Motd()

	// game with player name and first context
	g := game.New(sess.User(), context.Entrance())
	g.Player.Container.EventsCH <- context.Event{g.Player, context.ConnectionEvent, nil}
	defer func() { g.Player.Container.EventsCH <- context.Event{g.Player, context.DisconnectionEvent, nil} }()

	// main loops
	var err error
	var line string
	if err := g.Exec("look"); err != nil {
		log.Println(err.Error())
		return
	}
	disp.AppendComplete(g.AllCommands())

	// output loop
	go func() {
		for {
			select {
			case out := <-g.Player.Container.OutCH:
				// if the context receives en event, forward it to all contained players
				for _, ctx := range g.Player.Container.Contents {
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
						g.Quit <- true
						return
					}
					disp.AppendComplete(cmds)
				}
			}
		}
	}()

	// input loop: read line, exec and reset autocomplete
	go func() {
		oldctx := g.Player.Container
		for {
			line, err = disp.ReadLine(g.Player.Name + " | " + g.Player.Container.Name)
			if err != nil {
				log.Println(err.Error())
				g.Quit <- true
				return
			}
			if err = g.Exec(line); err != nil {
				log.Println(err.Error())
				g.Quit <- true
				return
			}
			if oldctx != g.Player.Container {
				oldctx = g.Player.Container
				disp.ResetComplete()
				disp.AppendComplete(g.AllCommands())
				j, err := json.Marshal(oldctx)
				if err != nil {
					log.Println(err.Error())
					return
				}
				log.Printf("len(json(ctx)): %v", len(j))
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
