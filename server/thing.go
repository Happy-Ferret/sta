package server

import "errors"

// Thing
type Thing struct {
	Name string
	Description string
	contentsMax int
	Contents []Item
	Tags []string
	Commands map[string]func(args []string) string
}

func (t Thing) Look() string {
	return "You are looking at " + t.Name + ": " + t.Description
}

func (t Thing) Deposit(it Item) {
	t.Contents[len(t.Contents)] = it
}

func (t Thing) Retrieve() (Item, error) {
	if len(t.Contents) == 0 {
		return Item{}, errors.New("Thing.Retrieve(): no item to retrieve")
	}
	defer func() { t.Contents = t.Contents[1:] }()
	return t.Contents[0], nil
}

func (t Thing) Call(args ...string) string {
	if f, ok := t.Commands[args[0]]; ok {
		return f(args)
	}
	return "such command does not exist"
}

