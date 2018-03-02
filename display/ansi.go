package display

var (
	// display styles
	styles = map[string]string{
		//styles
		"reset":         "\033[0m",
		"bright":        "\033[1m",
		"dim":           "\033[2m",
		"italic":        "\033[3m",
		"underscore":    "\033[4m",
		"blink":         "\033[5m",
		"reverse":       "\033[7m",
		"hidden":        "\033[8m",
		"strikethrough": "\033[9m",

		// fg colors
		"black":   "\033[30m",
		"red":     "\033[31m",
		"green":   "\033[32m",
		"yellow":  "\033[33m",
		"blue":    "\033[34m",
		"magenta": "\033[35m",
		"cyan":    "\033[36m",
		"white":   "\033[37m",

		// bg colors
		"blackBG":   "\033[40m",
		"redBG":     "\033[41m",
		"greenBG":   "\033[42m",
		"yellowBG":  "\033[43m",
		"blueBG":    "\033[44m",
		"magentaBG": "\033[45m",
		"cyanBG":    "\033[46m",
		"whiteBG":   "\033[47m",
	}

	// some other escape sequences
	escape = map[string]string{
		"clear":     "\033[2J",
		"up":        "\033[A",
		"down":      "\033[B",
		"right":     "\033[C",
		"left":      "\033[D",
		"bol":       "\r",
		"clearline": "\033[2K",
	}
)
