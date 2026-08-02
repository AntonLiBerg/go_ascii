package tests

import (
	component "go_ascii/internal"
	"go_ascii/internal/service/interaction"
	"go_ascii/internal/world"
	"testing"
)

func TestServiceOpensSingleNeighborDoorAndClearsInteractKey(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	gameWorld.UserInputProfile = world.UserInputProfile{KeyInteract: "e"}
	addTestEntity(t, &gameWorld, [2]int{1, 1}, map[string][]string{
		"pos": {}, "ascii": {"o"}, "player": {},
	})
	doorID := addTestEntity(t, &gameWorld, [2]int{2, 1}, map[string][]string{
		"pos": {}, "ascii": {"D"}, "impassable": {}, "interactable": {component.InteractionTypeDoor},
	})
	gameWorld.KeyDown = "e"

	result := interaction.ServiceInteraction{}.GetUpdateFunc(gameWorld)
	if result.UpdateFunc == nil {
		t.Fatal("expected interact update func")
	}
	result.UpdateFunc(&gameWorld)

	if _, isClosed := gameWorld.Impassable[doorID]; isClosed {
		t.Fatal("expected door to be open after interaction")
	}
	if gameWorld.KeyDown != "" {
		t.Fatalf("expected interact key to be cleared, got %q", gameWorld.KeyDown)
	}
	if !gameWorld.HasChanged {
		t.Fatal("expected world to be marked changed after opening door")
	}
}

func TestServiceOpensAndExitsTerminalRoom(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	topdown := world.UserInputProfile{KeyInteract: "e"}
	terminal := world.UserInputProfile{KeyExit: "e"}
	gameWorld.InputProfiles["topdown"] = topdown
	gameWorld.InputProfiles["terminal"] = terminal
	gameWorld.InputProfileByRoom[""] = "topdown"
	gameWorld.InputProfileByRoom["terminal"] = "terminal"
	gameWorld.UserInputProfile = topdown
	playerID := addTestEntity(t, &gameWorld, [2]int{1, 1}, map[string][]string{
		"pos": {}, "ascii": {"o"}, "player": {},
	})
	commandTableID := addTestEntity(t, &gameWorld, [2]int{2, 1}, map[string][]string{
		"pos": {}, "ascii": {"T"}, "impassable": {}, "interactable": {component.InteractionTypeCommandTable},
	})
	gameWorld.Terminals[component.Position{X: 2, Y: 1}] = "terminal"
	playerPosition := gameWorld.Pos[playerID]
	gameWorld.KeyDown = "e"

	open := interaction.ServiceInteraction{}.GetUpdateFunc(gameWorld)
	open.UpdateFunc(&gameWorld)

	if gameWorld.ViewRoom != "terminal" {
		t.Fatalf("expected terminal view, got %q", gameWorld.ViewRoom)
	}
	if gameWorld.UserInputProfile != terminal {
		t.Fatalf("expected terminal input profile, got %+v", gameWorld.UserInputProfile)
	}
	if gameWorld.Pos[playerID] != playerPosition || gameWorld.Pos[commandTableID] != (component.Position{X: 2, Y: 1}) {
		t.Fatal("expected opening terminal not to move physical entities")
	}

	gameWorld.HasChanged = false
	gameWorld.KeyDown = "e"
	closeTerminal := interaction.ServiceInteraction{}.GetUpdateFunc(gameWorld)
	closeTerminal.UpdateFunc(&gameWorld)

	if gameWorld.ViewRoom != "" {
		t.Fatalf("expected terminal view to close, got %q", gameWorld.ViewRoom)
	}
	if gameWorld.UserInputProfile != topdown {
		t.Fatalf("expected topdown input profile to be restored, got %+v", gameWorld.UserInputProfile)
	}
	if gameWorld.KeyDown != "" || !gameWorld.HasChanged {
		t.Fatal("expected exit key to be consumed and world marked changed")
	}
}

