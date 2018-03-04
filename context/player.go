package context

// NewPlayer returns a context representing a player.
func NewPlayer(name string) (player *Context) {
	player = New(name)
	player.Description = "A person called " + name + "."
	player.Properties["player"] = "player"
	return
}
