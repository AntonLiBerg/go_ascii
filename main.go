package main

import (
	aiapi "go_ascii/aiAPI"
	serv "go_ascii/service"
	wrld "go_ascii/world"
	"os"
	"time"

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

	keys := make(chan string)
	go func() {
		for {
			var key [1]byte
			os.Stdin.Read(key[:])
			keys <- string(key[:])
		}
	}()
	ticker := time.NewTicker(time.Second / 30)

	services := []serv.IService{}
	changes := []func(*wrld.World){}
	for {
		select {
		case key := <-keys:
			world.ClearUserInput()
			world.SetKeyDown(key)
		case <-ticker.C:

			for _, service := range services {
				_ = service.Update(&world)
			}
		}
	}
}
