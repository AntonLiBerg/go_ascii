package render

import (
	"fmt"
	"go_ascii/internal/game"
	wrld "go_ascii/internal/world"
)

type ServiceDrawOnTerminal struct{}

func (s ServiceDrawOnTerminal) GetUpdateFunc(w wrld.World) game.UpdateFunc {
	return game.UpdateFunc{
		Order: 100,
		UpdateFunc: func(w *wrld.World) {
			if w.IterationNr == 1 || w.HasChanged {
				UpdateTerminal(*w)
				w.HasChanged = false
			}
		},
	}
}

func UpdateTerminal(world wrld.World) {
	fmt.Print("\033[2J\033[H")

	maxY := 0
	for _, eID := range world.Entities {
		pos, okPos := world.Pos[eID]
		ascii, okASCII := world.Ascii[eID]
		if !okPos || !okASCII || pos.X < 0 || pos.Y < 0 {
			continue
		}

		fmt.Printf("\033[%d;%dH%c", pos.Y+1, pos.X+1, ascii.Ascii)
		if pos.Y > maxY {
			maxY = pos.Y
		}
	}

	fmt.Printf("\033[%d;1H", maxY+2)
}
