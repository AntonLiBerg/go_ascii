package render

import (
	"fmt"
	component "go_ascii/internal"
	"go_ascii/internal/game"
	"go_ascii/internal/world"
)

type ServiceDrawOnTerminal struct{}

func (s ServiceDrawOnTerminal) GetUpdateFunc(w world.World) game.UpdateFunc {
	return game.UpdateFunc{
		Order: 100,
		UpdateFunc: func(w *world.World) {
			if w.IterationNr == 1 || w.HasChanged {
				UpdateTerminal(*w)
				w.HasChanged = false
			}
		},
	}
}

func UpdateTerminal(world world.World) {
	fmt.Print("\033[2J\033[H")

	activeRoom := ""
	hasActiveRoom := false
	for playerID := range world.Player {
		if pos, ok := world.Pos[playerID]; ok {
			activeRoom = pos.Room
			hasActiveRoom = true
			break
		}
	}

	maxY := 0
	for _, eID := range world.Entities {
		pos, okPos := world.Pos[eID]
		ascii, okASCII := world.Ascii[eID]
		if !okPos || !okASCII || (hasActiveRoom && pos.Room != activeRoom) || pos.X < 0 || pos.Y < 0 || isCovered(world, eID, pos) {
			continue
		}

		fmt.Printf("\033[%d;%dH%c", pos.Y+1, pos.X+1, ascii.Ascii)
		if pos.Y > maxY {
			maxY = pos.Y
		}
	}

	fmt.Printf("\033[%d;1H", maxY+2)
}

func isCovered(world world.World, eID int, pos component.Position) bool {
	layer := world.Layer[eID].Nr
	_, isPlayer := world.Player[eID]
	for _, otherID := range world.Entities {
		if otherID == eID || world.Pos[otherID] != pos {
			continue
		}
		if _, hasASCII := world.Ascii[otherID]; !hasASCII {
			continue
		}

		otherLayer := world.Layer[otherID].Nr
		_, otherIsPlayer := world.Player[otherID]
		if otherLayer > layer || (otherLayer == layer && otherIsPlayer && !isPlayer) {
			return true
		}
		if otherLayer == layer && otherIsPlayer == isPlayer && otherID > eID {
			return true
		}
	}
	return false
}
