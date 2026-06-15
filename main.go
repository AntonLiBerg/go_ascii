package main

import (
	ai "go_ascii/aiAPI"
	gme "go_ascii/game"
	serv "go_ascii/service"
	"os"

	wrld "go_ascii/world"

	"golang.org/x/term"
)

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	runDemo()
}

func runDemo() {
	aMap, entities, components, userInputProfileMap, err := ai.GetAsciiMapAndEntitiesFromFile("./scenarios/demo/map.txt")
	if err != nil {
		panic(err)
	}
	world, err := wrld.NewWorld(aMap, entities, components)
	if err != nil {
		panic(err)
	}
	world.UserInputProfile.KeyQuitGame = userInputProfileMap["quitgame"]

	services := []serv.IService{
		&serv.ServiceDrawOnTerminal{},
		&serv.ServiceQuitGame{},
	}
	keys := make(chan string)
	go func() {
		for {
			var key [1]byte
			os.Stdin.Read(key[:])
			keys <- string(key[:])
		}
	}()
	gme.RunGame(world, services, keys)
}
