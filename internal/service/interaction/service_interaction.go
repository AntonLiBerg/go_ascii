package interaction

import (
	component "go_ascii/internal"
	"go_ascii/internal/game"
	"go_ascii/internal/world"
)

type ServiceInteraction struct{}

func (s ServiceInteraction) GetUpdateFunc(w world.World) game.UpdateFunc {
	exitKey := w.UserInputProfile.KeyExit
	if w.ViewRoom != "" && exitKey != "" && w.KeyDown == exitKey {
		return game.UpdateFunc{
			Order: 1,
			UpdateFunc: func(w world.World) (world.World, error) {
				next := w
				next.ViewRoom = ""
				for _, playerID := range world.GetPlayerIDs(w) {
					position, ok := w.Pos[playerID]
					if ok {
						profile, ok := w.InputProfileForRoom(position.Room)
						if ok {
							next.UserInputProfile = profile
						}
					}
					break
				}
				next.HasChanged = true
				if next.KeyDown == exitKey {
					next.KeyDown = ""
				}
				return next, nil
			},
		}
	}
	if w.ViewRoom != "" {
		return game.UpdateFunc{Order: 1}
	}

	interactKey := w.UserInputProfile.KeyInteract
	if interactKey == "" || w.KeyDown != interactKey {
		return game.UpdateFunc{Order: 1}
	}

	playerID := -1
	for _, eID := range world.GetPlayerIDs(w) {
		playerID = eID
		break
	}

	targets := world.GetInteractableNeighbors(w, playerID)
	return game.UpdateFunc{
		Order: 1,
		UpdateFunc: func(w world.World) (world.World, error) {
			next := w
			if len(targets) == 1 {
				targetID := targets[0]
				if interaction, ok := w.Interactable[targetID]; ok {
					switch interaction.InteractionType {
					case component.InteractionTypeDoor:
						next = interactWithDoor(next, targetID)
					case component.InteractionTypeCommandTable:
						next = interactWithCommandTable(next, targetID)
					}
				}
			}
			if next.KeyDown == interactKey {
				next.KeyDown = ""
			}
			return next, nil
		},
	}
}

func interactWithDoor(w world.World, doorID int) world.World {
	next := w.Clone()
	if _, isClosed := next.Impassable[doorID]; isClosed {
		delete(next.Impassable, doorID)
	} else {
		next.Impassable[doorID] = component.Impassable{}
	}
	next.HasChanged = true
	return next
}

func interactWithCommandTable(w world.World, commandTableID int) world.World {
	position, ok := w.Pos[commandTableID]
	if !ok {
		return w
	}
	roomName, ok := w.Terminals[position]
	if !ok {
		return w
	}
	profile, ok := w.InputProfileForRoom(roomName)
	if !ok {
		return w
	}
	next := w
	next.UserInputProfile = profile
	next.ViewRoom = roomName
	next.HasChanged = true
	return next
}
