package context

import (
	"errors"
	"regexp"
)

// Link is a type for joining two contexts
type Link struct {
	Name   string
	locked bool
	key    string
	target *Context
}

// AddLink adds a new link to a context
func (c *Context) AddLink(target *Context, name string, key string) {
	c.Links = append(c.Links, Link{name, key != "", key, target})
}

// AddDoubleLink adds links a->b and b->a
func AddDoubleLink(ctx1 *Context, ctx2 *Context, name string, key string) {
	ctx1.AddLink(ctx2, name, key)
	ctx2.AddLink(ctx1, name, key)
}

// GetLink returns a link from the context by name, or an error if no such link exists
func (c *Context) GetLink(name string) (*Link, error) {
	if name == "" {
		return nil, errors.New("No name was given.")
	}

	for _, l := range c.Links {
		if matched, err := regexp.Match("^"+name+".*$", []byte(l.Name)); err == nil && matched {
			return &l, nil
		}
	}
	return nil, errors.New("No such link.")
}

// IsLocked returns whether a link is locked
func (l *Link) IsLocked() bool {
	return l.locked
}

// HasKey returns whether the link requires a key
func (l *Link) HasKey() bool {
	return l.key != ""
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
