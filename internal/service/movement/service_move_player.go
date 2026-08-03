package movement

import (
	"fmt"
	component "go_ascii/internal"
	"go_ascii/internal/game"
	"go_ascii/internal/world"
)

type ServiceMovePlayer struct{}

func (s ServiceMovePlayer) GetUpdateFunc(w world.World) game.UpdateFunc {
	//
	// we pushed no button
	//
	if w.KeyDown == "" {
		return game.UpdateFunc{Order: 1}
	}

	//
	//button pushed but did we move?
	//
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

	//
	// move player
	//
	return game.UpdateFunc{
		Order: 1,
		UpdateFunc: func(w world.World) (world.World, error) {
			next := w.Clone()
			next.HasChanged = true
			for _, eID := range world.GetPlayerIDs(w) {
				var err error
				next, err = tryGoToPosition(next, eID, moveDelta)
				if err != nil {
					return w, err
				}
			}
			if next.KeyDown == keyToClear {
				next.KeyDown = ""
			}
			return next, nil
		},
	}
}

func tryGoToPosition(w world.World, moverID int, delta component.Position) (world.World, error) {
	moverPos, ok := w.Pos[moverID]
	if !ok {
		return w, fmt.Errorf("mover entity not found")
	}

	targetPos := component.Position{Room: moverPos.Room, X: moverPos.X + delta.X, Y: moverPos.Y + delta.Y}
	targetIDs := world.GetEntitiesAtPosition(w, targetPos)
	if len(targetIDs) == 0 || !canMakeMove(w, targetIDs) {
		return w, nil
	}
	//
	// moving into portal?
	//
	if portalTarget, isPortal := w.Portals[targetPos]; isPortal {
		targetPos = portalTarget
		targetIDs = world.GetEntitiesAtPosition(w, targetPos)
		if len(targetIDs) == 0 || !canMakeMove(w, targetIDs) {
			return w, nil
		}
	}

	next := w.Clone()
	next.Pos[moverID] = targetPos
	next.EByPos[targetPos] = moverID
	//
	// are we going through a portal?
	//
	if moverPos.Room != targetPos.Room {
		if profile, ok := next.InputProfileForRoom(targetPos.Room); ok {
			next.UserInputProfile = profile
			next.Room = targetPos.Room
		}
	}
	entitiesAtOldPosition := world.GetEntitiesAtPosition(next, moverPos)
	if len(entitiesAtOldPosition) == 0 {
		delete(next.EByPos, moverPos)
	} else {
		next.EByPos[moverPos] = entitiesAtOldPosition[len(entitiesAtOldPosition)-1]
	}
	return next, nil
}

func canMakeMove(w world.World, targetIDs []int) bool {
	for _, targetID := range targetIDs {
		if _, blocked := w.Impassable[targetID]; blocked {
			return false
		}
	}
	return true
}
