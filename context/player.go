package context

func NewPlayer(name string) *Context {
	return &Context{
		Name:           name,
		Description:    "A player called " + name,
		commandActions: map[string]commandFunc{},
		Contents:       []*Context{},
	}
}
