package movement

import (
	"fmt"
	component "go_ascii/internal"
	"go_ascii/internal/game"
	"go_ascii/internal/world"
)

type ServiceMovePlayer struct{}

func (s ServiceMovePlayer) GetUpdateFunc(w world.World) game.UpdateFunc {
	if w.ViewRoom != "" || w.KeyDown == "" {
		return game.UpdateFunc{Order: 1}
	}

	moveDelta := component.Position{}
	keyToClear := w.KeyDown
	switch w.KeyDown {
	case w.UserInputProfile.KeyMoveUp:
		moveDelta = component.Position{Y: -1}
	case w.UserInputProfile.KeyMoveLeft:
		moveDelta = component.Position{X: -1}
	case w.UserInputProfile.KeyMoveDown:
		moveDelta = component.Position{Y: 1}
	case w.UserInputProfile.KeyMoveRight:
		moveDelta = component.Position{X: 1}
	default:
		return game.UpdateFunc{Order: 1}
	}

	return game.UpdateFunc{
		Order: 1,
		UpdateFunc: func(w *world.World) {
			w.HasChanged = true
			for eID := range w.Player {
				if err := tryGoToPosition(w, eID, moveDelta); err != nil {
					panic(err)
				}
			}
			if w.KeyDown == keyToClear {
				w.KeyDown = ""
			}
		},
	}
}

func tryGoToPosition(w *world.World, moverID int, delta component.Position) error {
	moverPos, ok := w.Pos[moverID]
	if !ok {
		return fmt.Errorf("mover entity not found")
	}

	targetPos := component.Position{Room: moverPos.Room, X: moverPos.X + delta.X, Y: moverPos.Y + delta.Y}
	targetIDs := world.GetEntitiesAtPosition(*w, targetPos)
	if len(targetIDs) == 0 || !canMakeMove(w, targetIDs) {
		return nil
	}
	if portalTarget, isPortal := w.Portals[targetPos]; isPortal {
		targetPos = portalTarget
		targetIDs = world.GetEntitiesAtPosition(*w, targetPos)
		if len(targetIDs) == 0 || !canMakeMove(w, targetIDs) {
			return nil
		}
	}

	w.Pos[moverID] = targetPos
	w.EByPos[targetPos] = moverID
	if moverPos.Room != targetPos.Room {
		w.SetInputProfileForRoom(targetPos.Room)
	}
	entitiesAtOldPosition := world.GetEntitiesAtPosition(*w, moverPos)
	if len(entitiesAtOldPosition) == 0 {
		delete(w.EByPos, moverPos)
	} else {
		w.EByPos[moverPos] = entitiesAtOldPosition[len(entitiesAtOldPosition)-1]
	}
	return nil
}

func canMakeMove(w *world.World, targetIDs []int) bool {
	for _, targetID := range targetIDs {
		if _, blocked := w.Impassable[targetID]; blocked {
			return false
		}
	}
	return true
}
