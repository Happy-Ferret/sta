package server

type Event struct {
	trigger func() bool
	action func() string
}

