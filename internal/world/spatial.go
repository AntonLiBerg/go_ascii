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
		pos := component.Position{Room: targetPos.Room, X: targetPos.X + delta.X, Y: targetPos.Y + delta.Y}
		neighbors = append(neighbors, GetEntitiesAtPosition(world, pos)...)
	}
	return neighbors
}

func GetEntitiesAtPosition(world World, target component.Position) []int {
	var entities []int
	for _, eID := range world.Entities {
		if pos, ok := world.Pos[eID]; ok && pos == target {
			entities = append(entities, eID)
		}
	}
	return entities
}

func GetInteractableNeighbors(world World, target int) []int {
	var interactable []int
	for _, neighborID := range GetNeighbors(world, target) {
		if isInteractable(world, neighborID) {
			interactable = append(interactable, neighborID)
		}
	}
	return interactable
}

func isInteractable(world World, eID int) bool {
	if door, ok := world.Door[eID]; ok && door.IsInteractable {
		return true
	}
	if helm, ok := world.Helm[eID]; ok && helm.IsInteractable {
		return true
	}
	if commandTable, ok := world.CommandTable[eID]; ok && commandTable.IsInteractable {
		return true
	}
	return false
}
