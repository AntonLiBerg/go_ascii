package interaction

import (
	component "go_ascii/internal"
	"go_ascii/internal/game"
	"go_ascii/internal/world"
)

type ServiceInteraction struct{}

func (s ServiceInteraction) GetUpdateFunc(w world.World) game.UpdateFunc {
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
				if _, isDoor := w.Door[targetID]; isDoor {
					interactWithDoor(w, targetID)
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
