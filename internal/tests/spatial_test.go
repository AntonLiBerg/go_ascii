package tests

import (
	component "go_ascii/internal"
	"go_ascii/internal/world"
	"slices"
	"testing"
)

func TestGetNeighbors(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	addEntity := func(pos [2]int) int {
		t.Helper()
		eID := len(gameWorld.Entities)
		if err := gameWorld.AddEntity(pos, map[string][]string{"pos": {}}); err != nil {
			t.Fatalf("AddEntity returned error: %v", err)
		}
		return eID
	}

	target := addEntity([2]int{1, 1})
	positions := [][2]int{
		{0, 0},
		{0, 1},
		{0, 2},
		{1, 2},
		{2, 2},
		{2, 1},
		{2, 0},
		{1, 0},
	}
	want := make([]int, 0, len(positions))
	for _, position := range positions {
		want = append(want, addEntity(position))
	}

	if got := world.GetNeighbors(gameWorld, target); !slices.Equal(got, want) {
		t.Fatalf("expected neighbors %v, got %v", want, got)
	}
	if got := world.GetNeighbors(gameWorld, -1); len(got) != 0 {
		t.Fatalf("expected no neighbors for missing target, got %v", got)
	}
}

func TestGetNeighborsOnlyReturnsEntitiesInSameRoom(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	if err := gameWorld.AddEntityInRoom("bridge", [2]int{1, 1}, map[string][]string{"pos": {}}); err != nil {
		t.Fatalf("add target: %v", err)
	}
	if err := gameWorld.AddEntityInRoom("bridge", [2]int{2, 1}, map[string][]string{"pos": {}}); err != nil {
		t.Fatalf("add bridge neighbor: %v", err)
	}
	if err := gameWorld.AddEntityInRoom("comms", [2]int{2, 1}, map[string][]string{"pos": {}}); err != nil {
		t.Fatalf("add comms entity: %v", err)
	}

	if got := world.GetNeighbors(gameWorld, 0); !slices.Equal(got, []int{1}) {
		t.Fatalf("expected only bridge neighbor, got %v", got)
	}
	if gameWorld.Pos[2] != (component.Position{Room: "comms", X: 2, Y: 1}) {
		t.Fatalf("expected comms entity position, got %+v", gameWorld.Pos[2])
	}
}

func TestGetInteractableNeighborsUsesInteractableComponent(t *testing.T) {
	gameWorld := world.NewWorldEmpty()
	target := addTestEntity(t, &gameWorld, [2]int{1, 1}, map[string][]string{"pos": {}})
	addTestEntity(t, &gameWorld, [2]int{0, 1}, map[string][]string{"pos": {}, "door": {}})
	interactable := addTestEntity(t, &gameWorld, [2]int{2, 1}, map[string][]string{
		"pos": {}, "interactable": {"helm"},
	})

	if got := world.GetInteractableNeighbors(gameWorld, target); !slices.Equal(got, []int{interactable}) {
		t.Fatalf("expected only entity with Interactable component, got %v", got)
	}
}
