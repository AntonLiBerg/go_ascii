package main

import (
	"go_ascii/internal/game"
	"go_ascii/internal/scenario"
	"go_ascii/internal/service/interaction"
	"go_ascii/internal/service/movement"
	"go_ascii/internal/service/quit"
	"go_ascii/internal/service/render"
	"go_ascii/internal/world"
	"os"

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
	aMap, entities, components, userInputProfileMap, err := scenario.GetScenarioFromFiles(
		"./scenarios/demo/map.txt",
		"./scenarios/demo/content.txt",
	)
	if err != nil {
		panic(err)
	}
	gameWorld, err := world.NewWorld(aMap, entities, components)
	if err != nil {
		panic(err)
	}
	gameWorld.UserInputProfile = world.NewUserInputProfile(userInputProfileMap)

	services := []game.IService{
		quit.ServiceQuitGame{},
		movement.ServiceMovePlayer{},
		interaction.ServiceInteraction{},
		render.ServiceDrawOnTerminal{},
	}
	keys := make(chan string)
	go func() {
		for {
			var key [1]byte
			os.Stdin.Read(key[:])
			keys <- string(key[:])
		}
	}()
	game.RunGame(gameWorld, services, keys)
}
