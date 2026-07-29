package world

import (
	"slices"
	"testing"
)

func TestGetNeighbors(t *testing.T) {
	gameWorld := NewWorldEmpty()
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

	if got := GetNeighbors(gameWorld, target); !slices.Equal(got, want) {
		t.Fatalf("expected neighbors %v, got %v", want, got)
	}
	if got := GetNeighbors(gameWorld, -1); len(got) != 0 {
		t.Fatalf("expected no neighbors for missing target, got %v", got)
	}
}
