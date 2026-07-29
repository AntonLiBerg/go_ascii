package movement

import (
	"fmt"
	"go_ascii/internal/game"
	"go_ascii/internal/world"
)

type ServiceMovePlayer struct{}

func (s ServiceMovePlayer) GetUpdateFunc(w world.World) game.UpdateFunc {
	if w.KeyDown == "" {
		return game.UpdateFunc{Order: 1}
	}

	moveDelta := world.Position{}
	keyToClear := w.KeyDown
	switch w.KeyDown {
	case w.UserInputProfile.KeyMoveUp:
		moveDelta = world.Position{Y: -1}
	case w.UserInputProfile.KeyMoveLeft:
		moveDelta = world.Position{X: -1}
	case w.UserInputProfile.KeyMoveDown:
		moveDelta = world.Position{Y: 1}
	case w.UserInputProfile.KeyMoveRight:
		moveDelta = world.Position{X: 1}
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

func tryGoToPosition(w *world.World, moverID int, delta world.Position) error {
	moverPos, ok := w.Pos[moverID]
	if !ok {
		return fmt.Errorf("mover entity not found")
	}

	targetPos := world.Position{X: moverPos.X + delta.X, Y: moverPos.Y + delta.Y}
	targetID, ok := w.EByPos[targetPos]
	if !ok || !canMakeMove(w, targetID) {
		return nil
	}

	w.Pos[moverID] = targetPos
	w.EByPos[targetPos] = moverID
	w.Pos[targetID] = moverPos
	w.EByPos[moverPos] = targetID
	return nil
}

func canMakeMove(w *world.World, targetID int) bool {
	_, blocked := w.Impassable[targetID]
	return !blocked
}
