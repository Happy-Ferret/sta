package context

func NewPlayer(name string) *Context {
	return &Context{
		Name:        name,
		Description: "A person called " + name + ".",
		Commands:    map[string]CommandFunc{},
		Contents:    []*Context{},
		Properties:  map[string]string{},
	}
}
