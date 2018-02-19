package context

import (
	"errors"
)

type Link struct {
	Name string
	locked bool
	key string
	target *Context
}

func (c *Context) AddLink(target *Context, name string) {
	c.Links = append(c.Links, Link{name, false, "", target})
}

func (c *Context) GetLink(name string) (*Link, error) {
	for _, l := range c.Links {
		if l.Name == name {
			return &l, nil
		}
	}
	return nil, errors.New("No such link.")
}

func (l *Link) IsLocked() bool { return l.locked }
func (l *Link) GetTarget() *Context { return l.target }
func (l *Link) SetKey(key string) { l.key = key }

func (l *Link) Lock(key string) {
	if key == l.key {
		l.locked = true
	}
}

func (l *Link) Unlock(key string) {
	if key == l.key {
		l.locked = false
	}
}

