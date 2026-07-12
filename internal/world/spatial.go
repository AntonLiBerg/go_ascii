package world

import cmp "go_ascii/internal/world/component"

func GetNeighbors(world World, target int, filterComponents []cmp.ComponentName) []int {
	targetPos, ok := world.Pos[target]
	if !ok {
		return []int{}
	}

	deltas := []cmp.Position{
		{X: -1, Y: -1},
		{X: -1, Y: 0},
		{X: -1, Y: 1},
		{X: 0, Y: 1},
		{X: 1, Y: 1},
		{X: 1, Y: 0},
		{X: 1, Y: -1},
		{X: 0, Y: -1},
	}

	neighbors := []int{}
	for _, delta := range deltas {
		pos := cmp.Position{X: targetPos.X + delta.X, Y: targetPos.Y + delta.Y}
		eID, ok := world.EByPos[pos]
		if !ok || eID == target || !hasComponents(world, eID, filterComponents) {
			continue
		}
		neighbors = append(neighbors, eID)
	}

	return neighbors
}

func hasComponents(world World, eID int, components []cmp.ComponentName) bool {
	for _, component := range components {
		if !world.HasComponent(eID, component) {
			return false
		}
	}
	return true
}
