package main

import (
	"fmt"
	"os"
	"path/filepath"

	"go_ascii/internal/game"
	"go_ascii/internal/scenario"
	"go_ascii/internal/service/control"
	"go_ascii/internal/service/infobox"
	"go_ascii/internal/service/interaction"
	"go_ascii/internal/service/movement"
	"go_ascii/internal/service/quit"
	"go_ascii/internal/service/render"
	"go_ascii/internal/world"

	"golang.org/x/term"
)

const scenariosDirectory = "scenarios"

func main() {
	name, err := scenario.ReadCurrentScenario(filepath.Join(scenariosDirectory, "current_scenario.txt"))
	if err != nil {
		panic(err)
	}

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	if err := runScenario(name); err != nil {
		panic(err)
	}
}

func runScenario(name string) error {
	directory := filepath.Join(scenariosDirectory, name)
	serviceNames, err := scenario.ReadStartupServices(filepath.Join(directory, "startup_config.txt"))
	if err != nil {
		return err
	}
	services, err := servicesFromNames(serviceNames)
	if err != nil {
		return fmt.Errorf("load scenario %q services: %w", name, err)
	}

	aMap, entities, components, inputProfiles, uiLayout, uis, err := scenario.GetScenarioFromFiles(
		filepath.Join(directory, "map.txt"),
		filepath.Join(directory, "content.txt"),
		filepath.Join(directory, "ui.txt"),
	)
	if err != nil {
		return fmt.Errorf("load scenario %q: %w", name, err)
	}
	gameWorld, err := world.NewWorldWithUI(aMap, entities, components, inputProfiles, uiLayout, uis)
	if err != nil {
		return fmt.Errorf("build scenario %q world: %w", name, err)
	}

	keys := make(chan string)
	go func() {
		for {
			var key [1]byte
			os.Stdin.Read(key[:])
			keys <- string(key[:])
		}
	}()
	return game.RunGame(gameWorld, services, keys)
}

func servicesFromNames(names []string) ([]game.IService, error) {
	services := make([]game.IService, 0, len(names))
	for _, name := range names {
		service, err := serviceFromName(name)
		if err != nil {
			return nil, err
		}
		services = append(services, service)
	}
	return services, nil
}

func serviceFromName(name string) (game.IService, error) {
	switch name {
	case "ServiceQuitGame":
		return quit.ServiceQuitGame{}, nil
	case "ServiceMovePlayer":
		return movement.ServiceMovePlayer{}, nil
	case "ServiceInteraction":
		return interaction.ServiceInteraction{}, nil
	case "ServiceControl":
		return control.ServiceControl{}, nil
	case "ServiceDrawOnTerminal":
		return render.ServiceDrawOnTerminal{}, nil
	case "ServiceInfobox":
		return infobox.ServiceInfobox{}, nil
	default:
		return nil, fmt.Errorf("service %q is not registered", name)
	}
}
