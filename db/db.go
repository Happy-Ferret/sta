/*
Package db manages database access for github.com/ribacq/sta.
*/
package db

import (
	"github.com/ribacq/sta/context"
)

var (
	entrance = context.New("Dungeon Entrance")
	corridor = context.New("Corridor")
	key      = context.New("a small key")
	apple    = context.New("an apple")
)

// Entrance returns a context for the first room in the dungeon.
func Entrance() *context.Context {
	return entrance
}

func init() {
	entrance.Description = "You are at the entrance of a most fearsome dungeon. Be brave. The only way is a corridor going to the north."
	corridor.Description = "This is a mossy, dirty corridor carven in the mountainside. You can go south back to the entrance, but there is a solid iron gate barring the way up north."
	key.Description = "It is a small metal key. It doesn’t look like anything special."
	apple.Description = "It is a yellow apple. It smells good."

	entrance.AddLink(corridor, "north", "")
	corridor.AddLink(entrance, "south", "")
	corridor.Contents = append(corridor.Contents, key, apple)
}
