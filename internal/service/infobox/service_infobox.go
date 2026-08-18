package infobox

import (
	"go_ascii/internal/game"
	"go_ascii/internal/world"
)

type ServiceInfobox struct{}

func (s ServiceInfobox) GetUpdateFunc(w world.World) game.UpdateFunc {
	if w.InfoboxContent == "" {
		return game.UpdateFunc{}
	}
	return game.UpdateFunc{
		Order: 10,
		UpdateFunc: func(w world.World) (world.World, error) {
			nw := w.Clone()
			infobox, err := AI.SetInfoboxText(nw.UIContent["infobox"], nw.InfoboxContent)
			if err != nil {
				return w, err
			}
			nw.UIContent["infobox"] = infobox
			return nw, nil
		},
	}
}
