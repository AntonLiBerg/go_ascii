package interaction

import (
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
		"pos": {}, "ascii": {"D"}, "impassable": {}, "door": {},
	})
	gameWorld.KeyDown = "e"

	result := ServiceInteraction{}.GetUpdateFunc(gameWorld)
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

func TestServiceClosesOpenDoor(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	gameWorld.UserInputProfile = world.UserInputProfile{KeyInteract: "e"}
	addTestEntity(t, &gameWorld, [2]int{1, 1}, map[string][]string{"pos": {}, "player": {}})
	doorID := addTestEntity(t, &gameWorld, [2]int{2, 1}, map[string][]string{"pos": {}, "door": {}})
	gameWorld.KeyDown = "e"

	result := ServiceInteraction{}.GetUpdateFunc(gameWorld)
	result.UpdateFunc(&gameWorld)

	if _, isClosed := gameWorld.Impassable[doorID]; !isClosed {
		t.Fatal("expected door to be closed after interaction")
	}
}

func TestServiceIgnoresMultipleNeighborDoors(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	gameWorld.UserInputProfile = world.UserInputProfile{KeyInteract: "e"}
	addTestEntity(t, &gameWorld, [2]int{1, 1}, map[string][]string{"pos": {}, "player": {}})
	eastDoor := addTestEntity(t, &gameWorld, [2]int{2, 1}, map[string][]string{"pos": {}, "impassable": {}, "door": {}})
	westDoor := addTestEntity(t, &gameWorld, [2]int{0, 1}, map[string][]string{"pos": {}, "impassable": {}, "door": {}})
	gameWorld.KeyDown = "e"

	result := ServiceInteraction{}.GetUpdateFunc(gameWorld)
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
		name      string
		component string
	}{
		{name: "helm", component: "helm"},
		{name: "command table", component: "commandtable"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gameWorld := world.NewWorldEmpty()
			gameWorld.UserInputProfile = world.UserInputProfile{KeyInteract: "e"}
			addTestEntity(t, &gameWorld, [2]int{1, 1}, map[string][]string{"pos": {}, "player": {}})
			targetID := addTestEntity(t, &gameWorld, [2]int{2, 1}, map[string][]string{
				"pos": {}, "impassable": {}, tt.component: {},
			})
			gameWorld.KeyDown = "e"

			result := ServiceInteraction{}.GetUpdateFunc(gameWorld)
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

func addTestEntity(t *testing.T, gameWorld *world.World, position [2]int, components map[string][]string) int {
	t.Helper()
	eID := len(gameWorld.Entities)
	if err := gameWorld.AddEntity(position, components); err != nil {
		t.Fatalf("AddEntity returned error: %v", err)
	}
	return eID
}
