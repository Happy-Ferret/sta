package server

import "io"

type Game struct {
	c Character
	in io.Reader
	out io.Writer
}

