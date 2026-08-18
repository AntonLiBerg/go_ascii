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

type updateResult struct {
	serviceIndex int
	update       UpdateFunc
}

func RunGame(gameWorld world.World, services []IService, keyInput <-chan string) error {
	if len(services) == 0 {
		return fmt.Errorf("Services is empty")
	}
	state := s_readyToGetUpdateFunctions
	ticker := time.NewTicker(time.Second / 30)
	defer ticker.Stop()
	results := make([]updateResult, 0, len(services))
	updateFuncs := make(chan updateResult, len(services))

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

				for index, service := range services {
					go func(index int, service IService) {
						updateFuncs <- updateResult{serviceIndex: index, update: service.GetUpdateFunc(snapshot)}
					}(index, service)
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
				if results[i].update.Order != results[j].update.Order {
					return results[i].update.Order < results[j].update.Order
				}
				return results[i].serviceIndex < results[j].serviceIndex
			})

			for _, result := range results {
				if result.update.UpdateFunc != nil {
					nextWorld, err := result.update.UpdateFunc(gameWorld)
					if err != nil {
						return err
					}
					gameWorld = nextWorld
				}
			}
			state = s_readyToGetUpdateFunctions
			gameWorld.InfoboxContent = ""
			gameWorld.UIContent["infobox"] = gameWorld.UIEmpty["infobox"]
			gameWorld.IterationNr++
		}
	}
}
