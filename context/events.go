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
	TakeDropEvent
)

type takeDropEventContent struct {
	ctx  *Context
	drop bool
}

// handleEvents is meant to be lauched in a goroutine.
func (c *Context) handleEvents() {
	for {
		select {
		case e := <-c.EventsCH:
			if _, ok := c.Properties["player"]; !ok {
				for _, ctx := range c.Contents {
					if _, ok := ctx.Properties["player"]; ok {
						ctx.EventsCH <- e
					}
				}
			} else {
				switch e.Type {
				case CharacterDoesEvent:
					if verb, ok := e.Content.(string); ok && e.Source != c {
						c.OutCH <- "*" + e.Source.Name + "* " + verb + "."
					}
				case LookEvent:
					if ctx, ok := e.Content.(*Context); ok && e.Source != c {
						c.OutCH <- "*" + e.Source.Name + "* is looking at " + ctx.Name + "."
					}
				case TakeDropEvent:
					if cs, ok := e.Content.(takeDropEventContent); ok {
						if c != e.Source {
							if !cs.drop {
								c.OutCH <- "*" + e.Source.Name + "* takes *" + cs.ctx.Name + "*."
							} else {
								c.OutCH <- "*" + e.Source.Name + "* drops *" + cs.ctx.Name + "*."
							}
						} else {
							if !cs.drop {
								c.OutCH <- "You take *" + cs.ctx.Name + "*."
							} else {
								c.OutCH <- "You drop *" + cs.ctx.Name + "*."
							}
						}
					}
				}
			}
		}
	}
}
