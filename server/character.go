package server

// Character
type Character struct {
	T Thing
	HP int
	HPMax int
	Stats struct {
		Constitution int
		Attention int
		Magic int
		Influence int
		Wisdom int
	}
}

func (c Character) SayHi() string {
	return c.T.Look()
}

