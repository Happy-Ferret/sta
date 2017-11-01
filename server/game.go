package server

import "io"

//Game
type Game struct {
	C Character
	In io.Reader
	Out io.Writer
}

