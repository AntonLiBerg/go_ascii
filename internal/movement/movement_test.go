package movement

import (
	"go_ascii/internal/world"
	"testing"
)

func TestServiceMovesPlayerWithWASD(t *testing.T) {
	tests := []struct {
		name          string
		key           string
		wantPlayerPos world.Position
	}{
		{name: "up blocked", key: "w", wantPlayerPos: world.Position{X: 1, Y: 1}},
		{name: "left open", key: "a", wantPlayerPos: world.Position{X: 0, Y: 1}},
		{name: "down open", key: "s", wantPlayerPos: world.Position{X: 1, Y: 2}},
		{name: "right open", key: "d", wantPlayerPos: world.Position{X: 2, Y: 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gameWorld := newTestWorld(t)
			gameWorld.KeyDown = tt.key

			result := ServiceMovePlayer{}.GetUpdateFunc(gameWorld)
			if result.UpdateFunc == nil {
				t.Fatal("expected movement update func")
			}
			result.UpdateFunc(&gameWorld)

			playerID := getSinglePlayerID(t, gameWorld)
			if got := gameWorld.Pos[playerID]; got != tt.wantPlayerPos {
				t.Fatalf("expected player at %+v, got %+v", tt.wantPlayerPos, got)
			}
			if got := gameWorld.EByPos[tt.wantPlayerPos]; got != playerID {
				t.Fatalf("expected reverse position index to point to player %d, got %d", playerID, got)
			}
			if gameWorld.KeyDown != "" {
				t.Fatalf("expected key to be cleared after movement, got %q", gameWorld.KeyDown)
			}
		})
	}
}

func TestMovementDoesNotClearNewerKey(t *testing.T) {
	gameWorld := newTestWorld(t)
	gameWorld.KeyDown = "a"
	result := ServiceMovePlayer{}.GetUpdateFunc(gameWorld)

	gameWorld.KeyDown = "e"
	result.UpdateFunc(&gameWorld)

	if gameWorld.KeyDown != "e" {
		t.Fatalf("expected newer key to remain, got %q", gameWorld.KeyDown)
	}
}

func newTestWorld(t *testing.T) world.World {
	t.Helper()

	gameWorld := world.NewWorldEmpty()
	gameWorld.UserInputProfile = world.UserInputProfile{
		KeyMoveUp:    "w",
		KeyMoveLeft:  "a",
		KeyMoveDown:  "s",
		KeyMoveRight: "d",
	}

	addTestEntity(t, &gameWorld, [2]int{1, 1}, map[string][]string{
		"pos": {}, "ascii": {"o"}, "player": {},
	})
	addTestEntity(t, &gameWorld, [2]int{1, 0}, map[string][]string{
		"pos": {}, "ascii": {"#"}, "impassable": {},
	})
	for _, position := range [][2]int{{0, 1}, {1, 2}, {2, 1}} {
		addTestEntity(t, &gameWorld, position, map[string][]string{
			"pos": {}, "ascii": {"."},
		})
	}

	return gameWorld
}

func addTestEntity(t *testing.T, gameWorld *world.World, position [2]int, components map[string][]string) int {
	t.Helper()
	eID := len(gameWorld.Entities)
	if err := gameWorld.AddEntity(position, components); err != nil {
		t.Fatalf("AddEntity returned error: %v", err)
	}
	return eID
}

func getSinglePlayerID(t *testing.T, gameWorld world.World) int {
	t.Helper()
	if len(gameWorld.Player) != 1 {
		t.Fatalf("expected exactly one player, got %d", len(gameWorld.Player))
	}
	for id := range gameWorld.Player {
		return id
	}
	t.Fatal("expected player id")
	return 0
}
