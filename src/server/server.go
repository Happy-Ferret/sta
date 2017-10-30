package server

import "fmt"

func Launch() {
	fmt.Println("Hello world!")
	jirsad := Character{ thing: Thing{ name: "Jirsad", description: "An awesome hero" } }
	fmt.Println(jirsad.sayHi())
}

