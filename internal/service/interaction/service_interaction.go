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
			UpdateFunc: func(w *world.World) {
				w.ViewRoom = ""
				for playerID := range w.Player {
					w.SetInputProfileForRoom(w.Pos[playerID].Room)
					break
				}
				w.HasChanged = true
				if w.KeyDown == exitKey {
					w.KeyDown = ""
				}
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
	for eID := range w.Player {
		playerID = eID
		break
	}

	targets := world.GetInteractableNeighbors(w, playerID)
	return game.UpdateFunc{
		Order: 1,
		UpdateFunc: func(w *world.World) {
			if len(targets) == 1 {
				targetID := targets[0]
				interaction := w.Interactable[targetID]
				switch interaction.InteractionType {
				case component.InteractionTypeDoor:
					interactWithDoor(w, targetID)
				case component.InteractionTypeHelm:
					interactWithHelm(w, targetID)
				case component.InteractionTypeCommandTable:
					interactWithCommandTable(w, targetID)
				}
			}
			if w.KeyDown == interactKey {
				w.KeyDown = ""
			}
		},
	}
}

func interactWithDoor(w *world.World, doorID int) {
	if _, isClosed := w.Impassable[doorID]; isClosed {
		delete(w.Impassable, doorID)
	} else {
		w.Impassable[doorID] = component.Impassable{}
	}
	w.HasChanged = true
}

func interactWithHelm(*world.World, int) {
	//show UI
}

func interactWithCommandTable(w *world.World, commandTableID int) {
	position, ok := w.Pos[commandTableID]
	if !ok {
		return
	}
	roomName, ok := w.Terminals[position]
	if !ok || !w.SetInputProfileForRoom(roomName) {
		return
	}
	w.ViewRoom = roomName
	w.HasChanged = true
}
