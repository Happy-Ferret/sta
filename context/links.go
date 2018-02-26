package context

import (
	"errors"
)

// Link is a type for joining two contexts
type Link struct {
	Name   string
	locked bool
	Key    string
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
	for _, l := range c.Links {
		if l.Name == name {
			return &l, nil
		}
	}
	return nil, errors.New("no such link")
}

// IsLocked get whether a link is locked
func (l *Link) IsLocked() bool {
	return l.locked
}

// GetTarget returns the target of a link
func (l *Link) GetTarget() *Context {
	return l.target
}

// Lock locks a link if the key is correct
func (l *Link) Lock(key string) {
	if key == l.Key {
		l.locked = true
	}
}

// Unlock unlocks a link if the key is correct
func (l *Link) Unlock(key string) {
	if key == l.Key {
		l.locked = false
	}
}
