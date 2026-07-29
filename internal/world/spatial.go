package world

func GetNeighbors(world World, target int) []int {
	targetPos, ok := world.Pos[target]
	if !ok {
		return nil
	}

	deltas := [...]Position{
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
		pos := Position{X: targetPos.X + delta.X, Y: targetPos.Y + delta.Y}
		if eID, ok := world.EByPos[pos]; ok {
			neighbors = append(neighbors, eID)
		}
	}
	return neighbors
}
