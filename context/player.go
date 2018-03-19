package context

// NewPlayer returns a context representing a player.
func NewPlayer(name string, ctx *Context) (player *Context) {
	player = New(name)
	player.SetContainer(ctx)
	player.SetDescription("A person called " + name + ".")
	player.Properties["player"] = "player"
	return
}
