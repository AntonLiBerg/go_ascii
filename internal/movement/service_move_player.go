package movement

import (
	"fmt"
	"go_ascii/internal/game"
	wrld "go_ascii/internal/world"
	cmp "go_ascii/internal/world/component"
)

type ServiceMovePlayer struct{}

func (s ServiceMovePlayer) GetUpdateFunc(w wrld.World) game.UpdateFunc {
	moveDelta := cmp.Position{}
	keyToClear := ""

	switch {
	case w.UserInput[string(w.UserInputProfile.KeyMoveUp)]:
		moveDelta = cmp.Position{Y: -1}
		keyToClear = string(w.UserInputProfile.KeyMoveUp)
	case w.UserInput[string(w.UserInputProfile.KeyMoveLeft)]:
		moveDelta = cmp.Position{X: -1}
		keyToClear = string(w.UserInputProfile.KeyMoveLeft)
	case w.UserInput[string(w.UserInputProfile.KeyMoveDown)]:
		moveDelta = cmp.Position{Y: 1}
		keyToClear = string(w.UserInputProfile.KeyMoveDown)
	case w.UserInput[string(w.UserInputProfile.KeyMoveRight)]:
		moveDelta = cmp.Position{X: 1}
		keyToClear = string(w.UserInputProfile.KeyMoveRight)
	default:
		return game.UpdateFunc{Order: 1}
	}

	return game.UpdateFunc{
		Order: 1,
		UpdateFunc: func(w *wrld.World) {
			w.HasChanged = true
			for _, eID := range w.Entities {
				if !w.HasComponent(eID, cmp.C_PLAYER) {
					continue
				}
				if err := tryGoToPosition(w, eID, moveDelta); err != nil {
					panic(err)
				}
			}
			w.UserInput[keyToClear] = false
		},
	}
}

func tryGoToPosition(w *wrld.World, eMover int, posDelta cmp.Position) error {
	moverPos, ok := w.Pos[eMover]
	if !ok {
		return fmt.Errorf("Mover entity not found!")
	}

	targetPos := cmp.Position{X: moverPos.X + posDelta.X, Y: moverPos.Y + posDelta.Y}
	targetID, ok := w.EByPos[targetPos]
	if !ok {
		return nil
	}
	if !canMakeMove(w, targetID) {
		return nil
	}

	w.Pos[eMover] = targetPos
	w.EByPos[targetPos] = eMover

	w.Pos[targetID] = moverPos
	w.EByPos[moverPos] = targetID
	return nil
}

func canMakeMove(w *wrld.World, targetID int) bool {
	if _, blocked := w.Impassable[targetID]; blocked {
		return false
	}
	return true
}
