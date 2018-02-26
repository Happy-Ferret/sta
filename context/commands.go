package context

// Look command gives a description of the context and available un-hidden links.
func Look(c *Context, cmd []string) (out string, err error) {
	lookStr := c.Description + "\nLinks:"
	for _, l := range c.Links {
		lookStr += " " + l.Name
	}
	return lookStr, nil
}

// Take puts an item into player’s bag.
func Take(c *Context, cmd []string) (out string, err error) {
	return "Taken!", nil
}
