package server

// Event
type Event struct {
	Trigger func() bool
	Action func() string
}

