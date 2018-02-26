package main

import (
	"fmt"
	"github.com/gliderlabs/ssh"
	"github.com/ribacq/sta/context"
	"github.com/ribacq/sta/games"
	"github.com/ribacq/sta/parser"
	"golang.org/x/crypto/ssh/terminal"
	"log"
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

	// Game variable with name, context, reader and writer
	term := terminal.NewTerminal(sess, "> ")
	game := games.New(sess.User(), hall, term)
	game.WriteString("Hello " + game.Name + "\n")
	game.WriteString("\n" + game.Context.Look() + "\n")
mainLoop:
	for {
		line, err := game.Prompt()
		if err != nil {
			return
		}
		cmd := parser.Parse(line)

		if cmd.Cmd == "" {
			// command was nothing, skip
			continue
		} else if parser.IsBuiltin(cmd) {
			// command is a builtin, call parser
			out, err := parser.ExecBuiltin(cmd, game.Context)
			if err != nil {
				game.WriteString(err.Error())
				continue
			}
			game.WriteString(out.Message + "\n")
			switch out.Flag {
			case parser.QuitFlag:
				break mainLoop
			}
		} else if _, err := game.Context.GetLink(cmd.Cmd); err == nil {
			game.UseLink(cmd.Cmd)
			game.WriteString(game.Context.Look() + "\n")
		} else if game.Context.HasCommand(cmd) {
			// command is from context
			str, err := game.ExecCommand(cmd)
			if err != nil {
				game.WriteString(err.Error())
				continue
			}
			game.WriteString(str + "\n")
		} else {
			// other, command does not exist
			game.WriteString(fmt.Sprintf("Unknown command: %v %q\n", cmd.Cmd, cmd.Args) + "\n")
		}
	}
}

func main() {
	log.Println("server started")
	log.Fatal(ssh.ListenAndServe(":2222", gameHandler, ssh.HostKeyFile("./id_rsa")))
}
