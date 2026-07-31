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

	runScenario("skyship")
}

func runScenario(name string) {
	aMap, entities, components, inputProfiles, err := scenario.GetScenarioFromFiles(
		"./scenarios/"+name+"/map.txt",
		"./scenarios/"+name+"/content.txt",
	)
	if err != nil {
		panic(err)
	}
	gameWorld, err := world.NewWorld(aMap, entities, components, inputProfiles)
	if err != nil {
		panic(err)
	}
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
