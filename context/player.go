package context

func NewPlayer(name string) *Context {
	return &Context{
		Name:           name,
		Description:    "A player called " + name,
		CommandActions: map[string]commandFunc{},
		Contents:       []*Context{},
	}
}
