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
	temple   = context.New("Temple")
	key      = context.New("a small key")
	apple    = context.New("an apple")
)

// Entrance returns a context for the first room in the dungeon.
func Entrance() *context.Context {
	return entrance
}

func init() {
	entrance.Description = "You are at the entrance of a most fearsome dungeon. Be brave. The only way is a corridor going to the north."

	corridor.Description = "This is a mossy, dirty corridor carven in the mountainside. You can go south back to the entrance, and there is a solid iron door barring the way up north."

	temple.Description = "You are inside of a temple of blue-green translucent stone. A chandelier hangs down over an altar. The exit is an iron door to the south."

	key.Description = "It is a small metal key. It doesn’t look like anything special."
	key.MakeTakeable(true)
	key.Properties["key"] = "iron#0"

	apple.Description = "It is a yellow apple. It smells good."

	entrance.AddLink(corridor, "north", "")
	corridor.AddLink(entrance, "south", "")
	corridor.Contents = append(corridor.Contents, key, apple)
	corridor.AddLink(temple, "north", "iron#0")
	temple.AddLink(corridor, "south", "iron#0")
}
