package server

import "fmt"

func Launch() {
	fmt.Println("server.Launch() called")
}

func NewGame(characterName string) (g Game) {
	g = Game{}
	g.C.T.Name = characterName
	g.C.T.Description = "a character for a new game"
	g.C.T.Commands = make(map[string]func (args []string) string)
	g.C.T.Commands["look"] = func (args []string) string {
		return g.C.T.Look()
	}
	g.C.T.Commands["help"] = func (args []string) string {
		var ret string
		for cmd := range g.C.T.Commands {
			ret += ", " + cmd
		}
		return ret
	}
	return
}

