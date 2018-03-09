package context

// NewPlayer returns a context representing a player.
func NewPlayer(name string, ctx *Context) (player *Context) {
	player = New(name)
	player.Container = ctx
	player.Description = "A person called " + name + "."
	player.Properties["player"] = "player"
	return
}
