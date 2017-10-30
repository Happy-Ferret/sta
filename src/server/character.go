package server

type Character struct {
	thing Thing
	hp int
	hpMax int
	stats struct {
		constitution int
		attention int
		magic int
		influence int
		wisdom int
	}
}

func (c Character) sayHi() string {
	return c.thing.look()
}

