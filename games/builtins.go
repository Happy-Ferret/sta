package games

func helpCF(g *Game, cmd []string) (out string, err error) {
	return "help...", nil
}

func quitCF(g *Game, cmd []string) (out string, err error) {
	g.quit = true
	return "quit...", nil
}
