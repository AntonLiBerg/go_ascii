package movement

import (
	wrld "go_ascii/internal/world"
	cmp "go_ascii/internal/world/component"
	usr "go_ascii/internal/world/user"
	"testing"
)

func TestServiceMovesPlayerWithWASD(t *testing.T) {
	tests := []struct {
		name          string
		key           string
		wantPlayerPos cmp.Position
	}{
		{name: "up blocked", key: "w", wantPlayerPos: cmp.Position{X: 1, Y: 1}},
		{name: "left open", key: "a", wantPlayerPos: cmp.Position{X: 0, Y: 1}},
		{name: "down open", key: "s", wantPlayerPos: cmp.Position{X: 1, Y: 2}},
		{name: "right open", key: "d", wantPlayerPos: cmp.Position{X: 2, Y: 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			world := newTestWorld(t)
			world.SetKeyDown(tt.key)

			result := ServiceMovePlayer{}.GetUpdateFunc(world)
			if result.UpdateFunc == nil {
				t.Fatal("expected movement update func")
			}
			result.UpdateFunc(&world)

			playerID := getSinglePlayerID(t, world)
			if got := world.Pos[playerID]; got != tt.wantPlayerPos {
				t.Fatalf("expected player at %+v, got %+v", tt.wantPlayerPos, got)
			}
			if got := world.EByPos[tt.wantPlayerPos]; got != playerID {
				t.Fatalf("expected reverse position index to point to player %d, got %d", playerID, got)
			}
			if world.UserInput[tt.key] {
				t.Fatalf("expected key %q to be cleared after movement", tt.key)
			}
		})
	}
}

func newTestWorld(t *testing.T) wrld.World {
	t.Helper()

	world := wrld.NewWorldEmpty()
	world.UserInputProfile = usr.UserInputProfile{
		KeyMoveUp:    "w",
		KeyMoveLeft:  "a",
		KeyMoveDown:  "s",
		KeyMoveRight: "d",
	}

	addTestEntity(t, &world, [2]int{1, 1}, map[cmp.ComponentName][]string{
		cmp.C_POS: {}, cmp.C_ASCII: {"o"}, cmp.C_PLAYER: {},
	})
	addTestEntity(t, &world, [2]int{1, 0}, map[cmp.ComponentName][]string{
		cmp.C_POS: {}, cmp.C_ASCII: {"#"}, cmp.C_IMPASSABLE: {},
	})
	for _, position := range [][2]int{{0, 1}, {1, 2}, {2, 1}} {
		addTestEntity(t, &world, position, map[cmp.ComponentName][]string{
			cmp.C_POS: {}, cmp.C_ASCII: {"."},
		})
	}

	return world
}

func addTestEntity(t *testing.T, world *wrld.World, position [2]int, components map[cmp.ComponentName][]string) {
	t.Helper()
	if err := world.AddEntity(position, components); err != nil {
		t.Fatalf("AddEntity returned error: %v", err)
	}
}

func getSinglePlayerID(t *testing.T, world wrld.World) int {
	t.Helper()
	if len(world.Player) != 1 {
		t.Fatalf("expected exactly one player, got %d", len(world.Player))
	}
	for id := range world.Player {
		return id
	}
	t.Fatal("expected player id")
	return 0
}
