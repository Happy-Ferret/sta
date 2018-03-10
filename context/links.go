package context

import (
	"errors"
	"regexp"
)

// Link is a type for joining two contexts.
type Link struct {
	name   string
	key    string
	locked bool
	target *Context
	slaves []*Link
}

func (l Link) MarshalJSON() ([]byte, error) {
	return []byte("{\"name\": \"" + l.name + "\"}"), nil
}

// AddLink adds a new link to a context and returns it.
func (c *Context) AddLink(name, key string, locked bool, target *Context) *Link {
	l := &Link{
		name:   name,
		key:    key,
		locked: locked,
		target: target,
	}
	c.Links = append(c.Links, l)
	return l
}

// Name gets link name.
func (l *Link) Name() string {
	return l.name
}

// Locked returns whether a link is locked.
func (l *Link) Locked() bool {
	return l.locked
}

// HasKey returns whether the link requires a key.
func (l *Link) HasKey() bool {
	return l.key != ""
}

// Target returns target context after checking the key.
func (l *Link) Target() *Context {
	return l.target
}

// GetLink returns a link from the context by name, or an error if no such link exists
func (c *Context) GetLink(name string) (*Link, error) {
	if name == "" {
		return nil, errors.New("No name was given.")
	}

	for _, l := range c.Links {
		if matched, err := regexp.Match("^"+regexp.QuoteMeta(name)+".*$", []byte(l.name)); err == nil && matched {
			return l, nil
		}
	}
	return nil, errors.New("No such link.")
}

// Try to pass through a link using a key in player Contents.
// Returns the link target and a boolean indicating if the link could be unlocked.
func (l *Link) Try(player *Context) (target *Context, ok bool) {
	// link is not locked
	if !l.locked {
		return l.target, true
	}

	// link is locked: look for a key in player.Contents
	for _, ctx := range player.Contents {
		if val, ok := ctx.Properties["key"]; ok && val == l.key {
			return l.target, true
		}
	}
	// the required key was not found
	return nil, false
}
