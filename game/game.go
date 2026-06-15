package game

import (
	"fmt"
	serv "go_ascii/service"
	wrld "go_ascii/world"
	"sort"
	"time"
)

const (
	s_readyToGetUpdateFunctions = iota
	s_gettingUpdateFunctions
	s_applyingChanges
)

func RunGame(world wrld.World, services []serv.IService, keyInput <-chan string) error {
	if len(services) == 0 {
		return fmt.Errorf("Services is empty")
	}
	state := s_readyToGetUpdateFunctions
	ticker := time.NewTicker(time.Second / 30)
	results := make([]serv.UpdateFuncResult, 0, len(services))
	updateFuncs := make(chan serv.UpdateFuncResult, len(services))

	for {
		if world.StateUser == wrld.S_quit {
			return nil
		}
		switch state {
		case s_readyToGetUpdateFunctions:
			select {
			case key := <-keyInput:
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

			case key := <-keyInput:
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
