package server

type Thing struct {
	name string
	description string
	contentsMax int
	contents []Item
	tags []string
	commands map[string]func(args ...string)
}

func (t Thing) look() string {
	return "You are looking at " + t.name + ": " + t.description
}

func (t Thing) deposit(it Item) {
	t.contents[len(t.contents)] = it
}

func (t Thing) retrieve() Item {
	if len(t.contents) == 0 {
		return Item{}
	}
	defer func() { t.contents = t.contents[1:] }()
	return t.contents[0]
}

func (t Thing) call(command string) {
	if f, ok := t.commands[command]; ok {
		f()
	}
}

