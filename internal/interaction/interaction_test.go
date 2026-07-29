package interaction

import (
	"go_ascii/internal/world"
	"testing"
)

func TestServiceTurnsOnSingleNeighborAndClearsInteractKey(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	gameWorld.UserInputProfile = world.UserInputProfile{KeyInteract: "e"}
	addTestEntity(t, &gameWorld, [2]int{1, 1}, map[string][]string{
		"pos": {}, "ascii": {"o"}, "player": {},
	})
	machineID := addTestEntity(t, &gameWorld, [2]int{2, 1}, map[string][]string{
		"pos": {}, "ascii": {"R"}, "impassable": {}, "machine": {},
	})
	gameWorld.KeyDown = "e"

	result := ServiceTurnOnMachine{}.GetUpdateFunc(gameWorld)
	if result.UpdateFunc == nil {
		t.Fatal("expected interact update func")
	}
	result.UpdateFunc(&gameWorld)

	if !gameWorld.Machine[machineID] {
		t.Fatal("expected machine to be on after interaction")
	}
	if gameWorld.KeyDown != "" {
		t.Fatalf("expected interact key to be cleared, got %q", gameWorld.KeyDown)
	}
	if !gameWorld.HasChanged {
		t.Fatal("expected world to be marked changed after turning on machine")
	}
}

func TestServiceTurnsOnFirstNeighborMachine(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	gameWorld.UserInputProfile = world.UserInputProfile{KeyInteract: "e"}
	addTestEntity(t, &gameWorld, [2]int{1, 1}, map[string][]string{"pos": {}, "player": {}})
	eastMachine := addTestEntity(t, &gameWorld, [2]int{2, 1}, map[string][]string{"pos": {}, "machine": {}})
	westMachine := addTestEntity(t, &gameWorld, [2]int{0, 1}, map[string][]string{"pos": {}, "machine": {}})
	gameWorld.KeyDown = "e"

	result := ServiceTurnOnMachine{}.GetUpdateFunc(gameWorld)
	result.UpdateFunc(&gameWorld)

	if !gameWorld.Machine[westMachine] {
		t.Fatal("expected first machine in neighbor order to be on")
	}
	if gameWorld.Machine[eastMachine] {
		t.Fatal("expected later machine in neighbor order to remain off")
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
