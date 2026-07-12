package world

import (
	cmp "go_ascii/internal/world/component"
	"slices"
	"testing"
)

func TestGetNeighbors(t *testing.T) {
	world := NewWorldEmpty()
	target := world.MakeNewEntityId()
	world = world.AddPosition(target, cmp.Position{X: 1, Y: 1})

	positions := []cmp.Position{
		{X: 0, Y: 0},
		{X: 0, Y: 1},
		{X: 0, Y: 2},
		{X: 1, Y: 2},
		{X: 2, Y: 2},
		{X: 2, Y: 1},
		{X: 2, Y: 0},
		{X: 1, Y: 0},
	}
	want := make([]int, 0, len(positions))
	for _, position := range positions {
		neighbor := world.MakeNewEntityId()
		world = world.AddPosition(neighbor, position)
		want = append(want, neighbor)
	}

	if got := GetNeighbors(world, target, nil); !slices.Equal(got, want) {
		t.Fatalf("expected neighbors %v, got %v", want, got)
	}
	if got := GetNeighbors(world, target, []cmp.ComponentName{cmp.C_IMPASSABLE}); len(got) != 0 {
		t.Fatalf("expected no impassable neighbors, got %v", got)
	}
	if got := GetNeighbors(world, -1, nil); len(got) != 0 {
		t.Fatalf("expected no neighbors for missing target, got %v", got)
	}
}
