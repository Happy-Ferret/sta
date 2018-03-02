package context

var (
	entrance = New("Dungeon Entrance")
	corridor = New("Corridor")
	temple   = New("Temple")
	key      = New("a small key")
	apple    = New("an apple")
)

// Entrance returns a context for the first room in the dungeon.
func Entrance() *Context {
	return entrance
}

func init() {
	entrance.Description = "You are at the entrance of a most fearsome dungeon. Be brave and feel free to /ask for help/help/. The only way is a corridor going to the **north**."

	corridor.Description = "This is a mossy, dirty corridor carven in the mountainside. You can go **south** back to the entrance, and there is a solid iron door barring the way up **north**."

	temple.Description = "You are inside of a temple of blue-green translucent stone. A chandelier hangs down over an altar. The exit is an iron door to the **south**."

	key.Description = "It is a small metal *key*. It doesn’t look like anything special."
	key.MakeTakeable(true)
	key.Properties["key"] = "iron#0"

	apple.Description = "It is a yellow *apple*. It smells good."

	entrance.AddLink("north", "", false, corridor)
	corridor.AddLink("south", "", false, entrance)
	corridor.Contents = append(corridor.Contents, key, apple)

	corToTemple := corridor.AddLink("north", "iron#0", true, temple)
	templeToCor := temple.AddLink("north", "iron#0", true, corridor)
	corToTemple.slaves = append(corToTemple.slaves, templeToCor)
	templeToCor.slaves = append(templeToCor.slaves, corToTemple)
}
