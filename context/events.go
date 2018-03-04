package context

// Event is a type for events sent across contexts.
type Event struct {
	Source  *Context
	Type    int
	Content interface{}
}

// The different event types.
const (
	PingEvent = iota
	CharacterDoesEvent
	LookEvent
)

// handleEvents is meant to be lauched in a goroutine.
func (c *Context) handleEvents() {
	for {
		select {
		case e := <-c.EventsCH:
			if e.Source != c {
				if _, ok := c.Properties["player"]; !ok {
					for _, ctx := range c.Contents {
						if _, ok := ctx.Properties["player"]; ok {
							ctx.EventsCH <- e
						}
					}
				} else {
					switch e.Type {
					case CharacterDoesEvent:
						if verb, ok := e.Content.(string); ok {
							c.OutCH <- "*" + e.Source.Name + "* " + verb + "."
						}
					case LookEvent:
						if ctx, ok := e.Content.(*Context); ok {
							c.OutCH <- "*" + e.Source.Name + "* is looking at " + ctx.Name + "."
						}
					}
				}
			}
		}
	}
}
