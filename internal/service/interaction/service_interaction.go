package interaction

import (
	"fmt"
	component "go_ascii/internal"
	"go_ascii/internal/game"
	"go_ascii/internal/world"
)

type ServiceInteraction struct{}

func (s ServiceInteraction) GetUpdateFunc(w world.World) game.UpdateFunc {
	//
	//if using terminal, did you exit?
	//
	exitKey := w.UserInputProfile.KeyExit
	if exitKey != "" && w.KeyDown == exitKey {
		return game.UpdateFunc{
			Order: 1,
			UpdateFunc: func(w world.World) (world.World, error) {
				next := w
				playerID := world.GetPlayerIDs(w)[0]
				position, ok := w.Pos[playerID]
				if ok {
					profile, ok := w.InputProfileForRoom(position.Room)
					if ok {
						next.UserInputProfile = profile
					}
				}
				next.Room = position.Room
				next.HasChanged = true
				if next.KeyDown == exitKey {
					next.KeyDown = ""
				}
				return next, nil
			},
		}
	}
	//
	//not using terminal, was there an interaction?
	//
	interactKey := w.UserInputProfile.KeyInteract
	if interactKey == "" || w.KeyDown != interactKey {
		return game.UpdateFunc{Order: 1}
	}
	//
	//interact with one of the neighbors
	//
	playerID := world.GetPlayerIDs(w)[0]
	targets := world.GetInteractableNeighbors(w, playerID)
	return game.UpdateFunc{
		Order: 1,
		UpdateFunc: func(w world.World) (world.World, error) {
			next := w
			var err error
			if len(targets) == 1 {
				targetID := targets[0]
				if interaction, ok := w.Interactable[targetID]; ok {
					switch interaction.InteractionType {
					case component.InteractionTypeDoor:
						next = interactWithDoor(next, targetID)
					case component.InteractionTypeTerminal:
						next,err = interactWithTerminal(next, targetID)
					}
				}
			}
			if next.KeyDown == interactKey {
				next.KeyDown = ""
			}
			return next, err
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

func interactWithTerminal(w world.World, terminalID int) (world.World,error) {
	position, ok := w.Pos[terminalID]
	if !ok {
		return w,fmt.Errorf("Position for terminal not found!")
	}
	roomName, ok := w.Terminals[position]
	if !ok {
		return w,fmt.Errorf("Terminal not found!")
	}
	profile, ok := w.InputProfileForRoom(roomName)
	if !ok {
		return w,fmt.Errorf("inputprofile not found!")
	}
	next := w
	next.UserInputProfile = profile
	next.HasChanged = true
	next.Room = roomName
	return next,nil
}