func TestServiceClosesOpenDoor(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	gameWorld.UserInputProfile = world.UserInputProfile{KeyInteract: "e"}
	addTestEntity(t, &gameWorld, [2]int{1, 1}, map[string][]string{"pos": {}, "player": {}})
	doorID := addTestEntity(t, &gameWorld, [2]int{2, 1}, map[string][]string{
		"pos": {}, "interactable": {component.InteractionTypeDoor},
	})
	gameWorld.KeyDown = "e"

	result := interaction.ServiceInteraction{}.GetUpdateFunc(gameWorld)
	result.UpdateFunc(&gameWorld)

	if _, isClosed := gameWorld.Impassable[doorID]; !isClosed {
		t.Fatal("expected door to be closed after interaction")
	}
}

func TestServiceIgnoresMultipleNeighborDoors(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	gameWorld.UserInputProfile = world.UserInputProfile{KeyInteract: "e"}
	addTestEntity(t, &gameWorld, [2]int{1, 1}, map[string][]string{"pos": {}, "player": {}})
	eastDoor := addTestEntity(t, &gameWorld, [2]int{2, 1}, map[string][]string{
		"pos": {}, "impassable": {}, "interactable": {component.InteractionTypeDoor},
	})
	westDoor := addTestEntity(t, &gameWorld, [2]int{0, 1}, map[string][]string{
		"pos": {}, "impassable": {}, "interactable": {component.InteractionTypeDoor},
	})
	gameWorld.KeyDown = "e"

	result := interaction.ServiceInteraction{}.GetUpdateFunc(gameWorld)
	result.UpdateFunc(&gameWorld)

	if _, isClosed := gameWorld.Impassable[eastDoor]; !isClosed {
		t.Fatal("expected east door to remain closed")
	}
	if _, isClosed := gameWorld.Impassable[westDoor]; !isClosed {
		t.Fatal("expected west door to remain closed")
	}
	if gameWorld.HasChanged {
		t.Fatal("expected world not to change for multiple targets")
	}
	if gameWorld.KeyDown != "" {
		t.Fatalf("expected interact key to be cleared, got %q", gameWorld.KeyDown)
	}
}

func TestServiceAcceptsNoOpInteractables(t *testing.T) {
	tests := []struct {
		name            string
		interactionType string
	}{
		{name: "helm", interactionType: component.InteractionTypeHelm},
		{name: "command table", interactionType: component.InteractionTypeCommandTable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gameWorld := world.NewWorldEmpty()
			gameWorld.UserInputProfile = world.UserInputProfile{KeyInteract: "e"}
			addTestEntity(t, &gameWorld, [2]int{1, 1}, map[string][]string{"pos": {}, "player": {}})
			targetID := addTestEntity(t, &gameWorld, [2]int{2, 1}, map[string][]string{
				"pos": {}, "impassable": {}, "interactable": {tt.interactionType},
			})
			gameWorld.KeyDown = "e"

			result := interaction.ServiceInteraction{}.GetUpdateFunc(gameWorld)
			if result.UpdateFunc == nil {
				t.Fatal("expected interact update func")
			}
			result.UpdateFunc(&gameWorld)

			if gameWorld.KeyDown != "" {
				t.Fatalf("expected interact key to be cleared, got %q", gameWorld.KeyDown)
			}
			if gameWorld.HasChanged {
				t.Fatal("expected no-op interaction not to change world")
			}
			if _, blocked := gameWorld.Impassable[targetID]; !blocked {
				t.Fatal("expected target to remain impassable")
			}
		})
	}
}

func TestServiceIgnoresUnknownInteractionType(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	gameWorld.UserInputProfile = world.UserInputProfile{KeyInteract: "e"}
	addTestEntity(t, &gameWorld, [2]int{1, 1}, map[string][]string{"pos": {}, "player": {}})
	doorID := addTestEntity(t, &gameWorld, [2]int{2, 1}, map[string][]string{
		"pos": {}, "impassable": {}, "interactable": {"unsupported"},
	})
	gameWorld.KeyDown = "e"

	result := interaction.ServiceInteraction{}.GetUpdateFunc(gameWorld)
	result.UpdateFunc(&gameWorld)

	if _, isClosed := gameWorld.Impassable[doorID]; !isClosed {
		t.Fatal("expected unknown interaction not to toggle impassability")
	}
}

func addTestEntity(t *testing.T, gameWorld *world.World, position [2]int, components map[string][]string) int {
	t.Helper()
	eID := len(gameWorld.Entities)
	if err := gameWorld.AddEntity(position, components); err != nil {
		t.Fatalf("AddEntity returned error: %v", err)
	}
	return eID
}
