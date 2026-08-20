package movement_test

import (
	component "go_ascii/internal"
	"go_ascii/internal/service/movement"
	"go_ascii/internal/world"
	"testing"
)

func TestServiceMovesPlayerWithWASD(t *testing.T) {
	tests := []struct {
		name          string
		key           string
		wantPlayerPos component.Position
	}{
		{name: "up blocked", key: "w", wantPlayerPos: component.Position{X: 1, Y: 1}},
		{name: "left open", key: "a", wantPlayerPos: component.Position{X: 0, Y: 1}},
		{name: "down open", key: "s", wantPlayerPos: component.Position{X: 1, Y: 2}},
		{name: "right open", key: "d", wantPlayerPos: component.Position{X: 2, Y: 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gameWorld := newTestWorld(t)
			gameWorld.KeyDown = tt.key

			result := movement.ServiceMovePlayer{}.GetUpdateFunc(gameWorld)
			if result.UpdateFunc == nil {
				t.Fatal("expected movement update func")
			}
			input := gameWorld
			next := applyUpdate(t, result, gameWorld)

			playerID := getSinglePlayerID(t, next)
			if got := input.Pos[playerID]; got != (component.Position{X: 1, Y: 1}) {
				t.Fatalf("expected movement update not to mutate its input, got %+v", got)
			}
			gameWorld = next
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
	result := movement.ServiceMovePlayer{}.GetUpdateFunc(gameWorld)

	gameWorld.KeyDown = "e"
	gameWorld = applyUpdate(t, result, gameWorld)

	if gameWorld.KeyDown != "e" {
		t.Fatalf("expected newer key to remain, got %q", gameWorld.KeyDown)
	}
}

func TestMovementKeepsUnderlyingDoorInPlace(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	gameWorld.UserInputProfile = world.UserInputProfile{KeyMoveRight: "d"}
	playerID := addTestEntity(t, &gameWorld, [2]int{1, 1}, map[string][]string{
		"pos": {}, "ascii": {"o"}, "player": {},
	})
	doorID := addTestEntity(t, &gameWorld, [2]int{2, 1}, map[string][]string{
		"pos": {}, "ascii": {"D"},
	})
	addTestEntity(t, &gameWorld, [2]int{3, 1}, map[string][]string{
		"pos": {}, "ascii": {"."},
	})

	gameWorld.KeyDown = "d"
	firstMove := movement.ServiceMovePlayer{}.GetUpdateFunc(gameWorld)
	gameWorld = applyUpdate(t, firstMove, gameWorld)

	doorPos := component.Position{X: 2, Y: 1}
	if gameWorld.Pos[playerID] != doorPos || gameWorld.Pos[doorID] != doorPos {
		t.Fatalf("expected player and door at %+v, got player=%+v door=%+v", doorPos, gameWorld.Pos[playerID], gameWorld.Pos[doorID])
	}
	if gameWorld.EByPos[doorPos] != playerID {
		t.Fatal("expected player to be indexed above door")
	}

	gameWorld.KeyDown = "d"
	secondMove := movement.ServiceMovePlayer{}.GetUpdateFunc(gameWorld)
	gameWorld = applyUpdate(t, secondMove, gameWorld)

	if gameWorld.Pos[doorID] != doorPos {
		t.Fatalf("expected door to remain at %+v, got %+v", doorPos, gameWorld.Pos[doorID])
	}
	if gameWorld.EByPos[doorPos] != doorID {
		t.Fatal("expected door index to be restored after player leaves")
	}
}

func TestMovementTraversesPortalBetweenRooms(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	bridgeProfile := world.UserInputProfile{KeyMoveUp: "w", KeyMoveDown: "s", KeyInteract: "b"}
	commsProfile := world.UserInputProfile{KeyMoveUp: "w", KeyMoveDown: "s", KeyInteract: "c"}
	gameWorld.InputProfiles["bridge"] = bridgeProfile
	gameWorld.InputProfiles["comms"] = commsProfile
	gameWorld.InputProfileByRoom["bridge"] = "bridge"
	gameWorld.InputProfileByRoom["comms"] = "comms"
	gameWorld.UserInputProfile = bridgeProfile
	if err := gameWorld.AddEntityInRoom("bridge", [2]int{1, 1}, map[string][]string{
		"pos": {}, "ascii": {"o"}, "player": {},
	}); err != nil {
		t.Fatalf("add player: %v", err)
	}
	if err := gameWorld.AddEntityInRoom("bridge", [2]int{1, 0}, map[string][]string{
		"pos": {}, "ascii": {"."},
	}); err != nil {
		t.Fatalf("add bridge portal ground: %v", err)
	}
	if err := gameWorld.AddEntityInRoom("comms", [2]int{2, 3}, map[string][]string{
		"pos": {}, "ascii": {"."},
	}); err != nil {
		t.Fatalf("add comms portal ground: %v", err)
	}
	if err := gameWorld.AddEntityInRoom("comms", [2]int{2, 4}, map[string][]string{
		"pos": {}, "ascii": {"."},
	}); err != nil {
		t.Fatalf("add comms interior ground: %v", err)
	}
	bridgePortal := component.Position{Room: "bridge", X: 1, Y: 0}
	commsPortal := component.Position{Room: "comms", X: 2, Y: 3}
	gameWorld.Portals[bridgePortal] = commsPortal
	gameWorld.Portals[commsPortal] = bridgePortal

	gameWorld.KeyDown = "w"
	move := movement.ServiceMovePlayer{}.GetUpdateFunc(gameWorld)
	gameWorld = applyUpdate(t, move, gameWorld)

	if got := gameWorld.Pos[0]; got != commsPortal {
		t.Fatalf("expected player at comms portal %+v, got %+v", commsPortal, got)
	}
	if gameWorld.EByPos[commsPortal] != 0 {
		t.Fatal("expected player indexed at destination portal")
	}
	if gameWorld.EByPos[bridgePortal] != 1 {
		t.Fatal("expected bridge ground restored in position index")
	}
	if gameWorld.UserInputProfile != commsProfile {
		t.Fatalf("expected comms input profile, got %+v", gameWorld.UserInputProfile)
	}

	gameWorld.KeyDown = "s"
	move = movement.ServiceMovePlayer{}.GetUpdateFunc(gameWorld)
	gameWorld = applyUpdate(t, move, gameWorld)
	gameWorld.KeyDown = "w"
	move = movement.ServiceMovePlayer{}.GetUpdateFunc(gameWorld)
	gameWorld = applyUpdate(t, move, gameWorld)
	if got := gameWorld.Pos[0]; got != bridgePortal {
		t.Fatalf("expected reverse portal traversal to %+v, got %+v", bridgePortal, got)
	}
	if gameWorld.UserInputProfile != bridgeProfile {
		t.Fatalf("expected bridge input profile, got %+v", gameWorld.UserInputProfile)
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
