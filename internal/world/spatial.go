package world

import component "go_ascii/internal"

func GetNeighbors(world World, target int) []int {
	targetPos, ok := world.Pos[target]
	if !ok {
		return nil
	}

	deltas := [...]component.Position{
		{X: -1, Y: -1},
		{X: -1, Y: 0},
		{X: -1, Y: 1},
		{X: 0, Y: 1},
		{X: 1, Y: 1},
		{X: 1, Y: 0},
		{X: 1, Y: -1},
		{X: 0, Y: -1},
	}

	var neighbors []int
	for _, delta := range deltas {
		pos := component.Position{X: targetPos.X + delta.X, Y: targetPos.Y + delta.Y}
		if eID, ok := world.EByPos[pos]; ok {
			neighbors = append(neighbors, eID)
		}
	}
	return neighbors
}
