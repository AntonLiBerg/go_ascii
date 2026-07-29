package game

import (
	"fmt"
	"go_ascii/internal/world"
	"sort"
	"time"
)

const (
	s_readyToGetUpdateFunctions = iota
	s_gettingUpdateFunctions
	s_applyingChanges
)

func RunGame(gameWorld world.World, services []IService, keyInput <-chan string) error {
	if len(services) == 0 {
		return fmt.Errorf("Services is empty")
	}
	state := s_readyToGetUpdateFunctions
	ticker := time.NewTicker(time.Second / 30)
	defer ticker.Stop()
	results := make([]UpdateFunc, 0, len(services))
	updateFuncs := make(chan UpdateFunc, len(services))

	for {
		if gameWorld.ShouldQuit {
			return nil
		}
		switch state {
		case s_readyToGetUpdateFunctions:
			select {
			case key := <-keyInput:
				gameWorld.KeyDown = key

			case <-ticker.C:
				snapshot := gameWorld.Clone()
				results = results[:0]

				for _, service := range services {
					go func(service IService) {
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

			case key := <-keyInput:
				gameWorld.KeyDown = key
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
					result.UpdateFunc(&gameWorld)
				}
			}
			state = s_readyToGetUpdateFunctions
			gameWorld.IterationNr++
		}
	}
}
