package main

import (
	aiapi "go_ascii/aiAPI"
	wrld "go_ascii/world"
	"os"

	"golang.org/x/term"
)

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	world := wrld.NewWorld()
	ai := aiapi.New()
	aMap, entities, components := ai.GetAsciiMapAndEntitiesFromFile("./scenarios/demo/map.txt")

	for pos, ch := range aMap {
		eName := entities[ch]
		eComps := components[eName]
		ok := world.AddEntity(pos, eComps)
		if ok != nil {
			panic(ok.Error())
		}
	}
	for {
		aiapi.UpdateTerminal(world)
	}
}
