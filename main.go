package main

import (
	aiapi "go_ascii/aiAPI"
	serv "go_ascii/service"
	wrld "go_ascii/world"
	"os"
	"sort"
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

	const (
		s_readyToGetUpdateFunctions = iota
		s_gettingUpdateFunctions
		s_applyingChanges
	)
	ticker := time.NewTicker(time.Second / 30)
	services := []serv.IService{}
	state := s_readyToGetUpdateFunctions
	updateFuncs := make(chan serv.UpdateFuncResult, len(services))
	results := make([]serv.UpdateFuncResult, 0, len(services))

	for {
		switch state {
		case s_readyToGetUpdateFunctions:
			select {
			case key := <-keys:
				world.ClearUserInput()
				world.SetKeyDown(key)

			case <-ticker.C:
				snapshot := world.Clone()
				results = results[:0]

				for _, service := range services {
					go func(service serv.IService) {
						updateFuncs <- service.GetUpdateFunc(snapshot)
					}(service)
				}

				state = s_gettingUpdateFunctions
			}

		case s_gettingUpdateFunctions:
			select {
			case result := <-updateFuncs:
				results = append(results, result)

				if len(results) == len(services) {
					state = s_applyingChanges
				}

			case key := <-keys:
				world.ClearUserInput()
				world.SetKeyDown(key)
			}

		case s_applyingChanges:
			sort.SliceStable(results, func(i, j int) bool {
				return results[i].Order < results[j].Order
			})

			for _, result := range results {
				if result.Err != nil {
					panic(result.Err)
				}
				if result.UpdateFunc != nil {
					result.UpdateFunc(&world)
				}
			}

			state = s_readyToGetUpdateFunctions
		}
	}
}
