package tests

import (
	component "go_ascii/internal"
	"go_ascii/internal/service/interaction"
	"go_ascii/internal/world"
	"testing"
)

func TestServiceInteractionDoors(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*testing.T, *world.World) []int
		wantClosed  []bool
		wantChanged bool
	}{
		{
			name: "opens a single closed door",
			setup: func(t *testing.T, gameWorld *world.World) []int {
				addTestEntity(t, gameWorld, [2]int{1, 1}, map[string][]string{"pos": {}, "player": {}})
				return []int{addTestEntity(t, gameWorld, [2]int{2, 1}, map[string][]string{
					"pos": {}, "impassable": {}, "interactable": {component.InteractionTypeDoor},
				})}
			},
			wantClosed: []bool{false}, wantChanged: true,
		},
		{
			name: "closes a single open door",
			setup: func(t *testing.T, gameWorld *world.World) []int {
				addTestEntity(t, gameWorld, [2]int{1, 1}, map[string][]string{"pos": {}, "player": {}})
				return []int{addTestEntity(t, gameWorld, [2]int{2, 1}, map[string][]string{
					"pos": {}, "interactable": {component.InteractionTypeDoor},
				})}
			},
			wantClosed: []bool{true}, wantChanged: true,
		},
		{
			name: "ignores multiple doors",
			setup: func(t *testing.T, gameWorld *world.World) []int {
				addTestEntity(t, gameWorld, [2]int{1, 1}, map[string][]string{"pos": {}, "player": {}})
				return []int{
					addTestEntity(t, gameWorld, [2]int{2, 1}, map[string][]string{"pos": {}, "impassable": {}, "interactable": {component.InteractionTypeDoor}}),
					addTestEntity(t, gameWorld, [2]int{0, 1}, map[string][]string{"pos": {}, "impassable": {}, "interactable": {component.InteractionTypeDoor}}),
				}
			},
			wantClosed: []bool{true, true},
		},
		{
			name: "ignores unknown interaction",
			setup: func(t *testing.T, gameWorld *world.World) []int {
				addTestEntity(t, gameWorld, [2]int{1, 1}, map[string][]string{"pos": {}, "player": {}})
				return []int{addTestEntity(t, gameWorld, [2]int{2, 1}, map[string][]string{
					"pos": {}, "impassable": {}, "interactable": {"unsupported"},
				})}
			},
			wantClosed: []bool{true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gameWorld := world.NewWorldEmpty()
			gameWorld.UserInputProfile = world.UserInputProfile{KeyInteract: "e"}
			ids := tt.setup(t, &gameWorld)
			gameWorld.KeyDown = "e"
			input := gameWorld
			next := applyUpdate(t, interaction.ServiceInteraction{}.GetUpdateFunc(gameWorld), gameWorld)
			for i, id := range ids {
				_, gotClosed := next.Impassable[id]
				if gotClosed != tt.wantClosed[i] {
					t.Fatalf("door %d closed=%v, want %v", id, gotClosed, tt.wantClosed[i])
				}
			}
			if _, ok := input.Impassable[ids[0]]; tt.name == "opens a single closed door" && !ok {
				t.Fatal("expected update not to mutate input")
			}
			if next.HasChanged != tt.wantChanged {
				t.Fatalf("expected HasChanged=%v, got %v", tt.wantChanged, next.HasChanged)
			}
			if next.KeyDown != "" {
				t.Fatalf("expected interact key cleared, got %q", next.KeyDown)
			}
		})
	}
}

func TestServiceOpensTerminalAndRestoresInputProfile(t *testing.T) {
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
	terminalID := addTestEntity(t, &gameWorld, [2]int{2, 1}, map[string][]string{
		"pos": {}, "ascii": {"T"}, "impassable": {}, "interactable": {component.InteractionTypeTerminal},
	})
	gameWorld.Terminals[component.Position{X: 2, Y: 1}] = "terminal"
	playerPosition := gameWorld.Pos[playerID]
	gameWorld.KeyDown = "e"

	open := interaction.ServiceInteraction{}.GetUpdateFunc(gameWorld)
	gameWorld = applyUpdate(t, open, gameWorld)

	if gameWorld.UserInputProfile != terminal {
		t.Fatalf("expected terminal input profile, got %+v", gameWorld.UserInputProfile)
	}
	if gameWorld.Room != "terminal" {
		t.Fatalf("expected terminal room, got %q", gameWorld.Room)
	}
	if gameWorld.Pos[playerID] != playerPosition || gameWorld.Pos[terminalID] != (component.Position{X: 2, Y: 1}) {
		t.Fatal("expected opening terminal not to move physical entities")
	}

	gameWorld.HasChanged = false
	gameWorld.KeyDown = "e"
	closeTerminal := interaction.ServiceInteraction{}.GetUpdateFunc(gameWorld)
	gameWorld = applyUpdate(t, closeTerminal, gameWorld)

	if gameWorld.UserInputProfile != topdown {
		t.Fatalf("expected topdown input profile to be restored, got %+v", gameWorld.UserInputProfile)
	}
	if gameWorld.Room != playerPosition.Room {
		t.Fatalf("expected room %q after exit, got %q", playerPosition.Room, gameWorld.Room)
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
	gameWorld = applyUpdate(t, result, gameWorld)

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
	gameWorld = applyUpdate(t, result, gameWorld)

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

func TestServiceIgnoresUnknownInteractionType(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	gameWorld.UserInputProfile = world.UserInputProfile{KeyInteract: "e"}
	addTestEntity(t, &gameWorld, [2]int{1, 1}, map[string][]string{"pos": {}, "player": {}})
	doorID := addTestEntity(t, &gameWorld, [2]int{2, 1}, map[string][]string{
		"pos": {}, "impassable": {}, "interactable": {"unsupported"},
	})
	gameWorld.KeyDown = "e"

	result := interaction.ServiceInteraction{}.GetUpdateFunc(gameWorld)
	gameWorld = applyUpdate(t, result, gameWorld)

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
